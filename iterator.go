// iterator.go
// Iterator for sequential element processing through the list
// Created by Vokhmin D.A. 03.2025

package xlist

// Iterator : creates an iterator for sequential processing.
//
// 'workRange' - the first element of the range is the starting index,
// the second element is the end index, if only one element is passed,
// the range of work is from the first element to the last.
func (p *XList[T]) Iterator(workRange ...int) *Iterator[T] {

	if len(workRange) == 1 {
		return &Iterator[T]{parent: p, index: workRange[0], start: workRange[0], finish: -1}
	}

	if len(workRange) >= 2 {
		return &Iterator[T]{parent: p, index: workRange[0], start: workRange[0], finish: workRange[1]}
	}

	return &Iterator[T]{parent: p, index: -1, start: -1, finish: -1}
}

// Reset - resets the iterator with a new range of work.
// If empty, the iterator is reset to pass from the first to the last of the container elements.
func (p *Iterator[T]) Reset(workRange ...int) {
	p.lobj = nil
	p.index = -1
	p.start = -1
	p.finish = -1

	switch len(workRange) {
	case 1:
		p.start = workRange[0]
		return
	case 2:
		p.start = workRange[0]
		p.finish = workRange[1]
		return
	default:
		return
	}
}

func (p *Iterator[T]) setInitial() {
	if p.start == -1 {
		p.start = 0
	}

	if p.finish == -1 {
		p.finish = p.parent.Size() - 1
	}
}
func (p *Iterator[T]) setInitialForward() {
	p.setInitial()

	p.lobj = p.parent.goToPosition(p.start)
	if p.lobj != nil {
		p.index = p.start
	}
}
func (p *Iterator[T]) setInitialBackward() {
	p.setInitial()

	p.lobj = p.parent.goToPosition(p.finish)
	if p.lobj != nil {
		p.index = p.finish
	}
}

// SetIndex : sets the iterator to the specified index.
func (p *Iterator[T]) SetIndex(index int) (T, bool) {
	p.setInitial()

	xObj := p.parent.goToPosition(index)
	if xObj == nil || (index < p.start || index > p.finish) || index > p.parent.Size()-1 {
		var zero T
		return zero, false
	}

	p.lobj = xObj
	p.index = index

	return *xObj.obj, true
}

// SetFirst : returns the first element of the container.
// If Iterator was initialized with range, then returns the first element of the range.
func (p *Iterator[T]) SetFirst() (T, bool) {
	var xObj *xlistObj[T]

	p.setInitial()

	if p.start == 0 {
		xObj = p.parent.home
	} else {
		xObj = p.parent.goToPosition(p.start)
	}

	if xObj == nil {
		var zero T
		return zero, false
	}

	p.lobj = xObj
	p.index = p.start

	return *xObj.obj, true
}

// SetLast : returns the last element of the container.
// If Iterator was initialized with range, then returns the last element of the range.
func (p *Iterator[T]) SetLast() (T, bool) {
	var xObj *xlistObj[T]

	p.setInitial()

	if p.finish == p.parent.Size()-1 {
		xObj = p.parent.end
	} else {
		xObj = p.parent.goToPosition(p.finish)
	}

	if xObj == nil {
		var zero T
		return zero, false
	}

	p.lobj = xObj
	p.index = p.finish

	return *xObj.obj, true
}

// Index : returns current index.
func (p *Iterator[T]) Index() int {
	if p.lobj == nil {
		return -1
	}
	return p.index
}

// Value : returns current iterator value.
func (p *Iterator[T]) Value() (T, bool) {
	if p.lobj == nil {
		var zero T
		return zero, false
	}

	return *p.lobj.obj, true
}

func (p *Iterator[T]) Next() bool {
	if p.lobj == nil {
		p.setInitialForward()

		return p.lobj != nil
	}

	if p.lobj.next == nil || p.index >= p.finish {
		return false
	}

	p.lobj = p.lobj.next
	p.index++

	return true
}

func (p *Iterator[T]) Prev() bool {
	if p.lobj == nil {
		p.setInitialBackward()

		return p.lobj != nil
	}

	if p.lobj.prev == nil || p.index <= p.start {
		return false
	}

	p.lobj = p.lobj.prev
	p.index--

	return true
}

// NextValue : returns next value from container and 'true' if value is valid,
// 'false' in case of the list end value is reached (invalid value).
func (p *Iterator[T]) NextValue() (T, bool) {
	var zero T // empty object

	if p.lobj == nil {
		p.setInitialForward()

		////////////////////////
		if p.lobj == nil {
			return zero, false
		}

		return *p.lobj.obj, p.lobj != nil
	}

	if p.lobj == nil || (p.index+1 > p.finish && p.index > 0) {
		return zero, false
	}

	if p.lobj.next != nil {
		p.lobj = p.lobj.next
		p.index++

		return *p.lobj.obj, true
	}

	return zero, false
}

// PrevValue : returns previous value from container and 'true' if value is valid,
// 'false' in case of the list begin value is reached (invalid value).
func (p *Iterator[T]) PrevValue() (T, bool) {
	var zero T

	if p.lobj == nil {
		p.lobj = p.parent.end
	}

	if p.lobj == nil || (p.index+1 < p.start && p.index > 0) {
		return zero, false
	}

	if p.lobj.prev != nil && p.index > p.start {
		p.lobj = p.lobj.prev
		p.index--

		return *p.lobj.obj, true
	}

	return zero, false
}
