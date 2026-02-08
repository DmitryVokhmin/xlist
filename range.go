// range.go
// Provide `range` iterator for sequential element processing through the list
// Created by Vokhmin D.A. 02.2026

package xlist

import (
	"fmt"
	"iter"
)

type direction bool

type RangeOptions struct {
	index int
	count int
}

func WithPos(pos int) func(*RangeOptions) {
	return func(io *RangeOptions) {
		io.index = pos
	}
}

func WithCount(count int) func(*RangeOptions) {
	return func(io *RangeOptions) {
		io.count = count
	}
}

// ----------------------------------------------------------

// All returns a forward iterator over the list.
// Each element is yielded as (index, T).
// Options: WithPos/WithCount can limit the range.
//
// Example:
//
//	for i, obj := range list.All(xlist.WithPos(2), xlist.WithCount(3)) {
//		fmt.Println(i, obj)
//	}
func (p *XList[T]) All(opt ...func(*RangeOptions)) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		var tmp *xlistObj[T]
		var xobjList []*xlistObj[T]

		params := &RangeOptions{}

		p.mtx.RLock()
		defer p.mtx.RUnlock()

		tmp = p.home
		if len(opt) > 0 {
			for _, optSet := range opt {
				optSet(params)
			}

			if params.index < 0 || params.index >= p.size {
				panic(fmt.Sprintf("%v: index=%d, xlist size size=%d", ErrInvalidIndex, params.index, p.size))
			}
			// avoid negative indexes
			if params.count < 0 {
				params.count = -1 * params.count
			}

			// avoid over range
			if (params.index + params.count) > p.size-1 {
				params.count = p.size - params.index
			}

			// second param is speculative gos thru (to get CPU cache)
			xobjList = p.getObjectsAt(params.index, params.index+params.count)
			if len(xobjList) > 0 {
				tmp = xobjList[0]
			}
		}

		index := params.index
		count := 0
		for tmp != nil && (count < params.count || params.count == 0) {
			if !yield(index, *tmp.obj) {
				return
			}

			tmp = tmp.next
			index++
			count++
		}

	}
}

// Backward returns a reverse iterator over the list (from end to start).
// Each element is yielded as (index, T).
// Options: WithPos/WithCount can limit the range.
//
// Example:
//
//	for i, obj := range list.Backward(xlist.WithPos(2), xlist.WithCount(2)) {
//		fmt.Println(i, obj)
//	}
func (p *XList[T]) Backward(opt ...func(*RangeOptions)) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		var tmp *xlistObj[T]
		var xobjList []*xlistObj[T]

		params := &RangeOptions{index: -1}

		p.mtx.RLock()
		defer p.mtx.RUnlock()

		tmp = p.end
		index := p.size - 1
		count := p.size

		if len(opt) > 0 {
			for _, optSet := range opt {
				optSet(params)
			}

			if params.index == -1 { // if no index defined
				params.index = p.size - 1
			}

			// Validate index bounds (after -1 handling)
			if params.index < 0 || params.index >= p.size {
				panic(fmt.Sprintf("%v: index=%d, xlist size size=%d", ErrInvalidIndex, params.index, p.size))
			}

			// avoid negative indexes
			if params.count < 0 {
				params.count = -1 * params.count
			}

			// avoid over range with Count
			if (params.index + 1 - params.count) < 0 {
				params.count = params.index + 1 // ??? Проверить +/- 1
			}

			xobjList = p.getObjectsAt(params.index)
			if len(xobjList) > 0 {
				tmp = xobjList[0]
				index = params.index
				count = params.count
			}
		}

		for tmp != nil && (count > 0 /*count < params.count*/ || params.count == 0) {
			if !yield(index, *tmp.obj) {
				return
			}

			tmp = tmp.prev
			index--
			count--
		}

	}
}

// Values returns a forward iterator of values only (without indices).
// Equivalent to ToValues(list.All(...)).
//
// Example:
//
//	for v := range list.Values() {
//		fmt.Println(v)
//	}
func (p *XList[T]) Values(opt ...func(*RangeOptions)) iter.Seq[T] {
	return ToValues(p.All(opt...))
}

// ValuesBackward returns a reverse iterator of values only (without indices).
// Equivalent to ToValues(list.Backward(...)).
//
// Example:
//
//	for v := range list.ValuesBackward() {
//		fmt.Println(v)
//	}
func (p *XList[T]) ValuesBackward(opt ...func(*RangeOptions)) iter.Seq[T] {
	return ToValues(p.Backward(opt...))
}

// ToValues transforms a Seq2 iterator (with index and value) to a Seq (values only, without indices).
//
// Example:
//
//	for v := range ToValues(list.All()) {
//		fmt.Println(v)
//	}
func ToValues[T any](seq2 iter.Seq2[int, T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, value := range seq2 {
			if !yield(value) {
				return
			}
		}
	}
}

// Filter passes through elements for which the predicate returns true.
//
// Example:
//
//	filtered := Filter(list.All(), func(i int, obj *teststruct) bool {
//		return obj != nil && i%2 == 0
//	})
//
//	for _, obj := range filtered {
//		fmt.Println(obj)
//	}
func Filter[T any](it2 iter.Seq2[int, T], is_ok func(int, T) bool) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, value := range it2 {
			if is_ok(i, value) {
				if !yield(i, value) {
					return
				}
			}
		}
	}
}

// TakeWhile yields elements while the predicate remains true.
// Traversal stops when the first element fails the predicate.
//
// Example:
//
//	list := New[*int]()
//	list.Append(&val1, &val2, &val3)
//	prefix := TakeWhile(list.All(), func(_ int, val *int) bool {
//		return val != nil && *val < 10
//	})
//
//	for _, val := range prefix {
//		if val != nil {
//			fmt.Println(*val)
//		}
//	}
func TakeWhile[T any](it2 iter.Seq2[int, T], is_ok func(int, T) bool) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, value := range it2 {
			if !is_ok(i, value) {
				return
			}
			if !yield(i, value) {
				return
			}
		}
	}
}

// SkipWhile skips elements while the predicate is true, then yields the rest.
//
// Example:
//
//	list := New[*int]()
//	list.Append(&val1, &val2, &val3)
//	tail := SkipWhile(list.All(), func(_ int, val *int) bool {
//		return val != nil && *val < 10
//	})
//
//	for _, val := range tail {
//		if val != nil {
//			fmt.Println(*val)
//		}
//	}
func SkipWhile[T any](it2 iter.Seq2[int, T], is_ok func(int, T) bool) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		skipping := true
		for i, value := range it2 {
			if skipping {
				if is_ok(i, value) {
					continue
				}
				skipping = false
			}
			if !yield(i, value) {
				return
			}
		}
	}
}

// Map transforms a sequence of T into a sequence of V.
// Used to transform data extracted via Values().
//
// Example:
//
//	asStrings := Map(list.Values(), func(v int) string {
//		return fmt.Sprintf("value=%d", v)
//	})
//
//	for s := range asStrings {
//		fmt.Println(s)
//	}
func Map[T, V any](seq iter.Seq[T], transform func(T) V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if !yield(transform(v)) {
				return
			}
		}
	}
}

// AnyMatch is a terminal helper that checks whether any element matches the predicate.
// Example:
//
//	list := New[*int]()
//	list.Append(&val1, &val2, &val3)
//	hasEven := AnyMatch(list.All(), func(_ int, val *int) bool {
//		return val != nil && *val%2 == 0
//	})
//
//	if hasEven {
//	     fmt.Println("there is even")
//	}
func AnyMatch[T any](it2 iter.Seq2[int, T], is_ok func(int, T) bool) bool {
	for i, value := range it2 {
		if is_ok(i, value) {
			return true
		}
	}
	return false
}

// AllMatch is a terminal helper that checks whether all elements match the predicate.
// Example:
//
//	list := New[*int]()
//	list.Append(&val1, &val2, &val3)
//	allPositive := AllMatch(list.All(), func(_ int, val *int) bool {
//		return val != nil && *val > 0
//	})
//
//	if allPositive {
//	     fmt.Println("All are positive")
//	}
func AllMatch[T any](it2 iter.Seq2[int, T], is_ok func(int, T) bool) bool {
	for i, value := range it2 {
		if !is_ok(i, value) {
			return false
		}
	}
	return true
}
