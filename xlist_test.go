package xlist

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type teststruct struct {
	Num int
	Str string
}

func TestXList(t *testing.T) {
	var obj1, obj2, obj3, obj4, obj5, obj5cp, objRes *teststruct
	var objStk, xobjStk teststruct
	var xobj *teststruct
	var ok bool
	var err error

	obj1 = &teststruct{1, "obj1"}
	obj2 = &teststruct{2, "obj2"}
	obj3 = &teststruct{3, "obj3"}
	obj4 = &teststruct{4, "obj4"}
	obj5 = &teststruct{5, "obj5"}
	obj5cp = &teststruct{5, "obj5"}
	objStk = teststruct{10, "Stack"}

	// Test struct ----
	// Stack object
	listStk := New[teststruct]()
	listStk.Append(objStk)
	assert.Equal(t, false, listStk.IsEmpty())
	assert.Equal(t, 1, listStk.Size())
	xobjStk, ok = listStk.At(0)
	assert.Equal(t, true, ok)
	assert.Equal(t, objStk, xobjStk)

	xobjStk, ok = listStk.At(1)
	assert.Equal(t, false, ok)
	assert.Equal(t, teststruct{0, ""}, xobjStk) // Empty object

	//
	list := New[*teststruct]()

	// IsEmpty
	assert.Equal(t, list.IsEmpty(), true)

	list2 := New[*teststruct](obj1)
	assert.Equal(t, false, list2.IsEmpty())
	assert.Equal(t, 1, list2.Size())

	list2 = New[*teststruct](obj1, obj2)
	assert.Equal(t, 2, list2.Size())
	xobj, ok = list2.At(0)
	assert.Equal(t, obj1, xobj)
	assert.Equal(t, true, ok)

	xobj, ok = list2.At(1)
	assert.Equal(t, obj2, xobj)
	assert.Equal(t, true, ok)

	xobj, ok = list2.At(2)
	assert.Nil(t, xobj)
	assert.Equal(t, false, ok)

	list2 = New[*teststruct](obj1, obj2, obj3)
	assert.Equal(t, 3, list2.Size())

	xobj, _ = list2.At(0)
	assert.Equal(t, obj1, xobj)

	xobj, _ = list2.At(1)
	assert.Equal(t, obj2, xobj)

	xobj, _ = list2.At(2)
	assert.Equal(t, obj3, xobj)

	// Append, At
	list.Append(obj1)
	objRes, ok = list.At(0)
	assert.Equal(t, true, ok)
	assert.Equal(t, obj1, objRes)

	// check internal stuff
	assert.Equal(t, list.home, list.end)
	assert.Nil(t, list.home.next)
	assert.Nil(t, list.home.prev)

	// size
	assert.Equal(t, list.size, 1)

	// Clear
	assert.Equal(t, false, list.IsEmpty())
	list.Clear()
	assert.Equal(t, true, list.IsEmpty())

	// Size
	assert.Equal(t, 0, list.Size())

	// Append 2 obj
	list.Append(obj1, obj2)
	assert.Equal(t, 2, list.Size())

	// check internal stuff
	assert.Equal(t, obj1, *list.home.obj)
	assert.Equal(t, obj2, *list.end.obj)
	assert.Equal(t, obj2, *list.home.next.obj)
	assert.Equal(t, obj1, *list.end.prev.obj)
	assert.Nil(t, list.home.prev)
	assert.Nil(t, list.end.next)

	// Append
	list.Append(obj3)
	assert.Equal(t, 3, list.Size())
	assert.Equal(t, obj1, *list.home.obj)
	assert.Equal(t, obj3, *list.end.obj)
	assert.Equal(t, obj2, *list.end.prev.obj)
	assert.Equal(t, obj3, *list.end.prev.next.obj)
	assert.Nil(t, list.end.next)

	// Append
	list.Append(obj4, obj5)
	assert.Equal(t, obj1, *list.home.obj)
	assert.Equal(t, obj5, *list.end.obj)
	assert.Equal(t, obj4, *list.end.prev.obj)
	assert.Equal(t, obj3, *list.end.prev.prev.obj)
	assert.Equal(t, obj4, *list.end.prev.prev.next.obj)
	assert.Equal(t, obj5, *list.end.prev.prev.next.next.obj)
	assert.Nil(t, list.end.next)
	assert.Equal(t, 5, list.Size())

	// Check At()
	xobj, _ = list.At(0)
	assert.Equal(t, obj1, xobj)
	xobj, _ = list.At(1)
	assert.Equal(t, obj2, xobj)
	xobj, _ = list.At(2)
	assert.Equal(t, obj3, xobj)
	xobj, _ = list.At(3)
	assert.Equal(t, obj4, xobj)
	xobj, _ = list.At(4)
	assert.Equal(t, obj5, xobj)

	// Check GetLast()
	xobj, ok = list.LastObject()
	assert.Equal(t, true, ok)
	assert.Equal(t, obj5, xobj)

	// Clear
	list.Clear()
	assert.Equal(t, true, list.IsEmpty())
	assert.Equal(t, 0, list.Size())

	// Set
	list.Set(obj1)
	assert.Equal(t, 1, list.Size())
	xobj, ok = list.At(0)
	assert.Equal(t, true, ok)
	assert.Equal(t, obj1, xobj)

	xobj, ok = list.At(1)
	assert.Equal(t, false, ok)
	assert.Nil(t, xobj)

	list.Set(obj1, obj2)
	assert.Equal(t, 2, list.Size())
	xobj, ok = list.At(0)
	assert.Equal(t, true, ok)
	assert.Equal(t, obj1, xobj)
	xobj, ok = list.At(1)
	assert.Equal(t, true, ok)
	assert.Equal(t, obj2, xobj)
	xobj, ok = list.At(2) //Пустая структура
	assert.Equal(t, false, ok)
	assert.Nil(t, xobj)

	list.Set(obj5, obj2, obj3, obj4, obj1)
	assert.Equal(t, 5, list.Size())

	xobj, ok = list.At(0)
	assert.Equal(t, true, ok)
	assert.Equal(t, obj5, xobj)
	xobj, ok = list.At(1)
	assert.Equal(t, true, ok)
	assert.Equal(t, obj2, xobj)
	xobj, ok = list.At(2)
	assert.Equal(t, true, ok)
	assert.Equal(t, obj3, xobj)
	xobj, ok = list.At(3)
	assert.Equal(t, true, ok)
	assert.Equal(t, obj4, xobj)
	xobj, ok = list.At(4)
	assert.Equal(t, true, ok)
	assert.Equal(t, obj1, xobj)

	xobj, ok = list.At(5)
	assert.Equal(t, false, ok)
	assert.Nil(t, xobj)

	list.Clear()

	// Insert
	err = list.Insert(5, obj2)
	assert.ErrorIs(t, err, ErrInvalidIndex)
	err = list.Insert(0, obj2)
	assert.Nil(t, err)
	assert.Equal(t, 1, list.Size())

	// check internal stuff
	assert.Equal(t, list.home, list.end)
	assert.Nil(t, list.home.next)
	assert.Nil(t, list.home.prev)

	// Insert
	err = list.Insert(0, obj1)
	assert.Nil(t, err)
	assert.Equal(t, 2, list.Size())

	// check internal suff
	assert.Equal(t, obj1, *list.home.obj)
	assert.Equal(t, obj2, *list.end.obj)
	assert.Equal(t, obj2, *list.home.next.obj)
	assert.Equal(t, obj1, *list.end.prev.obj)
	assert.Nil(t, list.home.prev)
	assert.Nil(t, list.end.next)

	err = list.Insert(2, obj3)
	assert.Nil(t, err)
	assert.Equal(t, 3, list.Size())
	assert.Equal(t, obj1, *list.home.obj)
	assert.Equal(t, obj3, *list.end.obj)
	assert.Equal(t, obj2, *list.end.prev.obj)
	assert.Equal(t, obj3, *list.end.prev.next.obj)
	assert.Nil(t, list.end.next)

	err = list.Insert(1, obj4)
	assert.Nil(t, err)
	assert.Equal(t, 4, list.Size())
	xobj, ok = list.At(0)
	assert.Equal(t, true, ok)
	assert.Equal(t, obj1, xobj)
	xobj, ok = list.At(1)
	assert.Equal(t, true, ok)
	assert.Equal(t, obj4, xobj)
	xobj, ok = list.At(2)
	assert.Equal(t, true, ok)
	assert.Equal(t, obj2, xobj)

	// invalid insertion
	err = list.Insert(5, obj5)
	assert.ErrorIs(t, err, ErrInvalidIndex)
	assert.Equal(t, 4, list.Size())

	// Insert before the finish
	err = list.Insert(3, obj5) // moved current object to left
	assert.Nil(t, err)
	assert.Equal(t, obj2, list.AtPtr(2))
	assert.Equal(t, obj5, list.AtPtr(3))
	assert.Equal(t, obj3, list.AtPtr(4))

	// check internal stuff
	assert.Equal(t, obj1, *list.home.obj)
	assert.Equal(t, obj3, *list.end.obj)
	assert.Equal(t, obj5, *list.end.prev.obj)
	assert.Equal(t, obj2, *list.end.prev.prev.obj)
	assert.Equal(t, obj4, *list.end.prev.prev.prev.obj)
	assert.Equal(t, obj1, *list.end.prev.prev.prev.prev.obj)
	assert.Equal(t, obj4, *list.home.next.obj)
	assert.Equal(t, obj2, *list.home.next.next.obj)
	assert.Equal(t, obj5, *list.home.next.next.next.obj)
	assert.Equal(t, obj3, *list.home.next.next.next.next.obj)
	assert.Nil(t, list.end.next)
	assert.Equal(t, 5, list.Size())

	// Invalid index insert
	list.Clear()
	err = list.Insert(1, obj1, obj2, obj3, obj4, obj5)
	assert.ErrorIs(t, err, ErrInvalidIndex)

	// Multiple insert at 0
	err = list.Insert(0, obj1, obj2, obj5)
	assert.Nil(t, err)
	assert.Equal(t, 3, list.Size())

	// Multiple insert in the middle
	err = list.Insert(2, obj3, obj4)
	assert.Nil(t, err)
	assert.Equal(t, obj3, list.AtPtr(2))
	assert.Equal(t, obj4, list.AtPtr(3))

	// Multiple Insert at the finish
	err = list.Insert(5, obj1, obj2)
	assert.Nil(t, err)
	assert.Equal(t, 7, list.Size())
	assert.Equal(t, obj2, list.AtPtr(6))
	assert.Equal(t, obj2, list.LastObjectPtr())
	assert.Equal(t, obj1, list.AtPtr(5))

	assert.Equal(t, obj2, *list.end.obj)
	assert.Equal(t, obj1, *list.end.prev.obj)
	assert.Equal(t, obj5, *list.end.prev.prev.obj)
	assert.Equal(t, obj1, *list.end.prev.prev.next.obj)
	assert.Equal(t, obj2, *list.end.prev.prev.next.next.obj)

	// Replace first
	_ = list.Replace(0, obj5)
	assert.Equal(t, obj5, list.AtPtr(0))
	assert.Equal(t, obj2, list.AtPtr(1))

	// Replace in a middle
	err = list.Replace(2, obj5)
	assert.Nil(t, err)
	assert.Equal(t, obj2, list.AtPtr(1))
	assert.Equal(t, obj5, list.AtPtr(2))
	assert.Equal(t, obj4, list.AtPtr(3))

	// Replace last
	err = list.Replace(6, obj1)
	assert.Nil(t, err)
	assert.Equal(t, obj1, list.AtPtr(5))
	assert.Equal(t, obj1, list.AtPtr(6))

	// Replace last
	err = list.ReplaceLast(obj3)
	assert.Nil(t, err)
	assert.Equal(t, 7, list.Size())
	assert.Equal(t, obj1, list.AtPtr(5))
	assert.Equal(t, obj3, list.AtPtr(6))

	// Renew
	list.Clear()
	assert.Equal(t, true, list.IsEmpty())
	assert.Equal(t, 0, list.Size())

	list.Append(obj1, obj2, obj3, obj4, obj5)
	assert.Equal(t, 5, list.Size())

	// DeleteAt (invalid index)
	objRes, err = list.DeleteAt(5)
	assert.ErrorIs(t, err, ErrInvalidIndex)
	assert.Nil(t, objRes)

	// Delete At (last position)
	objRes, err = list.DeleteAt(4)
	assert.Nil(t, err)
	assert.Equal(t, obj5, objRes)
	assert.Equal(t, 4, list.Size())
	assert.Equal(t, obj4, list.AtPtr(3))

	//Check internal
	assert.Equal(t, obj4, *list.end.obj)
	assert.Nil(t, list.end.next)
	assert.Equal(t, obj3, *list.end.prev.obj)
	assert.Equal(t, obj4, *list.end.prev.next.obj)

	list.Append(obj5)

	// DeleteAt (middle)
	objRes, err = list.DeleteAt(2)
	assert.Nil(t, err)
	assert.Equal(t, 4, list.Size())
	assert.Equal(t, obj3, objRes)
	assert.Equal(t, obj2, list.AtPtr(1))
	assert.Equal(t, obj4, list.AtPtr(2))

	// Check internal connections
	assert.Equal(t, obj2, *list.home.next.obj)
	assert.Equal(t, obj4, *list.home.next.next.obj)
	assert.Equal(t, obj2, *list.home.next.next.prev.obj)

	err = list.Insert(2, obj3)
	assert.Nil(t, err)
	assert.Equal(t, 5, list.Size())
	assert.Equal(t, obj4, list.AtPtr(3))
	assert.Equal(t, obj3, list.AtPtr(2))
	assert.Equal(t, obj2, list.AtPtr(1))

	// DeleteAt (first position)
	objRes, err = list.DeleteAt(0)
	assert.Nil(t, err)
	assert.Equal(t, obj1, objRes)
	assert.Equal(t, 4, list.Size())
	assert.Equal(t, obj2, list.AtPtr(0))
	// Check internal connections
	assert.Equal(t, obj2, *list.home.obj)
	assert.Equal(t, obj3, *list.home.next.obj)
	assert.Equal(t, obj2, *list.home.next.prev.obj)
	assert.Nil(t, list.home.next.prev.prev)

	// DeleteLast
	objRes, err = list.DeleteLast()
	assert.Nil(t, err)
	assert.Equal(t, obj5, objRes)
	assert.Equal(t, 3, list.Size())
	assert.Equal(t, obj4, list.AtPtr(2))
	assert.Equal(t, obj4, list.LastObjectPtr())

	objRes, err = list.DeleteLast()
	assert.Nil(t, err)
	assert.Equal(t, obj4, objRes)
	assert.Equal(t, 2, list.Size())
	assert.Equal(t, obj3, list.AtPtr(1))
	assert.Equal(t, obj3, list.LastObjectPtr())

	objRes, err = list.DeleteLast()
	assert.Nil(t, err)
	assert.Equal(t, obj3, objRes)
	assert.Equal(t, 1, list.Size())
	assert.Equal(t, obj2, list.AtPtr(0))
	assert.Equal(t, obj2, list.LastObjectPtr())

	// Check internal
	assert.Equal(t, obj2, *list.home.obj)
	assert.Nil(t, list.home.next)
	assert.Nil(t, list.home.prev)

	assert.Equal(t, list.home, list.end)

	objRes, err = list.DeleteLast()
	assert.Nil(t, err)
	assert.Equal(t, obj2, objRes)
	assert.Equal(t, 0, list.Size())
	assert.Nil(t, list.AtPtr(0))

	// Check internal
	assert.Nil(t, list.home)
	assert.Nil(t, list.end)

	// Renew array
	err = list.Insert(0, obj1, obj2, obj3)
	assert.Nil(t, err)

	// Delete from the beginning until array is not empty
	objRes, err = list.DeleteAt(0)
	assert.Nil(t, err)
	assert.Equal(t, obj1, objRes)
	assert.Equal(t, 2, list.Size())
	assert.Equal(t, obj2, list.AtPtr(0))
	// check internal
	assert.Equal(t, obj2, *list.home.obj)
	assert.Nil(t, list.home.prev)
	assert.Equal(t, obj3, *list.home.next.obj)

	objRes, err = list.DeleteAt(0)
	assert.Nil(t, err)
	assert.Equal(t, obj2, objRes)
	assert.Equal(t, 1, list.Size())
	assert.Equal(t, obj3, list.AtPtr(0))

	objRes, err = list.DeleteAt(0)
	assert.Nil(t, err)
	assert.Equal(t, obj3, objRes)
	assert.Equal(t, 0, list.Size())

	// Delete from empty
	objRes, err = list.DeleteAt(0)
	assert.Nil(t, err)
	assert.Nil(t, objRes)

	// Copy
	list.Append(obj1, obj2, obj3, obj4)
	list2 = list.Copy()
	assert.Equal(t, 4, list2.Size())
	assert.Equal(t, obj1, list2.AtPtr(0))
	assert.Equal(t, obj2, list2.AtPtr(1))
	assert.Equal(t, obj3, list2.AtPtr(2))
	assert.Equal(t, obj4, list2.AtPtr(3))

	list.Clear()
	list2 = list.Copy()
	assert.Equal(t, 0, list2.Size())
	assert.Nil(t, list2.AtPtr(0))

	list.Append(obj1)
	list2 = list.Copy()
	assert.Equal(t, 1, list2.Size())
	assert.Equal(t, obj1, list2.AtPtr(0))

	// AppendList
	list.Clear()
	list2.Clear()
	list.Append(obj1, obj2, obj3)
	list2.Append(obj4, obj5)
	list3 := list.AppendList(list2)
	assert.Equal(t, 5, list3.Size())

	assert.Equal(t, obj1, list3.AtPtr(0))
	assert.Equal(t, obj2, list3.AtPtr(1))
	assert.Equal(t, obj3, list3.AtPtr(2))
	assert.Equal(t, obj4, list3.AtPtr(3))
	assert.Equal(t, obj5, list3.AtPtr(4))

	// Check internal stuff
	assert.Equal(t, obj1, *list3.home.obj)
	assert.Equal(t, obj3, *list3.home.next.next.obj)
	assert.Equal(t, obj4, *list3.home.next.next.next.obj)
	assert.Equal(t, obj2, *list3.home.next.next.next.prev.prev.obj)

	list.Clear()
	list3 = list.AppendList(list2)

	assert.Equal(t, 2, list3.Size())
	assert.Equal(t, obj4, list3.AtPtr(0))
	assert.Equal(t, obj5, list3.AtPtr(1))

	list.Clear()
	list.Append(obj1, obj2, obj3)
	err = list.SpliceAtPos(3, list2)

	assert.Nil(t, err)
	assert.Equal(t, 5, list.Size())
	assert.Equal(t, 0, list2.Size())
	assert.Equal(t, obj3, list.AtPtr(2))
	assert.Equal(t, obj4, list.AtPtr(3))
	assert.Equal(t, obj5, list.AtPtr(4))

	list.Clear()
	list2.Clear()
	list2.Append(obj1)
	list.Splice(list2)
	assert.Equal(t, 1, list.Size())
	assert.Equal(t, 0, list2.Size())

	list2.Append(obj2, obj3)
	list.Splice(list2)
	assert.Equal(t, 3, list.Size())

	// Contains
	list.Clear()
	list.Append(obj1, obj2, obj3, obj4)
	ok = list.Contains(obj1)
	assert.True(t, ok)

	ok = list.Contains(obj2)
	assert.True(t, ok)

	ok = list.Contains(obj3)
	assert.True(t, ok)

	ok = list.Contains(obj4)
	assert.True(t, ok)

	ok = list.Contains(obj5)
	assert.False(t, ok)

	// AppendUnique
	list.Clear()
	list.AppendUnique(obj1)
	assert.Equal(t, 1, list.Size())

	list.AppendUnique(obj1)
	assert.Equal(t, 1, list.Size())
	assert.Equal(t, obj1, list.AtPtr(0))
	assert.Nil(t, list.AtPtr(1))

	list.AppendUnique(obj2)
	assert.Equal(t, 2, list.Size())
	list.AppendUnique(obj2)
	assert.Equal(t, 2, list.Size())
	assert.Equal(t, obj2, list.AtPtr(1))
	assert.Nil(t, list.AtPtr(2))

	list.AppendUnique(obj2, obj3, obj4, obj5)
	assert.Equal(t, 5, list.Size())
	assert.Equal(t, obj5, list.AtPtr(4))

	list.AppendUnique(obj5cp)
	assert.Equal(t, 5, list.Size())
	assert.Equal(t, obj5, list.AtPtr(4))

	list.AppendUnique(nil)
	assert.Equal(t, 6, list.Size())
	assert.Equal(t, obj5, list.AtPtr(4))

	// Cycle: 1
	intlist := New[int]()
	cSlice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	intlist.Append(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	it := intlist.Iterator()

	for value, ok := it.SetIndex(0); ok; value, ok = it.NextValue() {
		index := it.Index()
		cValue := cSlice[index]

		assert.Equal(t, cValue, value)
	}

	// Cycle: 2
	it.Reset()
	for it.Next() {
		value, _ := it.Value()
		index := it.Index()
		cValue := cSlice[index]

		assert.Equal(t, cValue, value)
	}

	// Reverse Cycle 1
	it.Reset()
	for value, ok := it.SetLast(); ok; value, ok = it.PrevValue() {
		index := it.Index()
		cValue := cSlice[index]

		assert.Equal(t, cValue, value)
	}

	// Reverse Cycle 2
	it.Reset()
	for it.Prev() {
		value, _ := it.Value()
		index := it.Index()
		cValue := cSlice[index]

		assert.Equal(t, cValue, value)
	}

	//-------------------------------------------
	// Iterator with range
	it = intlist.Iterator(2, 6)

	for value, ok := it.SetFirst(); ok; value, ok = it.NextValue() {
		index := it.Index()
		cValue := cSlice[index]

		assert.Equal(t, cValue, value)
	}

	// Cycle: 2
	it.Reset(2, 6)
	for it.Next() {
		value, _ := it.Value()
		index := it.Index()
		cValue := cSlice[index]

		assert.Equal(t, cValue, value)
	}

	// Reverse Cycle 1
	it.Reset(2, 6)
	for value, ok := it.SetLast(); ok; value, ok = it.PrevValue() {
		index := it.Index()
		cValue := cSlice[index]

		assert.Equal(t, cValue, value)
	}

	// Reverse Cycle 2
	it.Reset(2, 6)
	for it.Prev() {
		value, _ := it.Value()
		index := it.Index()
		cValue := cSlice[index]

		assert.Equal(t, cValue, value)
	}

	// ========== ContainsSome Tests ==========
	// Returns false for empty list
	list.Clear()
	ok = list.ContainsSome(obj1, obj2)
	assert.False(t, ok)

	// Returns false when no objects match
	list.Append(obj1, obj2, obj3)
	ok = list.ContainsSome(obj4, obj5)
	assert.False(t, ok)

	// Returns true when at least one object matches
	ok = list.ContainsSome(obj1, obj4, obj5)
	assert.True(t, ok)

	// Returns true when multiple objects match
	ok = list.ContainsSome(obj1, obj2, obj3)
	assert.True(t, ok)

	// Returns false when called with no arguments
	ok = list.ContainsSome()
	assert.False(t, ok)

	// Works correctly with nil values
	list.Clear()
	list.Append(obj1, obj2)
	list.Append(*new(*teststruct)) // nil value
	ok = list.ContainsSome(*new(*teststruct))
	assert.True(t, ok)

	// ========== CopyRange Tests ==========
	list.Clear()
	list.Append(obj1, obj2, obj3, obj4, obj5)

	// Copies a range correctly
	copiedRange, err := list.CopyRange(1, 3)
	assert.Nil(t, err)
	assert.Equal(t, 3, copiedRange.Size())
	assert.Equal(t, obj2, copiedRange.AtPtr(0))
	assert.Equal(t, obj3, copiedRange.AtPtr(1))
	assert.Equal(t, obj4, copiedRange.AtPtr(2))

	// Returns error for invalid indices
	_, err = list.CopyRange(-1, 2)
	assert.ErrorIs(t, err, ErrInvalidIndex)

	_, err = list.CopyRange(0, 10)
	assert.ErrorIs(t, err, ErrInvalidIndex)

	// Returns error when fromPos > toPos
	_, err = list.CopyRange(3, 1)
	assert.ErrorIs(t, err, ErrInvalidIndex)

	// Handles single element range
	copiedRange, err = list.CopyRange(2, 2)
	assert.Nil(t, err)
	assert.Equal(t, 1, copiedRange.Size())
	assert.Equal(t, obj3, copiedRange.AtPtr(0))

	// Handles full range correctly
	copiedRange, err = list.CopyRange(0, 4)
	assert.Nil(t, err)
	assert.Equal(t, 5, copiedRange.Size())
	assert.Equal(t, obj1, copiedRange.AtPtr(0))
	assert.Equal(t, obj5, copiedRange.AtPtr(4))

	// Works with empty list
	list.Clear()
	_, err = list.CopyRange(0, 0)
	assert.ErrorIs(t, err, ErrInvalidIndex)

	// ========== DeepCopy and DeepCopyRange Tests ==========
	list.Clear()
	list.Append(obj1, obj2, obj3, obj4, obj5)

	// DeepCopy creates truly independent copies
	deepCopyFn := func(obj *teststruct) *teststruct {
		return &teststruct{Num: obj.Num, Str: obj.Str}
	}

	deepList := list.DeepCopy(deepCopyFn)
	assert.Equal(t, 5, deepList.Size())

	// Modify original and verify copy is unaffected
	obj1Original := list.AtPtr(0)
	obj1Original.Num = 999
	obj1Original.Str = "modified"

	obj1Deep := deepList.AtPtr(0)
	assert.Equal(t, 1, obj1Deep.Num)
	assert.Equal(t, "obj1", obj1Deep.Str)

	// Reset obj1 for subsequent tests
	obj1.Num = 1
	obj1.Str = "obj1"

	// DeepCopyRange works for various ranges
	deepRange, err := list.DeepCopyRange(1, 3, deepCopyFn)
	assert.Nil(t, err)
	assert.Equal(t, 3, deepRange.Size())
	assert.Equal(t, obj2.Num, deepRange.AtPtr(0).Num)
	assert.Equal(t, obj4.Num, deepRange.AtPtr(2).Num)

	// Returns error when deepCopyFn is nil
	emptyList := list.DeepCopy(nil)
	assert.Equal(t, 0, emptyList.Size())

	_, err = list.DeepCopyRange(0, 2, nil)
	assert.ErrorIs(t, err, ErrNoClosure)

	// Returns empty list for empty input
	list.Clear()
	deepList = list.DeepCopy(deepCopyFn)
	assert.Equal(t, 0, deepList.Size())

	// Handles struct pointers correctly
	list.Append(obj1, obj2, obj3)
	deepList = list.DeepCopy(deepCopyFn)
	assert.Equal(t, 3, deepList.Size())

	// Verify deep copy independence
	originalPtr := list.AtPtr(0)
	deepPtr := deepList.AtPtr(0)
	assert.NotSame(t, originalPtr, deepPtr)
	assert.Equal(t, originalPtr.Num, deepPtr.Num)
	assert.Equal(t, originalPtr.Str, deepPtr.Str)

	// ========== Swap Tests ==========
	list.Clear()
	list.Append(obj1, obj2, obj3, obj4, obj5)

	// Swaps two elements correctly
	err = list.Swap(1, 3)
	assert.Nil(t, err)
	assert.Equal(t, obj1, list.AtPtr(0))
	assert.Equal(t, obj4, list.AtPtr(1))
	assert.Equal(t, obj3, list.AtPtr(2))
	assert.Equal(t, obj2, list.AtPtr(3))
	assert.Equal(t, obj5, list.AtPtr(4))

	// Returns error for invalid indices
	err = list.Swap(-1, 2)
	assert.ErrorIs(t, err, ErrInvalidIndex)

	err = list.Swap(1, 10)
	assert.ErrorIs(t, err, ErrInvalidIndex)

	// Returns nil when swapping same index
	err = list.Swap(2, 2)
	assert.Nil(t, err)
	assert.Equal(t, obj3, list.AtPtr(2))

	// Works at boundaries (first/last elements)
	list.Clear()
	list.Append(obj1, obj2, obj3, obj4, obj5)
	err = list.Swap(0, 4)
	assert.Nil(t, err)
	assert.Equal(t, obj5, list.AtPtr(0))
	assert.Equal(t, obj1, list.AtPtr(4))

	// Preserves list size
	initialSize := list.Size()
	err = list.Swap(1, 3)
	assert.Nil(t, err)
	assert.Equal(t, initialSize, list.Size())

	// ========== Find Tests ==========
	list.Clear()
	list.Append(obj1, obj2, obj3, obj4, obj5)

	// Returns new list with matching elements
	foundList := list.Find(func(index int, object *teststruct) bool {
		return object.Num > 2
	})
	assert.Equal(t, 3, foundList.Size())
	assert.Equal(t, obj3, foundList.AtPtr(0))
	assert.Equal(t, obj4, foundList.AtPtr(1))
	assert.Equal(t, obj5, foundList.AtPtr(2))

	// Returns empty list when nothing matches
	foundList = list.Find(func(index int, object *teststruct) bool {
		return object.Num > 10
	})
	assert.Equal(t, 0, foundList.Size())

	// Preserves original list
	assert.Equal(t, 5, list.Size())
	assert.Equal(t, obj1, list.AtPtr(0))
	assert.Equal(t, obj5, list.AtPtr(4))

	// Works with various predicates
	foundList = list.Find(func(index int, object *teststruct) bool {
		return index%2 == 0
	})
	assert.Equal(t, 3, foundList.Size())
	assert.Equal(t, obj1, foundList.AtPtr(0))
	assert.Equal(t, obj3, foundList.AtPtr(1))
	assert.Equal(t, obj5, foundList.AtPtr(2))

	// ========== Range Iterator Tests ==========
	list.Clear()
	list.Append(obj1, obj2, obj3, obj4, obj5)

	// All(): direct pass from start to end, no options
	forwardCount := 0
	var forwardIndices []int
	for i, ptr := range list.All() {
		forwardCount++
		forwardIndices = append(forwardIndices, i)
		if i == 0 {
			assert.Equal(t, obj1, ptr)
		}
		if i == 4 {
			assert.Equal(t, obj5, ptr)
		}
	}
	assert.Equal(t, 5, forwardCount)
	assert.Equal(t, []int{0, 1, 2, 3, 4}, forwardIndices)

	// All(): start from first element (pos=1) to end
	n := 0
	var forwardFromFirst []int
	for i, ptr := range list.All(WithPos(1)) {
		forwardFromFirst = append(forwardFromFirst, i)
		if i == 1 {
			assert.Equal(t, obj2, ptr)
		}
		if i == 4 {
			assert.Equal(t, obj5, ptr)
		}
		n++
	}
	assert.Equal(t, 4, n)
	assert.Equal(t, []int{1, 2, 3, 4}, forwardFromFirst)

	// All(): from first element, count=3
	n = 0
	for i, ptr := range list.All(WithPos(1), WithCount(3)) {
		switch i {
		case 1:
			assert.Equal(t, obj2, ptr)
		case 2:
			assert.Equal(t, obj3, ptr)
		case 3:
			assert.Equal(t, obj4, ptr)
		}
		assert.Contains(t, []int{1, 2, 3}, i)
		n++
	}
	assert.Equal(t, 3, n)

	// All(): from first element, count=4 (boundary within range)
	n = 0
	for i, ptr := range list.All(WithPos(1), WithCount(4)) {
		switch i {
		case 1:
			assert.Equal(t, obj2, ptr)
		case 2:
			assert.Equal(t, obj3, ptr)
		case 3:
			assert.Equal(t, obj4, ptr)
		case 4:
			assert.Equal(t, obj5, ptr)
		}
		assert.Contains(t, []int{1, 2, 3, 4}, i)
		n++
	}
	assert.Equal(t, 4, n)

	// All(): from first element, count=5 (over range should clamp)
	n = 0
	for i, ptr := range list.All(WithPos(1), WithCount(5)) {
		switch i {
		case 1:
			assert.Equal(t, obj2, ptr)
		case 2:
			assert.Equal(t, obj3, ptr)
		case 3:
			assert.Equal(t, obj4, ptr)
		case 4:
			assert.Equal(t, obj5, ptr)
		}
		assert.Contains(t, []int{1, 2, 3, 4}, i)
		n++
	}
	assert.Equal(t, 4, n)

	// All(): negative count should be treated as absolute
	n = 0
	for i := range list.All(WithCount(-2)) {
		assert.Contains(t, []int{0, 1}, i)
		n++
	}
	assert.Equal(t, 2, n)

	// All(): boundary positions within range
	n = 0
	for i, ptr := range list.All(WithPos(0), WithCount(1)) {
		assert.Equal(t, 0, i)
		assert.Equal(t, obj1, ptr)
		n++
	}
	assert.Equal(t, 1, n)

	n = 0
	for i, ptr := range list.All(WithPos(4), WithCount(1)) {
		assert.Equal(t, 4, i)
		assert.Equal(t, obj5, ptr)
		n++
	}
	assert.Equal(t, 1, n)

	// All(): invalid positive position should panic
	assert.Panics(t, func() {
		for range list.All(WithPos(5)) {
		}
	})

	// All(): invalid negative position should panic
	assert.Panics(t, func() {
		for range list.All(WithPos(-3)) {
		}
	})

	// Backward(): direct reverse pass from end to start
	backwardCount := 0
	var backwardIndices []int
	for i, ptr := range list.Backward() {
		backwardCount++
		backwardIndices = append(backwardIndices, i)
		if i == 4 {
			assert.Equal(t, obj5, ptr)
		}
		if i == 0 {
			assert.Equal(t, obj1, ptr)
		}
	}
	assert.Equal(t, 5, backwardCount)
	assert.Equal(t, []int{4, 3, 2, 1, 0}, backwardIndices)

	// Backward(): start from first element (pos=1), go to start
	n = 0
	var backwardFromFirst []int
	for i, ptr := range list.Backward(WithPos(1)) {
		backwardFromFirst = append(backwardFromFirst, i)
		if i == 1 {
			assert.Equal(t, obj2, ptr)
		}
		if i == 0 {
			assert.Equal(t, obj1, ptr)
		}
		n++
	}
	assert.Equal(t, 2, n)
	assert.Equal(t, []int{1, 0}, backwardFromFirst)

	// Backward(): from pos=1, count=3 (over range should clamp)
	n = 0
	for i, ptr := range list.Backward(WithPos(1), WithCount(3)) {
		switch i {
		case 1:
			assert.Equal(t, obj2, ptr)
		case 0:
			assert.Equal(t, obj1, ptr)
		}
		assert.Contains(t, []int{1, 0}, i)
		n++
	}
	assert.Equal(t, 2, n)

	// Backward(): from pos=4, count=5 (full range)
	n = 0
	for i, ptr := range list.Backward(WithPos(4), WithCount(5)) {
		switch i {
		case 4:
			assert.Equal(t, obj5, ptr)
		case 3:
			assert.Equal(t, obj4, ptr)
		case 2:
			assert.Equal(t, obj3, ptr)
		case 1:
			assert.Equal(t, obj2, ptr)
		case 0:
			assert.Equal(t, obj1, ptr)
		}
		assert.Contains(t, []int{4, 3, 2, 1, 0}, i)
		n++
	}
	assert.Equal(t, 5, n)

	// Backward(): from pos=2, count=3
	n = 0
	for i, ptr := range list.Backward(WithPos(2), WithCount(3)) {
		switch i {
		case 2:
			assert.Equal(t, obj3, ptr)
		case 1:
			assert.Equal(t, obj2, ptr)
		case 0:
			assert.Equal(t, obj1, ptr)
		}
		assert.Contains(t, []int{2, 1, 0}, i)
		n++
	}
	assert.Equal(t, 3, n)

	// Backward(): from pos=2, count=4 (over range should clamp)
	n = 0
	for i, ptr := range list.Backward(WithPos(2), WithCount(4)) {
		switch i {
		case 2:
			assert.Equal(t, obj3, ptr)
		case 1:
			assert.Equal(t, obj2, ptr)
		case 0:
			assert.Equal(t, obj1, ptr)
		}
		assert.Contains(t, []int{2, 1, 0}, i)
		n++
	}
	assert.Equal(t, 3, n)

	// Backward(): negative count should be treated as absolute
	n = 0
	for i := range list.Backward(WithCount(-2)) {
		assert.Contains(t, []int{4, 3}, i)
		n++
	}
	assert.Equal(t, 2, n)

	// Backward(): boundary positions within range
	n = 0
	for i, ptr := range list.Backward(WithPos(0), WithCount(1)) {
		assert.Equal(t, 0, i)
		assert.Equal(t, obj1, ptr)
		n++
	}
	assert.Equal(t, 1, n)

	n = 0
	for i, ptr := range list.Backward(WithPos(4), WithCount(1)) {
		assert.Equal(t, 4, i)
		assert.Equal(t, obj5, ptr)
		n++
	}
	assert.Equal(t, 1, n)

	// Backward(): invalid positive position should panic
	assert.Panics(t, func() {
		for range list.Backward(WithPos(5)) {
		}
	})

	// Backward(): invalid negative position should panic
	assert.Panics(t, func() {
		for range list.Backward(WithPos(-3)) {
		}
	})

	// Values(): iterates without index
	list.Clear()
	list.Append(obj1, obj2, obj3)
	valuesCount := 0
	valuesResult := []*teststruct{}
	for v := range list.Values() {
		valuesCount++
		valuesResult = append(valuesResult, v)
	}
	assert.Equal(t, 3, valuesCount)
	assert.Equal(t, obj1, valuesResult[0])
	assert.Equal(t, obj2, valuesResult[1])
	assert.Equal(t, obj3, valuesResult[2])

	// ValuesBackward(): iterates in reverse without index
	valuesRev := []*teststruct{}
	for v := range list.ValuesBackward() {
		valuesRev = append(valuesRev, v)
	}
	assert.Equal(t, []*teststruct{obj3, obj2, obj1}, valuesRev)

	// Values(): RangeOptions
	valuesFromPos := []*teststruct{}
	for v := range list.Values(WithPos(1)) {
		valuesFromPos = append(valuesFromPos, v)
	}
	assert.Equal(t, []*teststruct{obj2, obj3}, valuesFromPos)

	valuesWithCount := []*teststruct{}
	for v := range list.Values(WithPos(1), WithCount(2)) {
		valuesWithCount = append(valuesWithCount, v)
	}
	assert.Equal(t, []*teststruct{obj2, obj3}, valuesWithCount)

	// ValuesBackward(): RangeOptions
	valuesBackFromPos := []*teststruct{}
	for v := range list.ValuesBackward(WithPos(1)) {
		valuesBackFromPos = append(valuesBackFromPos, v)
	}
	assert.Equal(t, []*teststruct{obj2, obj1}, valuesBackFromPos)

	valuesBackWithCount := []*teststruct{}
	for v := range list.ValuesBackward(WithPos(2), WithCount(2)) {
		valuesBackWithCount = append(valuesBackWithCount, v)
	}
	assert.Equal(t, []*teststruct{obj3, obj2}, valuesBackWithCount)

	list.Append(obj4, obj5)
	valuesBackClamp := []*teststruct{}
	for v := range list.ValuesBackward(WithPos(4), WithCount(10)) {
		valuesBackClamp = append(valuesBackClamp, v)
	}
	assert.Equal(t, []*teststruct{obj5, obj4, obj3, obj2, obj1}, valuesBackClamp)

	// Values(): invalid position should panic
	assert.Panics(t, func() {
		for range list.Values(WithPos(-1)) {
		}
	})

	// ValuesBackward(): invalid position should panic
	assert.Panics(t, func() {
		for range list.ValuesBackward(WithPos(5)) {
		}
	})

	// ToValues(): returns all values as-is (including nil)
	list.Clear()
	//list.Append(obj1, (*teststruct)(nil), obj2)
	list.Append(obj1, nil, obj2)

	vals := []*teststruct{}
	for v := range ToValues(list.All()) {
		if v != nil { // manually skip nil if needed
			vals = append(vals, v)
		}
	}
	assert.Equal(t, []*teststruct{obj1, obj2}, vals)

	valsFromPos := []*teststruct{}
	for v := range ToValues(list.All(WithPos(1), WithCount(2))) {
		if v != nil {
			valsFromPos = append(valsFromPos, v)
		}
	}
	assert.Equal(t, []*teststruct{obj2}, valsFromPos)

	valsBack := []*teststruct{}
	for v := range ToValues(list.Backward()) {
		if v != nil {
			valsBack = append(valsBack, v)
		}
	}
	assert.Equal(t, []*teststruct{obj2, obj1}, valsBack)

	// Filter(): keep even indices
	list.Clear()
	list.Append(obj1, obj2, obj3, obj4, obj5)
	filtered := []*teststruct{}
	for i, obj := range Filter(list.All(), func(index int, obj *teststruct) bool {
		return obj != nil && index%2 == 0
	}) {
		assert.Contains(t, []int{0, 2, 4}, i)
		filtered = append(filtered, obj)
	}
	assert.Equal(t, []*teststruct{obj1, obj3, obj5}, filtered)

	// TakeWhile(): stop on first failure
	taken := []*teststruct{}
	for _, obj := range TakeWhile(list.All(), func(_ int, obj *teststruct) bool {
		return obj != nil && obj.Num < 4
	}) {
		taken = append(taken, obj)
	}
	assert.Equal(t, []*teststruct{obj1, obj2, obj3}, taken)

	// SkipWhile(): skip prefix, then yield rest
	skipped := []*teststruct{}
	for _, obj := range SkipWhile(list.All(), func(_ int, obj *teststruct) bool {
		return obj != nil && obj.Num < 3
	}) {
		skipped = append(skipped, obj)
	}
	assert.Equal(t, []*teststruct{obj3, obj4, obj5}, skipped)

	// Map(): transform values
	strs := []string{}
	for s := range Map(list.Values(), func(v *teststruct) string {
		return fmt.Sprintf("%d-%s", v.Num, v.Str)
	}) {
		strs = append(strs, s)
	}
	assert.Equal(t, []string{"1-obj1", "2-obj2", "3-obj3", "4-obj4", "5-obj5"}, strs)

	// AnyMatch(): true and false cases
	hasEven := AnyMatch(list.All(), func(_ int, obj *teststruct) bool {
		return obj != nil && obj.Num%2 == 0
	})
	assert.Equal(t, true, hasEven)
	noneBig := AnyMatch(list.All(), func(_ int, obj *teststruct) bool {
		return obj != nil && obj.Num > 10
	})
	assert.Equal(t, false, noneBig)

	// AllMatch(): true and false cases
	allPositive := AllMatch(list.All(), func(_ int, obj *teststruct) bool {
		return obj != nil && obj.Num > 0
	})
	assert.Equal(t, true, allPositive)
	notAllBig := AllMatch(list.All(), func(_ int, obj *teststruct) bool {
		return obj != nil && obj.Num > 3
	})
	assert.Equal(t, false, notAllBig)

	// Terminal functions: invalid predicate should panic
	assert.Panics(t, func() {
		AnyMatch(list.All(), nil)
	})
	assert.Panics(t, func() {
		AllMatch(list.All(), nil)
	})

	// Empty list behavior
	list.Clear()
	emptyForwardCount := 0
	for range list.All() {
		emptyForwardCount++
	}
	assert.Equal(t, 0, emptyForwardCount)

	emptyBackwardCount := 0
	for range list.Backward() {
		emptyBackwardCount++
	}
	assert.Equal(t, 0, emptyBackwardCount)

	emptyValuesCount := 0
	for range list.Values() {
		emptyValuesCount++
	}
	assert.Equal(t, 0, emptyValuesCount)

	emptyValuesBackwardCount := 0
	for range list.ValuesBackward() {
		emptyValuesBackwardCount++
	}
	assert.Equal(t, 0, emptyValuesBackwardCount)

	// Terminal functions on empty list
	assert.Equal(t, false, AnyMatch(list.All(), func(_ int, obj *teststruct) bool {
		return obj != nil
	}))
	assert.Equal(t, true, AllMatch(list.All(), func(_ int, obj *teststruct) bool {
		return obj != nil
	}))

	// ========== Functional Chaining Tests ==========
	list.Clear()

	// Clear().Append().Modify() then Sort()
	chainedList := list.Clear().Append(obj3, obj1, obj2).Modify(func(index int, obj *teststruct) *teststruct {
		obj.Num = obj.Num * 10
		return obj
	})
	chainedList.Sort(func(a, b *teststruct) bool {
		return a.Num < b.Num
	})

	assert.Equal(t, 3, chainedList.Size())
	assert.Equal(t, 10, chainedList.AtPtr(0).Num)
	assert.Equal(t, 20, chainedList.AtPtr(1).Num)
	assert.Equal(t, 30, chainedList.AtPtr(2).Num)

	// Reset values for subsequent tests
	obj1.Num = 1
	obj2.Num = 2
	obj3.Num = 3

	// AppendList returns self for chaining
	list.Clear()
	list2.Clear()
	list.Append(obj1, obj2)
	list2.Append(obj3, obj4)

	result := list.AppendList(list2).Append(obj5)
	assert.Equal(t, 5, result.Size())
	assert.Equal(t, obj1, result.AtPtr(0))
	assert.Equal(t, obj5, result.AtPtr(4))

	// Splice returns self for chaining
	list.Clear()
	list2.Clear()
	list.Append(obj1, obj2)
	list2.Append(obj3)

	result = list.Splice(list2).Append(obj4)
	assert.Equal(t, 4, result.Size())
	assert.Equal(t, 0, list2.Size())
	assert.Equal(t, obj3, result.AtPtr(2))
	assert.Equal(t, obj4, result.AtPtr(3))

	concurrencyRWTest(t)
	concurrencyAppendTest(t)

	sortTest(t)
}

func concurrencyRWTest(t *testing.T) {
	const listSize = 10
	const attemptNum = 1000000

	list := XList[int]{}
	wg := sync.WaitGroup{}

	// A few goroutines write and read simultaneously
	for range listSize {
		list.Append(0)
	}

	dirChanges := func() {
		for i := 0; i < attemptNum; i++ {
			list.Modify(func(index int, obj int) int {
				obj += 1
				return obj
			})
		}

		wg.Done()
	}

	revChanges := func() {
		for i := 0; i < attemptNum; i++ {
			list.ModifyRev(func(index int, obj int) int {
				obj += 1
				return obj
			})
		}

		wg.Done()
	}

	wg.Add(2)

	dirChanges()
	revChanges()

	wg.Wait()

	// Checking
	for i := 0; i < listSize; i++ {
		v, ok := list.At(i)
		assert.Equal(t, true, ok)
		assert.Equal(t, attemptNum*2, v)
	}

}

func concurrencyAppendTest(t *testing.T) {
	const attemptNum = 1000000
	const threads = 10

	wg := sync.WaitGroup{}
	list := XList[int]{}

	appendChanges := func() {
		for i := 0; i < attemptNum; i++ {
			list.Append(i)
		}

		wg.Done()
	}

	wg.Add(threads)

	for range threads {
		go appendChanges()
	}

	wg.Wait()

	assert.Equal(t, attemptNum*threads, list.Size())

	time.Sleep(1 * time.Second)
}

func sortTest(t *testing.T) {
	const listSize = 10000
	xlist := XList[int64]{}

	// Generates random numbers
	gen := rand.New(rand.NewSource(time.Now().UnixNano()))

	t.Log(`-- Sort Test starting --`)
	t.Log(`Start generate random numbers...`)

	for range listSize {
		rnum := gen.Intn(listSize) + 1
		xlist.Append(int64(rnum))
	}

	t.Log("Done.")

	t.Log("Start sorting...")
	time.Sleep(100 * time.Millisecond)

	// Run sort
	start := time.Now()
	xlist.Sort(func(a, b int64) bool { return a < b })

	t.Log("Done.")
	t.Logf("Sort of %d elements done for %.2f seconds\n", listSize, time.Since(start).Seconds())

	// Check sort
	t.Log("Check sorting...")

	assert.Equal(t, listSize, xlist.Size())

	it := xlist.Iterator()

	var value, prevValue int64
	var ok bool

	// for i := 0; i < listSize; i++ {
	for value, ok = it.SetIndex(0); ok; value, ok = it.NextValue() {

		if it.Index() == 0 {
			prevValue = value
			continue
		}

		ss := fmt.Sprintf("Previous value must be less or equal than the current value at index %d", it.Index())
		assert.LessOrEqual(t, prevValue, value, ss)

		if prevValue > value {
			index := it.Index()

			v, _ := xlist.At(index)
			pv, _ := xlist.At(index - 1)

			ss = fmt.Sprintf("Prev value %d > %d at index %d", pv, v, index)
			t.Log(ss)
		}

		prevValue = value
	}

	t.Log("------------------------------------------------------")

	it.Reset()

	for it.Next() {
		value, _ = it.Value()
		index := it.Index()

		if index == 0 {
			prevValue = value
			continue
		}

		ss := fmt.Sprintf("Previous value must be less or equal than the current value at index %d", it.Index())
		assert.LessOrEqual(t, prevValue, value, ss)

		prevValue = value
	}

	t.Log("------------------------------------------------------")

	it.Reset()
	for it.Prev() {
		value, _ = it.Value()
		index := it.Index()

		if index == xlist.Size()-1 {
			prevValue = value
			continue
		}

		ss := fmt.Sprintf("Previous value must be less or equal than the current value at index %d", it.Index())
		// assert.LessOrEqual(t, prevValue, value, ss)
		assert.GreaterOrEqual(t, prevValue, value, ss)

		prevValue = value
	}

	t.Log("Done.")
}
