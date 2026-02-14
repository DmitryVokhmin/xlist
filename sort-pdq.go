// sort-pdq.go
// PDQSort (Pattern-Defeating Quicksort) adapted for XList doubly-linked list.
// Based on the algorithm from Go's slices package (zsortordered.go / sort.go).
// Created and adapted by Vokhmin D.A. 02.2026

package xlist

import (
	"math/bits"
	"runtime"
	"sync"
	"sync/atomic"
)

// pdqSortedHint : hint for pdqsortList when choosing the pivot.
type pdqSortedHint int

const (
	pdqUnknownHint    pdqSortedHint = iota
	pdqIncreasingHint               // data appears to be sorted ascending
	pdqDecreasingHint               // data appears to be sorted descending
)

// pdqXorshift : xorshift PRNG for breakPatternsList.
// Paper: https://www.jstatsoft.org/article/view/v008i14/xorshift.pdf
type pdqXorshift uint64

func (r *pdqXorshift) Next() uint64 {
	*r ^= *r << 13
	*r ^= *r >> 7
	*r ^= *r << 17
	return uint64(*r)
}

func pdqNextPowerOfTwo(length int) uint {
	return 1 << bits.Len(uint(length))
}

// PDQSort sorts the list using a PDQSort algorithm adapted for a doubly-linked list.
// Swaps only the obj pointers inside existing nodes (no element allocation).
//
// Parameters:
//   - compare: A function that compares two elements.
//     Returns true when `a` should be before `b`, otherwise false.
func (p *XList[T]) PDQSort(compare func(a, b T) bool) {
	n := p.size
	if n < 2 {
		return
	}

	p.mtx.Lock()
	defer p.mtx.Unlock()

	less := func(a, b *xlistObj[T]) bool {
		return compare(*a.obj, *b.obj)
	}

	// Atomic counter limits parallel goroutines to GOMAXPROCS-1.
	// The calling goroutine counts as one worker.
	maxWorkers := int32(runtime.GOMAXPROCS(0) - 1)
	if maxWorkers < 1 {
		maxWorkers = 1
	}
	var active atomic.Int32

	var wg sync.WaitGroup
	pdqsortListP(p.home, p.end, n, bits.Len(uint(n)), less, &active, maxWorkers, &wg)
	wg.Wait()
}

// advanceT moves a node forward by steps steps.
func advanceT[T comparable](node *xlistObj[T], steps int) *xlistObj[T] {
	for i := 0; i < steps; i++ {
		node = node.next
	}
	return node
}

// swapObjs swaps the obj pointers of two nodes (O(1), no allocation).
func swapObjs[T comparable](a, b *xlistObj[T]) {
	a.obj, b.obj = b.obj, a.obj
}

// parallelThreshold is the minimum segment size to spawn a goroutine for.
// Below this threshold sequential recursion is faster due to goroutine overhead.
const parallelThreshold = 2048

// pdqsortListP sorts segment [lo..hi] of length n using PDQSort.
// Spawns goroutines for the smaller partition when the segment is large enough
// and the active goroutine count is below maxWorkers.
func pdqsortListP[T comparable](lo, hi *xlistObj[T], n, limit int, less func(a, b *xlistObj[T]) bool, active *atomic.Int32, maxWorkers int32, wg *sync.WaitGroup) {
	const maxInsertion = 12

	wasBalanced := true
	wasPartitioned := true

	for {
		if n <= maxInsertion {
			insertionSortList(lo, hi, less)
			return
		}

		if limit == 0 {
			heapSortList(lo, n, less)
			return
		}

		if !wasBalanced {
			breakPatternsList(lo, n)
			limit--
		}

		pivotOffset, hint := choosePivotList(lo, n, less)

		if hint == pdqDecreasingHint {
			reverseRangeList(lo, hi)
			pivotOffset = (n - 1) - pivotOffset
			hint = pdqIncreasingHint
		}

		if wasBalanced && wasPartitioned && hint == pdqIncreasingHint {
			if partialInsertionSortList(lo, hi, less) {
				return
			}
		}

		// Partition equal elements if lo.prev >= pivot.
		if lo.prev != nil {
			pivotNode := advanceT(lo, pivotOffset)
			if !less(lo.prev, pivotNode) {
				newLo, newN := partitionEqualList(lo, hi, n, pivotOffset, less)
				lo = newLo
				n = newN
				continue
			}
		}

		mid, midOffset, alreadyPartitioned := partitionList(lo, hi, n, pivotOffset, less)
		wasPartitioned = alreadyPartitioned

		leftLen := midOffset          // [lo, mid)
		rightLen := n - midOffset - 1 // (mid, hi]
		balanceThreshold := n / 8

		if leftLen < rightLen {
			wasBalanced = leftLen >= balanceThreshold
			if leftLen > 1 {
				spawnOrRun(lo, mid.prev, leftLen, limit, less, active, maxWorkers, wg)
			}
			lo = mid.next
			n = rightLen
		} else {
			wasBalanced = rightLen >= balanceThreshold
			if rightLen > 1 {
				spawnOrRun(mid.next, hi, rightLen, limit, less, active, maxWorkers, wg)
			}
			hi = mid.prev
			n = leftLen
		}
	}
}

// spawnOrRun runs pdqsortListP in a goroutine if the segment is large enough
// and active goroutine count is below maxWorkers; otherwise runs synchronously.
func spawnOrRun[T comparable](lo, hi *xlistObj[T], n, limit int, less func(a, b *xlistObj[T]) bool, active *atomic.Int32, maxWorkers int32, wg *sync.WaitGroup) {
	if n >= parallelThreshold && active.Load() < maxWorkers {
		if active.Add(1) <= maxWorkers {
			wg.Add(1)
			go func() {
				defer func() {
					active.Add(-1)
					wg.Done()
				}()
				pdqsortListP(lo, hi, n, limit, less, active, maxWorkers, wg)
			}()
			return
		}
		active.Add(-1) // lost the race — undo increment
	}
	pdqsortListP(lo, hi, n, limit, less, active, maxWorkers, wg)
}

// insertionSortList sorts the segment [lo, hi] using insertion sort.
func insertionSortList[T comparable](lo, hi *xlistObj[T], less func(a, b *xlistObj[T]) bool) {
	if lo == hi {
		return
	}
	for cur := lo.next; cur != nil; cur = cur.next {
		for j := cur; j != lo && less(j, j.prev); j = j.prev {
			swapObjs(j, j.prev)
		}
		if cur == hi {
			break
		}
	}
}

// heapSortList sorts [lo, n) using heapsort.
// Collects node pointers into a temporary slice for O(1) random access.
// This is the fallback for degenerate inputs — called rarely.
func heapSortList[T comparable](lo *xlistObj[T], n int, less func(a, b *xlistObj[T]) bool) {
	nodes := make([]*xlistObj[T], n)
	cur := lo
	for i := 0; i < n; i++ {
		nodes[i] = cur
		cur = cur.next
	}
	for i := (n - 1) / 2; i >= 0; i-- {
		siftDownList(nodes, i, n, less)
	}
	for i := n - 1; i >= 0; i-- {
		nodes[0].obj, nodes[i].obj = nodes[i].obj, nodes[0].obj
		siftDownList(nodes, 0, i, less)
	}
}

func siftDownList[T comparable](nodes []*xlistObj[T], root, n int, less func(a, b *xlistObj[T]) bool) {
	for {
		child := 2*root + 1
		if child >= n {
			break
		}
		if child+1 < n && less(nodes[child], nodes[child+1]) {
			child++
		}
		if !less(nodes[root], nodes[child]) {
			return
		}
		nodes[root].obj, nodes[child].obj = nodes[child].obj, nodes[root].obj
		root = child
	}
}

// partitionList partitions [lo..hi] (length n) around pivot at offset pivotOffset from lo.
// Uses two counters (iOff, jOff) to track positions — no O(n) distance calls.
// Returns: pivot node after partition, its offset from lo, alreadyPartitioned flag.
func partitionList[T comparable](lo, hi *xlistObj[T], n, pivotOffset int, less func(a, b *xlistObj[T]) bool) (*xlistObj[T], int, bool) {
	// Move pivot to lo.
	pivotNode := advanceT(lo, pivotOffset)
	swapObjs(lo, pivotNode)
	pivotNode = lo

	// i starts at lo+1 (offset 1), j starts at hi (offset n-1).
	i := lo.next
	iOff := 1
	j := hi
	jOff := n - 1

	// Advance i forward while i < pivot.
	for iOff <= jOff && less(i, pivotNode) {
		i = i.next
		iOff++
	}
	// Retreat j backward while j >= pivot.
	for jOff > 0 && !less(j, pivotNode) {
		j = j.prev
		jOff--
	}

	if iOff > jOff {
		// Already partitioned: place pivot at j.
		swapObjs(j, pivotNode)
		return j, jOff, true
	}

	swapObjs(i, j)
	i = i.next
	iOff++
	if jOff > 0 {
		j = j.prev
		jOff--
	}

	for {
		for iOff <= jOff && less(i, pivotNode) {
			i = i.next
			iOff++
		}
		for jOff > 0 && !less(j, pivotNode) {
			j = j.prev
			jOff--
		}
		if iOff > jOff {
			break
		}
		swapObjs(i, j)
		i = i.next
		iOff++
		if jOff > 0 {
			j = j.prev
			jOff--
		}
	}

	swapObjs(j, pivotNode)
	return j, jOff, false
}

// partitionEqualList partitions elements equal to pivot to the front of [lo..hi].
// Returns the first node of the "greater than pivot" section and its length.
func partitionEqualList[T comparable](lo, hi *xlistObj[T], n, pivotOffset int, less func(a, b *xlistObj[T]) bool) (*xlistObj[T], int) {
	pivotNode := advanceT(lo, pivotOffset)
	swapObjs(lo, pivotNode)
	pivotNode = lo

	i := lo.next
	iOff := 1
	j := hi
	jOff := n - 1

	for {
		for iOff <= jOff && !less(pivotNode, i) {
			i = i.next
			iOff++
		}
		for jOff > 0 && less(pivotNode, j) {
			j = j.prev
			jOff--
		}
		if iOff > jOff {
			break
		}
		swapObjs(i, j)
		i = i.next
		iOff++
		if jOff > 0 {
			j = j.prev
			jOff--
		}
	}

	// i now points to the first element greater than pivot.
	if iOff > n-1 {
		return nil, 0
	}
	return i, n - iOff
}

// partialInsertionSortList tries to sort [lo, hi] with up to maxSteps fixes.
// Returns true if the segment is fully sorted.
func partialInsertionSortList[T comparable](lo, hi *xlistObj[T], less func(a, b *xlistObj[T]) bool) bool {
	const (
		maxSteps         = 5
		shortestShifting = 50
	)

	// Count length once — O(n) but called only when hint==increasing.
	n := 0
	for c := lo; c != hi.next; c = c.next {
		n++
	}

	i := lo.next
	for step := 0; step < maxSteps; step++ {
		for i != nil && i != hi.next && !less(i, i.prev) {
			i = i.next
		}
		if i == nil || i == hi.next {
			return true
		}
		if n < shortestShifting {
			return false
		}

		swapObjs(i, i.prev)
		// Shift left.
		for j := i.prev; j != lo && less(j, j.prev); j = j.prev {
			swapObjs(j, j.prev)
		}
		// Shift right.
		for j := i.next; j != nil && j != hi.next && less(j, j.prev); j = j.next {
			swapObjs(j, j.prev)
		}
	}
	return false
}

// breakPatternsList scatters 3 elements around the center to defeat patterns.
func breakPatternsList[T comparable](lo *xlistObj[T], n int) {
	if n < 8 {
		return
	}
	random := pdqXorshift(uint64(n))
	modulus := pdqNextPowerOfTwo(n)

	center := advanceT(lo, n/4*2-1)
	cur := center
	for idx := 0; idx < 3; idx++ {
		other := int(uint(random.Next()) & (modulus - 1))
		if other >= n {
			other -= n
		}
		otherNode := advanceT(lo, other)
		swapObjs(cur, otherNode)
		if cur.next != nil && idx < 2 {
			cur = cur.next
		}
	}
}

// choosePivotList returns the offset (from lo) of the chosen pivot and a sorted hint.
// Returns offset instead of pointer to avoid a second advanceT call in the caller.
func choosePivotList[T comparable](lo *xlistObj[T], n int, less func(a, b *xlistObj[T]) bool) (int, pdqSortedHint) {
	const (
		shortestNinther = 50
		maxSwaps        = 4 * 3
	)

	swaps := 0

	iOff := n / 4 * 1
	jOff := n / 4 * 2
	kOff := n / 4 * 3

	i := advanceT(lo, iOff)
	j := advanceT(i, jOff-iOff) // reuse i to avoid re-walking from lo
	k := advanceT(j, kOff-jOff)

	if n >= 8 {
		if n >= shortestNinther {
			i, iOff = medianAdjacentList(i, iOff, &swaps, less)
			j, jOff = medianAdjacentList(j, jOff, &swaps, less)
			k, kOff = medianAdjacentList(k, kOff, &swaps, less)
		}
		j, jOff = medianList3(i, iOff, j, jOff, k, kOff, &swaps, less)
	}

	switch swaps {
	case 0:
		return jOff, pdqIncreasingHint
	case maxSwaps:
		return jOff, pdqDecreasingHint
	default:
		return jOff, pdqUnknownHint
	}
}

// order2List returns (a,aOff,b,bOff) where a <= b, counting a swap if needed.
func order2List[T comparable](a *xlistObj[T], aOff int, b *xlistObj[T], bOff int, swaps *int, less func(a, b *xlistObj[T]) bool) (*xlistObj[T], int, *xlistObj[T], int) {
	if less(b, a) {
		*swaps++
		return b, bOff, a, aOff
	}
	return a, aOff, b, bOff
}

// medianList3 returns the node and offset of the median of a, b, c.
func medianList3[T comparable](a *xlistObj[T], aOff int, b *xlistObj[T], bOff int, c *xlistObj[T], cOff int, swaps *int, less func(a, b *xlistObj[T]) bool) (*xlistObj[T], int) {
	a, aOff, b, bOff = order2List(a, aOff, b, bOff, swaps, less)
	b, bOff, c, cOff = order2List(b, bOff, c, cOff, swaps, less)
	a, aOff, b, bOff = order2List(a, aOff, b, bOff, swaps, less)
	_, _ = c, cOff
	return b, bOff
}

// medianAdjacentList finds the median of node.prev, node, node.next.
// Returns the median node and its offset.
func medianAdjacentList[T comparable](node *xlistObj[T], off int, swaps *int, less func(a, b *xlistObj[T]) bool) (*xlistObj[T], int) {
	return medianList3(node.prev, off-1, node, off, node.next, off+1, swaps, less)
}

// reverseRangeList reverses the segment [lo, hi] by swapping obj pointers.
func reverseRangeList[T comparable](lo, hi *xlistObj[T]) {
	i, j := lo, hi
	for i != j && i.prev != j {
		swapObjs(i, j)
		i = i.next
		j = j.prev
	}
}
