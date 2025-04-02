# XList is a classic two-linked list

## Introduction


Xlist is a container that represents a classic doubly linked list, where each element is connected to the previous and next ones. 
This kind of container is efficient for storing and sequential elements processing.

Creation of this container was inspired by rich functionality of NSArray from Apple dev-library.  


## Installation and usage

Install go module:
```shell
go get github.com/DmitryVokhmin/xlist
```

# API description

## Core Methods

### New( ...T )
#### *creates a new empty XList container*
```go
New(objects ...T) *XList[T]
```
Creates a new empty `XList` container. If initial objects are provided, they will be added to the container.

Example:
```go
list := xlist.New[int](1, 2, 3)
fmt.Println(list.Size()) // Output: 3
```


### At( int )  
#### *returns value at specified position*
```go
At(index int) (T, bool)
```
Returns Value and Ok flag: true - value is valid, false - no value.
First element is at 0 position.

Example:
```Go
value, ok := list.At(10)
if ok {
    fmt.Println(value)
} else {
    fmt.Println("No value at position 10")
}
```



### AtPtr( int ) 
#### *returns pointer to a value at specified position*
_(experimental future)_

```go
AtPtr(index int) T
```
Returns a value pointer or nil if no value. First element is at 0 position.
Designed specifically to work with pointers in container.
AtPtr(...) can return 'nil', so no need to return additional validity flag like `At(...)`.

Example:
```Go
value := list.AtPtr(10)
if value != nil {
    fmt.Println(*value)
} else {
    fmt.Println("No value at position 10")
}
```



### IsEmpty() 
#### *Checks if the container is empty or not.*
```go
IsEmpty() bool
```
This function returns 'true' if container is empty

Example:
```Go
if list.IsEmpty() {
    fmt.Println("List is empty")
} else {
    fmt.Println("List is not empty")
}
```



### Size() 
#### *returns number of elements inside container*
```go
Size() int 
```
Example:
```Go
fmt.Println("List size is", list.Size())
```



### LastObject()
#### *returns last object in container*
```Go
LastObject() (T, bool)
```

Example:
```Go
value, ok := list.LastObject()
if ok {
    fmt.Println(value)
} else {
    fmt.Println("List is empty")
}
```



### LastObjectPtr() 
#### *returns last object pointer in container*

```Go
LastObjectPtr() T
```
Example:
```Go
value := list.LastObjectPtr()
if value != nil {
    fmt.Println(*value)
} else {
    fmt.Println("List is empty")
}
```



### Clear()
#### *Clears container content*

```Go
Clear()
```

Example:
```Go
list.Clear()
```



### Set(...T)
#### *set 'objects' to container*

```Go
Set(objects ...T)
```

Set object to the container. In case of empty `objects` objects receiver will be unchanged.

**( !!! ) It resets all the container content and removes old values.**
Consider using Append() if you want to keep old values

Example:
```Go
n := list.Size() // n == 2 (2 elements in container)
list.Set(1, 2, 3, 4, 5)
n := list.Size() // n == 5 
```



### Append(...T)
#### *appends 'objects' to container*

```Go
Append(objects ...T)
```

Appends new values 'objects' to container. In case of empty objects receiver will be unchanged. 

Example:
```Go
n := list.Size() // n == 2 (2 elements in container)
list.Append(1, 2, 3, 4, 5)
n := list.Size() // n == 7 
```



### AppendUnique(...T)
#### *appends unique 'objects' to container (adds elements if they don't exist)*

```Go
AppendUnique(objects ...T)
```

Adds new objects to the container but skips any that already exist within the container

Example:
```Go
list.Set(1, 6)
n := list.Size() // n == 2 (2 elements in container)

list.AppendUnique(1, 2, 3, 4, 5)
n := list.Size() // n == 6
```



### Contains(...T)
#### *checks whether the set of objects is in the container*

```Go
Contains(objects ...T) bool
```

Function checks whether the set of objects (the whole set) is in the container.

Example:

```Go
if list.Contains(1, 2, 3, 4, 5) {
    fmt.Println("All elements are in the container")	
} 
```



### Insert(int, ...T)
#### *inserts 'objects' at position 'pos'*

```Go
Insert(pos int, objects ...T) error
```

Function inserts objects before the position 'pos'.
If the position exceeds the container's size, objects will be appended to the end.

Example:

```Go
if err = list.Insert(8, obj1, obj2, obj3, obj4, obj5); err != nil {
	fmt.Println("Insert error: %w",err)
}
```



### Replace(int, T)
#### *replaces element at specified position*

```Go
Replace(pos int, obj T) error
```

Function replaces element at position 'pos' to 'obj'.
Returns an error if the position is out of the container bounds.

Example:

```Go
if err = list.Replace(8, obj); err != nil {
	fmt.Println("Replace error: %w",err)
}
```



### ReplaceLast(T)
#### *replaces last element in the container*

```Go
ReplaceLast(obj T) error
```
Function replaces the last element in the container with 'obj'. Returns an error if the container is empty.

Example:
```Go
if err = list.ReplaceLast(newObj); err != nil {
    fmt.Println("ReplaceLast error:", err)
}
```



### DeleteAt(int)
#### *deletes element at specified position*
```Go
DeleteAt(pos int) (T, error)
```

Function removes an element at position 'pos' and returns its value. Returns an error if the position is out of bounds or the container is empty.

Example:
```Go
value, err := list.DeleteAt(5)
if err != nil {
    fmt.Println("DeleteAt error:", err)
} else {
    fmt.Println("Deleted value:", value)
}
```



### DeleteLast()
#### *deletes last element from the container*
```Go
DeleteLast() (T, error)
```

Function removes the last element from the container and returns its value. Returns an error if the container is empty.

Example:
```Go
value, err := list.DeleteLast()
if err != nil {
    fmt.Println("DeleteLast error:", err)
} else {
    fmt.Println("Last value deleted:", value)
}
```



### Add( *XList[T] )
#### *combines two lists and returns a new one*
```Go
Add(dList *XList[T]) *XList[T]
```
Function adds 'dList' elements to the receiver's list and returns a new instance. It copies elements from both lists, and the original lists remain unchanged.

Example:
```Go
list1 := xlist.New(1, 2, 3)
list2 := xlist.New(4, 5, 6)
combinedList := list1.Add(list2)

// combinedList now contains [1, 2, 3, 4, 5, 6]
// list1 and list2 remain unchanged
```



### Move( *XList[T] )
#### *moves content from another list to the end of this list*
```Go
Move(dList *XList[T])
```

Function moves content from 'dList' to the end of the receiver. After this operation, 'dList' becomes empty.

Example:
```Go
list1 := xlist.New(1, 2, 3)
list2 := xlist.New(4, 5, 6)
list1.Move(list2)

// list1 now contains [1, 2, 3, 4, 5, 6]
// list2 is now empty
```


### MoveAtPos( int, *XList[T] )
#### *inserts content from another list at specified position*
```Go
MoveAtPos(pos int, dList *XList[T]) error
```

Function inserts (moves) content from 'dList' to receiver at position 'pos'. After this operation, 'dList' becomes empty. Returns an error if the position is out of bounds.

Example:
```Go
list1 := xlist.New(1, 2, 5, 6)
list2 := xlist.New(3, 4)

err := list1.MoveAtPos(2, list2)
if err != nil {
    fmt.Println("MoveAtPos error:", err)
} else {
	
// list1 now contains [1, 2, 3, 4, 5, 6]
// list2 is now empty
}
```



### Copy()
#### *creates a copy of the list*
```Go
Copy() *XList[T]
```
Function returns a complete copy of the list. Depending on the list's configuration, it can be a shallow or deep copy.

Example:
```Go
list := xlist.New(1, 2, 3, 4, 5)

list.DoDeepCopy()
deepCopyList := list.Copy() // deep copy

list.DoShallowCopy()
shallowCopyList := list.Copy() // shallow copy

// deepCopyList contains [1, 2, 3, 4, 5]
// shallowCopyList contains [1, 2, 3, 4, 5]
```



### CopyRange( int, int )
#### *creates a copy of a range of elements*
```go
CopyRange(fromPos int, toPos int) (*XList[T], error)
```
Creates a new container with a copy of elements from position `fromPos` to `toPos` (inclusive). Returns an error if the range is invalid.

Example:
```go
list := xlist.New[int](1, 2, 3, 4, 5)
sublist, err := list.CopyRange(1, 3)
if err != nil {
    fmt.Println("CopyRange error:", err)
} else {
    fmt.Println(sublist.Size()) // Output: 3
    fmt.Println(sublist.AtPtr(0)) // Output: 2
    fmt.Println(sublist.AtPtr(1)) // Output: 3
    fmt.Println(sublist.AtPtr(2)) // Output: 4
}
```



### Swap( int, int )
#### *swaps two elements in the list*
```go
Swap(i, j int) error
```
Swaps two elements at positions `i` and `j`. Returns an error if either index is out of bounds.

Example:
```go
list := xlist.New[int](1, 2, 3, 4)
err := list.Swap(1, 2)
if err != nil {
    fmt.Println("Swap error:", err)
} else {
    fmt.Println(list.AtPtr(1)) // Output: 3
    fmt.Println(list.AtPtr(2)) // Output: 2
}
```


## Marking Methods
Container element marking methods allow working with groups of objects without implementing additional logic.

### MarkAtIndex( int )
#### *marks an element at the specified index*
```go
MarkAtIndex(index int)
```
Marks the element at the specified index. This can be used for custom marking logic.

Example:
```go
list := xlist.New[int](1, 2, 3, 4)

list.MarkAtIndex(2)

fmt.Println(list.IsMarkedAtIndex(2)) // Output: true
```



### UnmarkAtIndex( int )
#### *unmarks an element at the specified index*
```go
UnmarkAtIndex(index int)
```
Clears the mark of the element at the specified index.

Example:
```go
list := xlist.New[int](1, 2, 3, 4)

list.MarkAtIndex(2)
list.UnmarkAtIndex(2)

fmt.Println(list.IsMarkedAtIndex(2)) // Output: false
```



### IsMarkedAtIndex(int)
#### *checks whether an element is marked at the specified index*
```go
IsMarkedAtIndex(index int) bool
```
Returns `true` if the element at the specified index is marked, otherwise `false`.

Example:
```go
list := xlist.New[int](1, 2, 3, 4)

list.MarkAtIndex(2)

fmt.Println(list.IsMarkedAtIndex(2)) // Output: true
fmt.Println(list.IsMarkedAtIndex(1)) // Output: false
```



### MarkAll()
#### *marks all elements in the list*
```go
MarkAll()
```
Marks all elements in the list.

Example:
```go
list := xlist.New[int](1, 2, 3, 4)

list.MarkAll()

fmt.Println(list.IsMarkedAtIndex(0)) // Output: true
fmt.Println(list.IsMarkedAtIndex(1)) // Output: true
fmt.Println(list.IsMarkedAtIndex(2)) // Output: true
fmt.Println(list.IsMarkedAtIndex(3)) // Output: true
```



### UnmarkAll()
#### *unmarks all elements in the list*
```go
UnmarkAll()
```
Clears the mark of all elements in the list.

Example:
```go
list := xlist.New[int](1, 2, 3, 4)

list.MarkAll()
list.UnmarkAll()

fmt.Println(list.IsMarkedAtIndex(0)) // Output: false
fmt.Println(list.IsMarkedAtIndex(1)) // Output: false
fmt.Println(list.IsMarkedAtIndex(2)) // Output: false
fmt.Println(list.IsMarkedAtIndex(3)) // Output: false
```


## Integrations with Go primitives

### Slice()
#### *get a Go slice from collection*

```Go
Slice() []T
```
Functions returns all collection objects as a slice.

Example:
```Go
list := xlist.New[int](1, 2, 3, 4, 5)

slice := list.Slice()

fmt.Println(slice) // Output: [1 2 3 4 5]
```

## Iterator (Sequential traversal of elements)
### Iterator()
#### *creates an iterator for sequential traversal*
```go
Iterator() *Iterator[T]
```
Creates an iterator of the list that supports forward and reverse iteration.
Sets iterator work range from the beginning to the end. 

Example:
```go
list := xlist.New[int](1, 2, 3, 4)

it := list.Iterator()
for it.Next() {
    value, _ := it.Value()
    fmt.Println(value)
}

// Output:
// 1
// 2
// 3
// 4
```

### Iterator( int )
#### *creates an iterator with work interval from the specified index to the last element*

```Go
Iterator(start int) *Iterator[T]
```

Creates an iterator for traversal from a specific index to the last element of the work interval.

Example:

```Go
list := xlist.New[int](1, 2, 3, 4, 5)

it := list.Iterator(2)
for it.Next() {
    value, _ := it.Value()
    fmt.Println(value)
}

// Output:
// 3
// 4
// 5
```


### Iterator( int, int )
#### *creates an iterator for a specific range*
```go
Iterator(start int, finish int) *Iterator[T]
```
Creates an iterator for traversal within a specific range of indices.

Example:
```go
list := xlist.New[int](1, 2, 3, 4, 5)

it := list.Iterator(1, 3)
for it.Next() {
    value, _ := it.Value()
    fmt.Println(value)
}

// Output:
// 2
// 3
// 4
```



### Reset( ...int )
#### *resets the iterator with a new range of work*

```Go
Reset(workRange ...int)
```
Reset - resets the iterator with a new range of work. Accepts 0, 1, or 2 arguments.

If no arguments are provided, the iterator is reset to traverse from the first to the last element of the container. 
If 1 argument is provided, it is treated as the start index of the range. 
If 2 arguments are provided, they define the start and end indexes of the range respectively.

Example:

```Go
list := xlist.New[int](1, 2, 3, 4, 5)

it := list.Iterator()
for it.Next() {
    value, _ := it.Value()
    fmt.Println(value)
}

fmt.Println("Reset the iterator with new work range:")

it.Reset(1, 3)
for it.Next() {
	value, _ := it.Value()
}

// Output:
// 1
// 2
// 3
// 4
// 5
// Reset the iterator with new work range
// 2
// 3
// 4
```



### SetIndex( int ) (T, bool)
#### *sets the iterator to the specified index*

```Go
SetIndex(index int) (T, bool)
```

Set the iterator to the specified index. 
Returns the element at the new index and a boolean indicating whether the returned value was valid.

Example:

```Go
list := xlist.New[int](1, 2, 3, 4, 5)

iter := list.Iterator()

if value, ok := iter.SetIndex(2); ok {
    fmt.Println(value) // Output: 3	
} else {
	fmt.Println("Invalid value")
}

fmt.Println(iter.Value()) // Output: 3
```



### SetFirst()
#### *sets the iterator to the first element of the container or work range*

```Go
SetFirst() (T, bool)
```

Sets the iterator to the first element of the container or work range. 
Returns the first element and a boolean indicating if the returned value is valid.

Example:
```Go
list := xlist.New[int](1, 2, 3, 4, 5)

iter := list.Iterator()

iter.Next()
fmt.Println(iter.Value()) // Output: 2

iter.SetFirst()
fmt.Println(iter.Value()) // Output: 1

iter = list.Iterator(2)
fmt.Println(iter.Value()) // Output: 3

iter.Next()
fmt.Println(iter.Value()) // Output: 4

iter.SetFirst()
fmt.Println(iter.Value()) // Output: 3
```



### SetLast()
#### *sets the iterator to the last element of the container or work range*
```Go
SetLast() (T, bool)
```
Sets the iterator to the last element of the container or work range. Returns the last element and a boolean indicating if the returned value is valid.

Example:
```Go
list := xlist.New[int](1, 2, 3, 4, 5)

iter := list.Iterator()
last, _ := iter.SetLast()
fmt.Println(last) // Output: 5

// Iterator with custom range
iter = list.Iterator(1, 3)
last, _ = iter.SetLast()
fmt.Println(last) // Output: 4
```



### Index()
#### *returns the current index of the iterator*
```Go
Index() int
```
Returns the current index of the iterator. If the iterator hasn't been initialized yet, it returns -1.

Example:
```Go
list := xlist.New[int](1, 2, 3, 4, 5)

iter := list.Iterator()
fmt.Println(iter.Index()) // Output: -1 (not initialized yet)

iter.Next()
fmt.Println(iter.Index()) // Output: 0

iter.Next()
fmt.Println(iter.Index()) // Output: 1
```



### Value()
#### *returns the current value of the iterator*
```Go
Value() (T, bool)
```
Returns the current value of the iterator and a boolean indicating if the value is valid.

Example:
```Go
list := xlist.New[int](1, 2, 3, 4, 5)

iter := list.Iterator()
value, ok := iter.Value() // ok is false, not initialized yet

iter.Next()
value, ok = iter.Value()
fmt.Println(value, ok) // Output: 1 true
```



### Next()
#### *advances the iterator to the next element*
```Go
Next() bool
```
Advances the iterator to the next element. Returns true if successful, false if the end of the range has been reached.

Example:
```Go
list := xlist.New[int](1, 2, 3)

iter := list.Iterator()
for iter.Next() {
    value, _ := iter.Value()
    fmt.Println(value)
}

// Output:
// 1
// 2
// 3
```



### Prev()
#### *moves the iterator to the previous element*
```Go
Prev() bool
```
Moves the iterator to the previous element. Returns true if successful, false if the beginning of the range has been reached.

Example:
```Go
list := xlist.New[int](1, 2, 3, 4, 5)

iter := list.Iterator()
iter.SetLast() // Start from the end

for iter.Prev() {
    value, _ := iter.Value()
    fmt.Println(value)
}

// Output:
// 4
// 3
// 2
// 1
```



### NextValue()
#### *returns the next value and moves the iterator forward*
```Go
NextValue() (T, bool)
```
Returns the next value from the container and moves the iterator forward the iterator. Returns a value and a boolean indicating if the value is valid. Returns false when the end of the range is reached.

Example:
```Go
list := xlist.New[int](1, 2, 3)

iter := list.Iterator()
for {
    value, ok := iter.NextValue()
    if !ok {
        break
    }
    fmt.Println(value)
}

// Output:
// 1
// 2
// 3
```



### PrevValue()
#### *returns the previous value and moves the iterator backward*
```Go
PrevValue() (T, bool)
```
Returns the previous value from the container and moves the iterator backward. Returns a value and a boolean indicating if the value is valid. Returns false when the beginning of the range is reached.

Example:
```Go
list := xlist.New[int](1, 2, 3, 4, 5)

iter := list.Iterator()
iter.SetLast() // Start from the end

for {
    value, ok := iter.PrevValue()
    if !ok {
        break
    }
    fmt.Println(value)
}

// Output:
// 4
// 3
// 2
// 1
```



## Bulk processing methods

### Find( func(int, T) bool )
#### *finds elements that match specific criteria*
```Go
Find(is func(index int, object T) bool) *XList[T]
```
Searches for objects in the list according to criteria defined in the 'is' function and returns a new list with objects that were found.

Example:
```Go
list := xlist.New[int](1, 2, 3, 4, 5)

// Find all even numbers
evenNumbers := list.Find(func(index int, value int) bool {
    return value%2 == 0
})

fmt.Println(evenNumbers.Size()) // Output: 2
fmt.Println(evenNumbers.AtPtr(0)) // Output: 2
fmt.Println(evenNumbers.AtPtr(1)) // Output: 4
```



### Modify( func(int, T) T )
#### *modifies each element in the collection*
```Go
Modify(change func(index int, object T) T) 
```
Modifies each element in the collection by applying the provided function. This is useful when the XList is used in a highly concurrent environment, since each 'change' function's logic performs under the list's internal mutex.

Example:
```Go
list := xlist.New[int](1, 2, 3, 4, 5)

// Multiply each element by 2
list.Modify(func(index int, value int) int {
    return value * 2
})

fmt.Println(list.AtPtr(0)) // Output: 2
fmt.Println(list.AtPtr(1)) // Output: 4
fmt.Println(list.AtPtr(2)) // Output: 6
fmt.Println(list.AtPtr(3)) // Output: 8
fmt.Println(list.AtPtr(4)) // Output: 10
```



### ModifyRev( func(int, T) T )
#### *modifies each element in reverse order*
```Go
ModifyRev(change func(index int, object T) T)
```
Modifies each element in the collection by going in reverse order. The function applies the changes starting from the end of the list and moving towards the beginning.
Useful when current result depends on previous modified values. 

Example:
```Go
list := xlist.New[int](1, 2, 3, 4, 5)

// Add index value in reverse (starting from the end)
list.ModifyRev(func(index int, value int) int {
    return value + index
})

fmt.Println(list.AtPtr(0)) // Output: 1 + 4 = 5
fmt.Println(list.AtPtr(1)) // Output: 2 + 3 = 5
fmt.Println(list.AtPtr(2)) // Output: 3 + 2 = 5
fmt.Println(list.AtPtr(3)) // Output: 4 + 1 = 5
fmt.Println(list.AtPtr(4)) // Output: 5 + 0 = 5
```


## Sorting Methods

### Sort( func(T, T) bool )
#### *sorts the list using a custom comparator*
```go
Sort(compare func(a, b T) bool)
```
Sorts the list in-place using the provided comparator function. The comparator should return `true` if the first argument is less than the second.

Example:
```go
list := xlist.New[int](5, 3, 8, 1, 2)

list.Sort(func(a, b int) bool { return a < b })

fmt.Println(list.AtPtr(0)) // Output: 1
fmt.Println(list.AtPtr(1)) // Output: 2
fmt.Println(list.AtPtr(2)) // Output: 3
fmt.Println(list.AtPtr(3)) // Output: 5
fmt.Println(list.AtPtr(4)) // Output: 8
```

---


### Title
#### *subtitle*
```Go
signature
```
Description

Example:
```Go
Пример использования
```

