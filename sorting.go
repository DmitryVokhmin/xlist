// sorting.go
// Sorts the items in the container
// Created by Vokhmin D.A. 03.2025

package xlist

import (
	"math"
	"runtime"
	"sync"
)

// Sort sorts the list according to the provided comparison function.
//
// Parameters:
//   - compare: A function that compares two elements.
//   - Returns true when `a` should be before `b`, otherwise false.
func (p *XList[T]) Sort(compare func(a, b T) bool) {
	// Performance benchmarks:
	//   - Single thread:  10,000 items ~0.18s, 100,000 items ~28.00s
	//   - Two threads:    10,000 items ~0.15s, 100,000 items ~28.00s
	//   - 16 threads:     10,000 items ~0.06s, 100,000 items ~3.53s, 200,000 items ~12-16s
	//   - 16 threads with indexes: 10,000 items ~0.05s, 100,000 items ~1.04-1.47s,
	//     200,000 items ~2.67-3s, 1,000,000 items ~44.69-55.46s
	p.ScanSort(compare)
}

// ScanSort : sort list with simple move and order algoritm.
// Parameters:
//   - compare: A function that compares two elements.
//   - Returns true when `a` should be before `b`, otherwise false.
func (p *XList[T]) ScanSort(compare func(a, b T) bool) {
	wg := &sync.WaitGroup{}

	if p.Size() < 2 {
		return
	}

	// Lock component work until sort is done
	p.mtx.Lock()
	defer p.mtx.Unlock()

	// Channel for supplying elements for sorting
	unsortedChanel := make(chan *xlistObj[T], runtime.NumCPU()) // Add buffer for better performance

	// Detach all elements from the list
	unsortedHome := p.home.next // starts from the second element, first element put in sorted array
	p.home.next = nil
	p.end = p.home
	indexGrains := p.grainsAccordingToSize(p.Size())
	p.size = 1

	p.sortContext = &sortContext[T]{
		canRead: true,
		grains:  indexGrains,
		indexes: make([]*indexPair[T], indexGrains+1),
	}

	p.sortContext.cond = sync.NewCond(&sync.Mutex{})

	// Pushes elements to unsorted chanel
	go p.elementProvider(unsortedChanel, unsortedHome)

	for wid := range runtime.NumCPU() {

		wg.Add(1)
		go func(threadId int) {
			defer wg.Done()
			p.placeElementWorker(unsortedChanel, compare)
		}(wid)
	}

	wg.Wait()
}

// elementProvider : pushes elements to chanel for sort workers
func (p *XList[T]) elementProvider(unsortedChanel chan *xlistObj[T], home *xlistObj[T]) {
	defer close(unsortedChanel) // Ensure channel is closed when done

	for xObj := home; xObj != nil; {
		tmp := xObj
		xObj = xObj.next
		unsortedChanel <- tmp
	}
}

// SortAddElement : add element to the list according to 'compare' closure.
//
// Parameters:
//   - unsortedChanel: Channel providing unsorted elements
//   - compare: Function to compare elements
func (p *XList[T]) placeElementWorker(unsortedChanel <-chan *xlistObj[T], compare func(a, b T) bool) {
	var xObj *xlistObj[T]
	var ptr *xlistObj[T]
	var pptr *xlistObj[T]
	var index int
	var gSize int

	gixMap := map[int]*indexPair[T]{} // temporary index store, put to main indexes on changes, strictly

	// Acquire read lock to prevent changes during reading
	p.sortContext.changeMtx.RLock()
	defer p.sortContext.changeMtx.RUnlock() // Ensure lock is released when function exits

	for xObj = range unsortedChanel {
		if p.sortContext.grains > 0 && p.Size() > p.sortContext.grains*sortIndexingFromSize {
			gSize = p.Size() / p.sortContext.grains
		}

		ptr, index = p.searchInterval(xObj, compare)
		for ; ptr != nil; ptr = ptr.next {
			// Store actual granule index with pointer (renew indexes only when locked for changes
			if gSize > 0 && index > 0 && index%gSize == 0 {
				i := int(math.Ceil(float64(index)/float64(gSize))) - 1
				gixMap[i] = &indexPair[T]{index, ptr}
			}

			if !p.sortContext.canRead { // Ожидание вставки элемента в массив
				p.waitUntilChangesDone()
			}

			if compare(*xObj.obj, *ptr.obj) {
				// Insert an element before 'ptr'
				p.waitAllReaderStopsLockChanges()
				// until waiting all readers stop, insert position may be changed by other thread,
				// so we must search it again. +/- 1-2 elements move
				if xptr, ok := p.search2InsertPos(xObj, ptr, compare); ok {
					ptr = xptr
				}

				xObj.next = ptr
				xObj.prev = ptr.prev

				if ptr.prev != nil {
					ptr.prev.next = xObj
				} else {
					// First element
					p.home = xObj
				}

				ptr.prev = xObj
				p.size++
				pptr = nil

				p.refreshIndexesFromMap(&gixMap)
				gixMap = map[int]*indexPair[T]{} // clear local cache map of indexes

				p.startReadsChangesDone()

				break // Stop search cycle and get a new element
			}

			pptr = ptr
			index++
		}

		// If we reached the end of the list, add the element to the end
		if pptr != nil {
			p.waitAllReaderStopsLockChanges()

			xObj.prev = pptr
			xObj.next = nil

			pptr.next = xObj
			p.end = xObj

			p.size++

			p.startReadsChangesDone()
		}
	}
}

// waitUntilChangesDone : waits until changes to the list are complete.
// It releases the read lock to allow changes, then waits for the signal
// that changes are done before reacquiring the read lock.
func (p *XList[T]) waitUntilChangesDone() {
	p.sortContext.changeMtx.RUnlock() // Release read lock to allow changes

	p.sortContext.cond.L.Lock()
	for !p.sortContext.canRead { // Check for false awakenings
		p.sortContext.cond.Wait() // Wait for signal
	}
	p.sortContext.cond.L.Unlock()

	p.sortContext.changeMtx.RLock() // Reacquire read lock
}

// waitAllReaderStopsLockChanges : stops reading and waits until all readers stop,
// then acquires the write lock for making changes to the list.
func (p *XList[T]) waitAllReaderStopsLockChanges() {
	p.sortContext.canRead = false     // Signal readers to stop
	p.sortContext.changeMtx.RUnlock() // Release read lock
	p.sortContext.changeMtx.Lock()    // Acquire write lock
}

// startReadsChangesDone : signals that changes are complete and reading can resume.
// It releases the write lock and broadcasts to all waiting readers.
func (p *XList[T]) startReadsChangesDone() {
	p.sortContext.canRead = true     // Signal that reading can resume
	p.sortContext.changeMtx.Unlock() // Release write lock

	p.sortContext.cond.Broadcast()  // Notify all waiting readers
	p.sortContext.changeMtx.RLock() // Reacquire read lock
}

// refreshIndexesFromMap : updates the global indexes from the local index map.
func (p *XList[T]) refreshIndexesFromMap(gixmap *map[int]*indexPair[T]) {
	for k, v := range *gixmap {
		if k > len(p.sortContext.indexes)-1 {
			p.sortContext.indexes = append(p.sortContext.indexes, v)
		} else {
			p.sortContext.indexes[k] = v
		}
	}
}

func (p *XList[T]) searchInterval(xObj *xlistObj[T], compare func(a, b T) bool) (*xlistObj[T], int) {
	// If no indexes start from the beginning
	if len(p.sortContext.indexes) == 0 {
		return p.home, 0
	}

	// Pass thru the indexes to get appropriate interval
	for i, ptr := range p.sortContext.indexes {
		// Check the validity of the pointer
		if ptr == nil {
			if i == 0 {
				return p.home, 0
			}
			// Return the last valid element
			prev := p.sortContext.indexes[i-1]
			return prev.obj, prev.ix
		}

		// If we found an appropriate interval
		if compare(*xObj.obj, *ptr.obj.obj) {
			if i == 0 {
				return p.home, 0
			}
			prev := p.sortContext.indexes[i-1]
			return prev.obj, prev.ix
		}
	}

	// If we reached the end, return the last element
	lastIndex := len(p.sortContext.indexes) - 1
	if lastElement := p.sortContext.indexes[lastIndex]; lastElement != nil {
		return lastElement.obj, lastElement.ix
	}

	return p.home, 0
}

// search2InsertPos : refines the insertion position for an element.
// This is used after acquiring the write lock to ensure the correct insertion point.
func (p *XList[T]) search2InsertPos(insObj, placeObj *xlistObj[T], compare func(a, b T) bool) (*xlistObj[T], bool) {
	if insObj == nil || placeObj == nil || insObj.obj == nil || placeObj.obj == nil {
		return nil, false
	}

	if placeObj.prev == nil && placeObj.next == nil {
		return placeObj, true
	}

	tmp := placeObj

	for tmp.prev != nil && compare(*insObj.obj, *tmp.prev.obj) {
		tmp = tmp.prev
	}

	for tmp != nil && compare(*tmp.obj, *insObj.obj) {
		if tmp.next == nil {
			break
		}
		tmp = tmp.next
	}

	return tmp, true
}

// grainsAccordingToSize returns the number of index grains based on list size.
// This determines how many index points to maintain for efficient searching.
func (p *XList[T]) grainsAccordingToSize(size int) int {
	if size < 0 {
		return 0
	}

	switch {
	case size < 100:
		return 0
	case size <= 1000:
		return 4
	case size <= 10000:
		return 10
	case size <= 100000:
		return 100
	case size <= 500000:
		return 1000
	case size <= 1000000:
		return 5000
	case size <= 50000000:
		return 7500
	case size <= 100000000:
		return 10000
	case size <= 500000000:
		return 50000
	default:
		return 100000
	}
}
