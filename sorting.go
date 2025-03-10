// sorting.go
// Sorts the items in the container
// Created by Vokhmin D.A. 03.2025

package xlist

import (
	"math"
	"runtime"
	"sync"
)

// Sort : sort list according to 'compare' closure
//
//	compare - is a closure that compares 2 elements.
//	true - when `a` must be before `b`, otherwise `false`
func (p *XList[T]) Sort(compare func(a, b T) bool) {

	// One thread
	// 0.18 sec = 10 000
	// ~28.00 sec = 100 000
	// two threads
	// 0.15 sec = 10 000
	//~28.00 sec = 100 000
	//
	// 16 threads
	// 0.06 sec = 10 000
	// 3.53 sec = 100 000
	// 12 - 16 sec = 200 000
	//
	// 16 threads (with indexes)
	// 0.05 sec = 10 000
	// 1.04 - 1.47 sec = 100 000
	// 2.67 - 3 sec = 200 000
	// 44.69-55.46 sec = 1 000 000
	p.ScanSort(compare)
}

// ScanSort : sort list with simple move and order algoritm.
//
//	compare - is a closure that compares 2 elements.
//	true - when `a` must be before `b`, otherwise `false`
func (p *XList[T]) ScanSort(compare func(a, b T) bool) {
	wg := &sync.WaitGroup{}

	if p.Size() < 2 {
		return
	}

	// Close component work until sort is done
	p.mtx.Lock()
	defer p.mtx.Unlock()

	// Канал подачи элементов на сортировку
	unsortedChanel := make(chan *xlistObj[T])

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
			p.placeElementWorker(unsortedChanel, compare, threadId)
			wg.Done()
		}(wid)
	}

	wg.Wait()

}

// elementProvider : pushes elements to chanel for sort workers
func (p *XList[T]) elementProvider(unsortedChanel chan *xlistObj[T], home *xlistObj[T]) {
	xObj := home
	for xObj != nil {

		tmp := xObj
		xObj = xObj.next

		unsortedChanel <- tmp
	}

	close(unsortedChanel)
}

// SortAddElement : add element to the list according to 'compare' closure.
//
//	compare - is a closure that compares 2 elements.
//	true - when `a` must be before `b`, otherwise `false`
func (p *XList[T]) placeElementWorker(unsortedChanel <-chan *xlistObj[T], compare func(a, b T) bool, wid int) {
	var xObj *xlistObj[T]
	var ptr *xlistObj[T]
	var pptr *xlistObj[T]
	var index int
	var gSize int

	gixMap := map[int]*indexPair[T]{} // temporary index store, put to main indexes on changes, strictly

	p.sortContext.changeMtx.RLock() // read lock to prevent changes

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
				// Вставляем элемент перед 'ptr'
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
					// first element
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

		// Если цикл закончился, но ничего не вставлено, то добавляем элемент в конец
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

	p.sortContext.changeMtx.RUnlock()
}

// waitUntilChangesDone : wait until changes done
func (p *XList[T]) waitUntilChangesDone() {
	p.sortContext.changeMtx.RUnlock() // read unlock to allow change

	p.sortContext.cond.L.Lock()
	for !p.sortContext.canRead { // checking for false awakenings
		p.sortContext.cond.Wait() // wait for signal
	}
	p.sortContext.cond.L.Unlock()

	p.sortContext.changeMtx.RLock() // read lock to prevent changes
}

// stopReadsAnsWaitChanges : stop reading wait until all of the reader stops ( p.sortFlow.changeMtx.RUnlock() )
// then wait until changes done
func (p *XList[T]) waitAllReaderStopsLockChanges() {
	p.sortContext.canRead = false     // Flag to enter at stop reading block
	p.sortContext.changeMtx.RUnlock() // read unlock in changes thread
	p.sortContext.changeMtx.Lock()    // Lock for changing and reading
}

// startReadsChangesDone : start reading all changes done
func (p *XList[T]) startReadsChangesDone() {
	p.sortContext.canRead = true     // Flag to do not enter at "stop reading block"
	p.sortContext.changeMtx.Unlock() // Unlock for changing and reading

	p.sortContext.cond.Broadcast() // Broadcast start to all readers
	p.sortContext.changeMtx.RLock()
}

func (p *XList[T]) refreshIndexesFromMap(gixmap *map[int]*indexPair[T]) {
	for k, v := range *gixmap {
		if k > len(p.sortContext.indexes)-1 {
			p.sortContext.indexes = append(p.sortContext.indexes, v)
		} else {
			p.sortContext.indexes[k] = v
		}
	}
}

func (p *XList[T]) searchInterval(xObj *xlistObj[T], compare func(a, b T) bool) (*xlistObj[T], int) { // search interval
	lastElementIndex := 0

	// Ищем интервал элементов в котором будем производить поиск
	if len(p.sortContext.indexes) == 0 {
		return p.home, 0 // if no indexes, scan all elements from the beginning
	} else {
		for i, ptr := range p.sortContext.indexes {
			if ptr == nil {
				break
			}

			if compare(*xObj.obj, *ptr.obj.obj) {
				if i == 0 {
					return p.home, 0
				}
				// index := i * p.Size() / p.sortContext.grains
				ipe := p.sortContext.indexes[i-1]
				return ipe.obj, ipe.ix
			}
			lastElementIndex = i
		}
	}

	// Nothing found, return last element
	if lastElementIndex == 0 {
		return p.home, 0
	}

	if le := p.sortContext.indexes[lastElementIndex]; le != nil {
		return le.obj, le.ix
	}

	return p.home, 0
}

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
