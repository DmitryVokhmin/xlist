// corefunc.go
// xlist core functions
// Created by Vokhmin D.A. 01.2025

package xlist

import (
	"crypto/sha256"
	"encoding/json"
)

// ------ Core functions ------

// At : returns the value at the specified index.
// Returns the value and a bool flag: true if the value is valid, false if index is out of range.
// This method is recommended for value types (e.g., XList[int], XList[string])
// where you need to distinguish between a valid zero value and a missing element.
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

// AtPtr returns the value at the specified index, or zero value if not found.
// This method is designed for pointer types (e.g., XList[*User], XList[*MyStruct])
// where nil naturally indicates absence - simply check if the returned pointer is nil.
// For value types, use At() instead to properly distinguish zero values from missing elements.
func (p *XList[T]) AtPtr(index int) T {
	xobj, _ := p.At(index)
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

// LastObject returns the last element in the container.
// Returns the value and a bool flag: true if the value is valid, false if container is empty.
// This method is recommended for value types (e.g., XList[int], XList[string])
// where you need to distinguish between a valid zero value and an empty container.
func (p *XList[T]) LastObject() (T, bool) {
	return p.At(p.size - 1)
}

// LastObjectPtr returns the last element in the container, or zero value if container is empty.
// This method is designed for pointer types (e.g., XList[*User], XList[*MyStruct])
// where nil naturally indicates absence - simply check if the returned pointer is nil.
// For value types, use LastObject() instead to properly distinguish zero values from empty container.
func (p *XList[T]) LastObjectPtr() T {
	xobj, _ := p.LastObject()
	return xobj
}

// Clear : clear container.
func (p *XList[T]) Clear() *XList[T] {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.home = nil
	p.end = nil
	p.size = 0

	return p
}

// Set : set 'objects' to container.
// In case of empty objects receiver will be unchanged.
// Returns self for method chaining; return value can be ignored.
func (p *XList[T]) Set(objects ...T) *XList[T] {
	if len(objects) == 0 {
		return p
	}

	p.Clear()
	p.Append(objects...)

	return p
}

// Append : appends 'objects' to container.
// In case of empty objects receiver will be unchanged.
// Returns self for method chaining; return value can be ignored.
func (p *XList[T]) Append(objects ...T) *XList[T] {
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

	return p
}

// AppendUnique : appends element if it doesn't exist in current collection.
// Returns self for method chaining; return value can be ignored.
func (p *XList[T]) AppendUnique(objects ...T) *XList[T] {
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
		return p
	}

	// Check object for uniqueness and add it
	for _, obj := range objects {
		hash = getHash(&obj)
		if _, found := isObj[hash]; !found {
			p.Append(obj)
		}
	}

	return p
}

// Contains : checks whether the set of objects (the whole set) in the list
// it returns false if any of them not in the list.
func (p *XList[T]) Contains(objects ...T) bool {
	return p.containsInternal(false, objects...)
}

// ContainsSome : checks whether any of objects of `objects` is in the list
// it returns false if no one of them not in the list.
func (p *XList[T]) ContainsSome(objects ...T) bool {
	return p.containsInternal(true, objects...)
}

// containsInternal : internal realisation (for optimal Contains and ContainsSome)
func (p *XList[T]) containsInternal(containsSome bool, objects ...T) bool { // nosonar
	if len(objects) == 0 {
		if containsSome {
			return false // Search for nothing, find nothing.
		}
		return true // The empty set is a subset of every set.
	}

	if p.home == nil {
		return false
	}

	// Special case for one object
	if len(objects) == 1 {
		target := objects[0]
		p.mtx.RLock()
		defer p.mtx.RUnlock()

		lobj := p.home
		for lobj != nil {
			if *lobj.obj == target { // direct compare T
				return true
			}
			lobj = lobj.next
		}
		return false
	}

	lookingFor := make(map[T]struct{}, len(objects))
	for _, obj := range objects {
		lookingFor[obj] = struct{}{}
	}

	p.mtx.RLock()
	defer p.mtx.RUnlock()

	xobj := p.home
	for xobj != nil {
		if _, found := lookingFor[*xobj.obj]; found {
			if containsSome {
				return true
			}

			delete(lookingFor, *xobj.obj)
			if len(lookingFor) == 0 {
				return true // all objects found
			}
		}
		xobj = xobj.next
	}

	return false
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

// DeleteAt : deletes and returns the element at the specified position, or an error if the position is invalid.
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
	if xobj == nil {
		var zero T
		return zero, ErrElementNotFound
	}

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

// AppendList  adds objects to the end of the list (mutating).
// Returns self for method chaining; return value can be ignored.
// (-) Add
func (p *XList[T]) AppendList(dList *XList[T]) *XList[T] {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	if dList.isEmpty() && p.isEmpty() {
		return &XList[T]{}
	}

	targetCp := p

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
		sourceCp.home.prev = targetCp.end
	}

	// Set the targetCp.end to the tail of connected chain
	targetCp.end = sourceCp.end

	targetCp.size += sourceCp.size

	return targetCp
}

// Splice : move content from 'dList' to receiver at its tail (appends container - mutating).
// (!) 'dList' is destroyed, it becomes empty.
// (-) Move
func (p *XList[T]) Splice(dList *XList[T]) *XList[T] {

	if dList == nil || dList.isEmpty() {
		return p
	}

	_ = p.SpliceAtPos(p.size, dList)

	return p
}

// SpliceAtPos : inserts (moves) content from 'dList' to receiver at position 'pos'.
// (!) 'dList' is destroyed, it becomes empty.
// (-) MoveAtPos
func (p *XList[T]) SpliceAtPos(pos int, dList *XList[T]) error {
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

	// In case of empty receiver
	if p.isEmpty() {
		if pos != 0 {
			return ErrInvalidIndex
		}
		p.home = dList.home
		p.end = dList.end
		p.size = dList.size
		resetSrc(dList)

		return nil
	}

	// Connect chain to the tail
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
	if xobj == nil { // (!!!) Можно не проверять, так как предыдущие проверки гарантируют не nil
		return ErrElementNotFound
	}

	// left side of dList
	if xobj.prev != nil {
		xobj.prev.next = dList.home
	} else {
		p.home = dList.home
	}

	if dList.home != nil {
		dList.home.prev = xobj.prev
	}

	// right side of dList
	//dList.end = xobj
	dList.end.next = xobj
	xobj.prev = dList.end

	// Reset dList
	resetSrc(dList)

	return nil
}

// Copy : returns a copy of the list.
// It makes shallow copies of objects, so be careful when changing container objects.
// Consider 'DeepCopy' method to copy the container objects themselves.
func (p *XList[T]) Copy() *XList[T] {
	if p.isEmpty() {
		return &XList[T]{}
	}
	na, _ := p.CopyRange(0, p.size-1)
	return na
}

// CopyRange : returns a new container with elements of receiver for specified range [fromPos, toPos].
// It makes shallow copies of objects, so be careful when changing container objects .
func (p *XList[T]) CopyRange(fromPos int, toPos int) (*XList[T], error) {
	return p.DeepCopyRange(fromPos, toPos, func(obj T) T {
		return obj
	})
}

// DeepCopy :  returns a new container with new elements of receiver.
// It makes deep copies of objects, so you must provide a closure 'deepCopyFn' to make a deep copy of type T.
func (p *XList[T]) DeepCopy(deepCopyFn func(T) T) *XList[T] {
	if p.isEmpty() || deepCopyFn == nil {
		return &XList[T]{}
	}

	na, _ := p.DeepCopyRange(0, p.size-1, deepCopyFn)
	return na
}

// DeepCopyRange : returns a new container with elements from the range [fromPos, toPos].
// You must provide a closure 'deepCopyFn' that knows how to make a deep copy of type T.
func (p *XList[T]) DeepCopyRange(fromPos int, toPos int, deepCopyFn func(T) T) (*XList[T], error) {

	if deepCopyFn == nil {
		return nil, ErrNoClosure
	}

	p.mtx.RLock()
	defer p.mtx.RUnlock()

	if fromPos < 0 || fromPos > p.size-1 || toPos < 0 || toPos > p.size-1 || fromPos > toPos {
		return nil, ErrInvalidIndex
	}

	// toPos is required for speculative iteration to get CPU cache
	xobjs := p.getObjectsAt(fromPos, toPos)
	if len(xobjs) == 0 || xobjs[0] == nil {
		return nil, ErrElementNotFound
	}

	result := &XList[T]{}
	xobj := xobjs[0]
	i := fromPos

	for xobj != nil {
		result.Append(deepCopyFn(*xobj.obj))

		if i == toPos {
			break
		}

		xobj = xobj.next
		i++
	}
	return result, nil
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
