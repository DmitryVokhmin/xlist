// corefunc.go
// xlist core functions
// Created by Vokhmin D.A. 01.2025

package xlist

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"reflect"
)

// ------ Core functions ------
//TODO: add list.DoDeepCopy(), list.DoShallowCopy()

// At : returns value at specified position.
// Returns Value and Ok flag: true - value is valid, false - no value
func (p *XList[T]) At(index int) (T, bool) {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	return p.at(index)
}

// At : returns value at specified position.
// Returns Value and Ok flag: true - value is valid, false - no value
func (p *XList[T]) at(index int) (T, bool) {
	lobj := p.goToPosition(index)
	if lobj == nil {
		var zero T
		return zero, false
	}

	return *lobj.obj, true
}

func (p *XList[T]) AtPtr(index int) T {
	var xobj T
	xobj, _ = p.At(index)

	return xobj
}

// IsEmpty : returns 'true' if container is empty
func (p *XList[T]) IsEmpty() bool {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	if p.home == nil && p.end == nil {
		return true
	}

	return false
}

// isEmpty : for internal (use without mutex)
func (p *XList[T]) isEmpty() bool {
	if p.home == nil && p.end == nil {
		return true
	}

	return false
}

// Size : returns size of container.
func (p *XList[T]) Size() int {
	return p.size
}

// LastObject : returns last object in container.
func (p *XList[T]) LastObject() (T, bool) {
	return p.At(p.size - 1)
}

// LastObjectPtr : returns last object pointer in container.
func (p *XList[T]) LastObjectPtr() T {
	xobj, _ := p.LastObject()
	if reflect.ValueOf(xobj).Kind() == reflect.Ptr {
		return xobj
	}

	panic(fmt.Sprintf("%v, %s", xobj, ErrIsNotAPointer.Error()))
}

// Clear : clear container.
func (p *XList[T]) Clear() {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.home = nil
	p.end = nil
	p.size = 0
}

// Set : set 'objects' to container.
// In case of empty objects receiver will be unchanged.
func (p *XList[T]) Set(objects ...T) {
	if len(objects) == 0 {
		return
	}

	p.Clear()
	p.Append(objects...)
}

// Append : appends 'objects' to container.
// In case of empty objects receiver will be unchanged.
func (p *XList[T]) Append(objects ...T) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	for _, obj := range objects {
		lobj := &xlistObj[T]{
			obj: &obj,
		}

		p.size++

		if p.isEmpty() {
			p.home = lobj
			p.end = p.home

			continue
		}

		lobj.prev = p.end
		p.end.next = lobj
		p.end = lobj

	}
}

// AppendUnique : appends element if it doesn't exist.
func (p *XList[T]) AppendUnique(objects ...T) {
	var hash [32]byte
	isObj := make(map[any]bool)

	getHash := func(obj *T) [32]byte {
		if obj == nil {
			return [32]byte{}
		}

		data, err := json.Marshal(*obj)
		if err != nil {
			panic(err)
		}

		hash = sha256.Sum256(data)

		return hash
	}

	// Create hash map
	lobj := p.home
	p.mtx.RLock()
	for lobj != nil {
		hash = getHash(lobj.obj)
		isObj[hash] = true

		lobj = lobj.next
	}
	p.mtx.RUnlock()

	if len(isObj) == 0 {
		p.Append(objects...)
		return
	}

	// Check object for uniqueness and add it
	for _, obj := range objects {
		hash = getHash(&obj)
		if _, found := isObj[hash]; !found {
			p.Append(obj)
		}
	}
}

// Contains : checks whether the set of objects (the whole set) in the list
// it returns false if any of them not in the list.
func (p *XList[T]) Contains(objects ...T) bool {
	isObj := make(map[any]bool)

	p.mtx.RLock()
	lobj := p.home
	for lobj != nil {
		isObj[*lobj.obj] = true
		lobj = lobj.next
	}
	p.mtx.RUnlock()

	if len(isObj) == 0 {
		return false
	}

	for _, obj := range objects {
		if _, found := isObj[obj]; !found {
			return false
		}
	}

	return true
}

// Insert : inserts object before the 'pos' position
// if position is out of right range, append element - no error
func (p *XList[T]) Insert(pos int, objects ...T) error {
	if pos < 0 || pos > p.size {
		return ErrInvalidIndex
	}

	// insert last element
	if p.size == pos {
		p.Append(objects...)
		return nil
	}

	p.mtx.Lock()
	defer p.mtx.Unlock()

	for _, obj := range objects {
		lobj := &xlistObj[T]{
			obj: &obj,
		}

		// insert first element, position ignores
		xobj := p.home
		if xobj == nil { // for empty
			if pos != 0 {
				return ErrInvalidIndex
			}

			p.home = lobj
			p.end = lobj

			p.size = 1
			pos++
			continue
		}

		// go to insert position
		xobj = p.goToPosition(pos)
		if xobj == nil {
			return ErrInvalidIndex
		}

		// insert new element
		lobj.next = xobj
		lobj.prev = xobj.prev

		if xobj.prev != nil {
			xobj.prev.next = lobj
		} else {
			p.home = lobj // put object at 0 pos
		}

		xobj.prev = lobj
		p.size++
		pos++ // move position for multiple insert
	}

	return nil
}

// Replace : replaces element at position 'pos' to 'obj'.
// Returns 'true' if replaced, 'false' if not
func (p *XList[T]) Replace(pos int, obj T) error {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.isEmpty() {
		return ErrElementNotFound
	}

	xobj := p.goToPosition(pos)
	if xobj == nil {
		return ErrElementNotFound
	}
	xobj.obj = &obj

	return nil
}

// ReplaceLast : replaces last element, returns 'true' if replaced, 'false' if not.
func (p *XList[T]) ReplaceLast(obj T) error {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.isEmpty() {
		return ErrElementNotFound
	}

	xobj := p.goToPosition(p.size - 1)
	if xobj == nil {
		return ErrElementNotFound
	}
	xobj.obj = &obj

	return nil
}

func (p *XList[T]) DeleteAt(pos int) (T, error) {
	var zero T

	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.isEmpty() {
		return zero, nil
	}

	if pos < 0 || pos >= p.size {
		return zero, ErrInvalidIndex
	}

	xobj := p.goToPosition(pos)

	if xobj.prev != nil {
		xobj.prev.next = xobj.next
	}

	if xobj.next != nil {
		xobj.next.prev = xobj.prev
	}

	// First element
	if p.home == xobj {
		p.home = xobj.next
	}

	// Last element
	if p.end == xobj {
		p.end = xobj.prev
	}

	p.size--

	return *xobj.obj, nil
}

func (p *XList[T]) DeleteLast() (T, error) {
	p.mtx.RLock()
	if p.end == nil {
		var zero T
		p.mtx.RUnlock()

		return zero, ErrElementNotFound
	}
	p.mtx.RUnlock()

	return p.DeleteAt(p.Size() - 1)
}

// Add : adds 'dList' to the receiver's list and returns new instance
// (!) it does copy elements.
func (p *XList[T]) Add(dList *XList[T]) *XList[T] {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	if dList.isEmpty() && p.isEmpty() {
		return &XList[T]{}
	}

	targetCp := p.Copy()

	if dList.isEmpty() {
		return targetCp
	}

	sourceCp := dList.Copy()

	if p.isEmpty() {
		return sourceCp
	}

	// Connect 2 chains
	if targetCp.end != nil {
		targetCp.end.next = sourceCp.home
	}

	if sourceCp.home != nil {
		sourceCp.home.prev = p.end
	}

	// Set the p.finish to the finish of connected chain
	targetCp.end = sourceCp.end

	targetCp.size += sourceCp.size

	return targetCp
}

// Move : move content from 'dList' to receiver at finish.
// (!) 'dList' is destroyed, it becomes empty.
// Operation is analogue Add, but no new object creates, so it more effective.
func (p *XList[T]) Move(dList *XList[T]) {
	if dList == nil || dList.isEmpty() {
		return
	}

	_ = p.MoveAtPos(p.size, dList)
}

// MoveAtPos : inserts (moves) content from 'dList' to receiver at position 'pos'.
// (!) 'dList' is destroyed, it becomes empty.
func (p *XList[T]) MoveAtPos(pos int, dList *XList[T]) error {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if dList.isEmpty() {
		return nil
	}

	if pos < 0 || pos > p.size {
		return ErrInvalidIndex
	}

	resetSrc := func(dList *XList[T]) {
		dList.home = nil
		dList.end = nil
		dList.size = 0
	}

	// Connect chain to the finish
	if pos == p.size {
		if p.end != nil {
			p.end.next = dList.home
		}
		if dList.home != nil {
			dList.home.prev = p.end
		}

		p.end = dList.end
		p.size += dList.size

		resetSrc(dList)

		return nil
	}

	// Insert chain
	xobj := p.goToPosition(pos)
	if xobj == nil {
		return ErrElementNotFound
	}

	// left side of dList
	if xobj.prev != nil {
		xobj.prev.next = dList.home
	}

	if dList.home != nil {
		dList.home.prev = xobj.prev
	}

	// right side of dList
	dList.end = xobj
	xobj.prev = dList.end

	// Reset dList
	resetSrc(dList)

	return nil
}

// Copy : returns a copy of the list.
// If shallowCopy is true, the objects are not copied, only the references.
func (p *XList[T]) Copy() *XList[T] {
	if p.isEmpty() {
		return &XList[T]{}
	}
	na, _ := p.CopyRange(0, p.size-1)
	return na
}

func (p *XList[T]) CopyRange(fromPos int, toPos int) (*XList[T], error) {
	var newList *XList[T]

	p.mtx.RLock()
	defer p.mtx.RUnlock()

	xobj := p.home

	i := 0
	for i <= toPos && xobj != nil {

		if i == fromPos {
			newList = &XList[T]{}
		}

		if newList != nil {
			if p.shallowCopy {
				// shallow copy
				newList.Append(*xobj.obj)
			} else {
				// deep copy
				nrxObj := xobj.obj
				newList.Append(*nrxObj)
			}
		}

		xobj = xobj.next
		i++
	}

	if newList == nil || i != toPos+1 {
		return nil, ErrInvalidIndex
	}

	return newList, nil
}

// Swap : swapping 2 elements in the list.
func (p *XList[T]) Swap(i, j int) error {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if i < 0 || j < 0 || i > p.size-1 || j > p.size-1 {
		return ErrInvalidIndex
	}

	if i == j {
		return nil
	}

	p.swap(i, j)

	return nil
}

// swap : swapping 2 elements in the list.
// Internal implementation
func (p *XList[T]) swap(i, j int) {
	if i == j {
		return
	}

	objI := p.goToPosition(i)
	objJ := p.goToPosition(j)

	if objI != nil && objJ != nil {
		objI.obj, objJ.obj = objJ.obj, objI.obj
	}
}
