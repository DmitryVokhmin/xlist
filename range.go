// range.go
// Provide `range` iterator for sequential element processing through the list
// Created by Vokhmin D.A. 02.2026

package xlist

import "iter"

// ForwardPtr returns an iterator for forward traversal of the list with index and pointer.
// Returns all elements including those with nil pointers.
// Use this method when you need element index or access to nil elements.
//
// Example:
//
//	for i, ptr := range list.ForwardPtr() {
//	    if ptr == nil {
//	        fmt.Printf("[%d] = nil\n", i)
//	    } else {
//	        fmt.Printf("[%d] = %v\n", i, *ptr)
//	    }
//	}
func (p *XList[T]) ForwardPtr() iter.Seq2[int, *T] {
	return func(yield func(int, *T) bool) {
		p.mtx.RLock()
		defer p.mtx.RUnlock()

		tmp := p.home
		index := 0

		for tmp != nil {
			if !yield(index, tmp.obj) {
				return
			}
			tmp = tmp.next
			index++
		}
	}
}

// BackwardPtr returns an iterator for backward traversal of the list with index and pointer.
// Returns all elements including those with nil pointers.
// Indexes start from (size - 1) and decrement to 0.
// Use this method when you need element index or access to nil elements in reverse order.
//
// Example:
//
//	for i, ptr := range list.BackwardPtr() {
//	    if ptr != nil {
//	        fmt.Printf("[%d] = %v\n", i, *ptr)
//	    }
//	}
func (p *XList[T]) BackwardPtr() iter.Seq2[int, *T] {
	return func(yield func(int, *T) bool) {
		p.mtx.RLock()
		defer p.mtx.RUnlock()

		tmp := p.end
		index := p.size - 1

		for tmp != nil {
			if !yield(index, tmp.obj) {
				return
			}
			tmp = tmp.prev
			index--
		}
	}
}

// Values returns an iterator over non-nil values without index.
// Elements with nil pointers are automatically skipped (!).
// Use this method for simple iteration when you don't need element index.
//
// Example:
//
//	for v := range list.Values() {
//	    process(v) // v is T, not *T
//	}
func (p *XList[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		p.mtx.RLock()
		defer p.mtx.RUnlock()

		tmp := p.home
		for tmp != nil {
			if tmp.obj != nil {
				if !yield(*tmp.obj) {
					return
				}
			}
			tmp = tmp.next
		}
	}
}

// ValuesRev returns a reverse iterator over non-nil values without index.
// Elements with nil pointers are automatically skipped ( ! ).
// Traversal starts from the last element and moves towards the first.
// Use this method for simple backward iteration when you don't need element index.
//
// Example:
//
//	for v := range list.ValuesRev() {
//	    process(v) // v is T, not *T, in reverse order
//	}
func (p *XList[T]) ValuesRev() iter.Seq[T] {
	return func(yield func(T) bool) {
		p.mtx.RLock()
		defer p.mtx.RUnlock()

		tmp := p.end
		for tmp != nil {
			if tmp.obj != nil {
				if !yield(*tmp.obj) {
					return
				}
			}
			tmp = tmp.prev
		}
	}
}
