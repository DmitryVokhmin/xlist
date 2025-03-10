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

func Test(t *testing.T) {
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

	// Add
	list.Clear()
	list2.Clear()
	list.Append(obj1, obj2, obj3)
	list2.Append(obj4, obj5)
	list3 := list.Add(list2)
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
	list3 = list.Add(list2)

	assert.Equal(t, 2, list3.Size())
	assert.Equal(t, obj4, list3.AtPtr(0))
	assert.Equal(t, obj5, list3.AtPtr(1))

	list.Clear()
	list.Append(obj1, obj2, obj3)
	err = list.MoveAtPos(3, list2)

	assert.Nil(t, err)
	assert.Equal(t, 5, list.Size())
	assert.Equal(t, 0, list2.Size())
	assert.Equal(t, obj3, list.AtPtr(2))
	assert.Equal(t, obj4, list.AtPtr(3))
	assert.Equal(t, obj5, list.AtPtr(4))

	list.Clear()
	list2.Clear()
	list2.Append(obj1)
	list.Move(list2)
	assert.Equal(t, 1, list.Size())
	assert.Equal(t, 0, list2.Size())

	list2.Append(obj2, obj3)
	list.Move(list2)
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
	for it.SetNext() {
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
	for it.SetPrev() {
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
	for it.SetNext() {
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
	for it.SetPrev() {
		value, _ := it.Value()
		index := it.Index()
		cValue := cSlice[index]

		assert.Equal(t, cValue, value)
	}

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

	for it.SetNext() {
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
	for it.SetPrev() {
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

// ------- Marking elements -------

// MarkAtIndex : mark element at specified index
func (p *XList[T]) MarkAtIndex(index int) {
	xObj := p.goToPosition(index)
	if xObj != nil {
		xObj.mark = true
	}
}

// UnmarkAtIndex : clear mark of element at specified index
func (p *XList[T]) UnmarkAtIndex(index int) {
	xObj := p.goToPosition(index)
	if xObj != nil {
		xObj.mark = false
	}
}

// IsMarkedAtIndex : returns 'true' if element at specified index is marked
func (p *XList[T]) IsMarkedAtIndex(index int) bool {
	xObj := p.goToPosition(index)
	if xObj != nil {
		return xObj.mark
	}

	return false
}

// MarkAll : mark all elements
func (p *XList[T]) MarkAll() {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	for xobj := p.home; xobj != nil; xobj = xobj.next {
		xobj.mark = true
	}
}

// UnmarkAll : clear mark of all elements
func (p *XList[T]) UnmarkAll() {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	for xobj := p.home; xobj != nil; xobj = xobj.next {
		xobj.mark = false
	}
}
