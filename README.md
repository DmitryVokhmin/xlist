# XList is a classic two-linked list

## Introduction


Xlist is a container representing a classic doubly linked list, where each element is connected to both its previous and next elements. This container is efficient for storing elements and processing them sequentially.

The creation of this container was inspired by the rich functionality of NSArray from Apple’s development library.

It provides support for:
- CRUD operations (Create, Read, Update, and Delete)
- Managing unique container objects
- Bulk processing (modifying container objects using closures)
- Searching for objects
- Highly efficient (multithreaded) sorting operations
- Iterators for efficient sequential operations



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
Clear() *XList[T]
```
Returns self for method chaining; return value can be ignored.

Example:
```Go
list.Clear()

// or with chaining
list.Clear().Append(1, 2, 3)
```



### Set(...T)
#### *set 'objects' to container*

```Go
Set(objects ...T) *XList[T]
```

Set object to the container. In case of empty `objects` objects receiver will be unchanged.
Returns self for method chaining; return value can be ignored.

**( !!! ) It resets all the container content and removes old values.**
Consider using Append() if you want to keep old values

Example:
```Go
n := list.Size() // n == 2 (2 elements in container)
list.Set(1, 2, 3, 4, 5)
n := list.Size() // n == 5

// or with chaining
list.Set(1, 2, 3).Append(4, 5)
```



### Append(...T)
#### *appends 'objects' to container*

```Go
Append(objects ...T) *XList[T]
```

Appends new values 'objects' to container. In case of empty objects receiver will be unchanged.
Returns self for method chaining; return value can be ignored.

Example:
```Go
n := list.Size() // n == 2 (2 elements in container)
list.Append(1, 2, 3, 4, 5)
n := list.Size() // n == 7

// or with chaining
list.Append(1, 2).Append(3, 4)
```



### AppendUnique(...T)
#### *appends unique 'objects' to container (adds elements if they don't exist)*

```Go
AppendUnique(objects ...T) *XList[T]
```

Adds new objects to the container but skips any that already exist within the container.
Returns self for method chaining; return value can be ignored.

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



### ContainsSome(...T)
#### *checks whether any of objects is in the container*

```Go
ContainsSome(objects ...T) bool
```

Function checks whether any of the provided objects is in the container. Returns `true` if at least one object is found.

Example:

```Go
list := xlist.New[int](1, 2, 3)

if list.ContainsSome(5, 6, 2) {
    fmt.Println("At least one element is in the container") // This will print
}

if !list.ContainsSome(10, 20, 30) {
    fmt.Println("None of the elements are in the container") // This will print
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



### AppendList( *XList[T] )
#### *adds elements from another list to the end of this list*
```Go
AppendList(dList *XList[T]) *XList[T]
```
Function adds copies of elements from 'dList' to the end of the receiver's list. The source list 'dList' remains unchanged.
Returns self for method chaining; return value can be ignored.

Example:
```Go
list1 := xlist.New(1, 2, 3)
list2 := xlist.New(4, 5, 6)
list1.AppendList(list2)

// list1 now contains [1, 2, 3, 4, 5, 6]
// list2 remains unchanged [4, 5, 6]
```



### Splice( *XList[T] )
#### *moves content from another list to the end of this list*
```Go
Splice(dList *XList[T]) *XList[T]
```

Function moves content from 'dList' to the end of the receiver. After this operation, 'dList' becomes empty.
Returns self for method chaining; return value can be ignored.

Example:
```Go
list1 := xlist.New(1, 2, 3)
list2 := xlist.New(4, 5, 6)
list1.Splice(list2)

// list1 now contains [1, 2, 3, 4, 5, 6]
// list2 is now empty
```


### SpliceAtPos( int, *XList[T] )
#### *inserts content from another list at specified position*
```Go
SpliceAtPos(pos int, dList *XList[T]) error
```

Function inserts (moves) content from 'dList' to receiver at position 'pos'. After this operation, 'dList' becomes empty. Returns an error if the position is out of bounds.

Example:
```Go
list1 := xlist.New(1, 2, 5, 6)
list2 := xlist.New(3, 4)

err := list1.SpliceAtPos(2, list2)
if err != nil {
    fmt.Println("SpliceAtPos error:", err)
} else {
    // list1 now contains [1, 2, 3, 4, 5, 6]
    // list2 is now empty
}
```



### Copy()
#### *creates a shallow copy of the list*
```Go
Copy() *XList[T]
```
Function returns a shallow copy of the list. For deep copying, use `DeepCopy()` method.

Example:
```Go
list := xlist.New(1, 2, 3, 4, 5)

copyList := list.Copy()

// copyList contains [1, 2, 3, 4, 5]
// Modifying copyList doesn't affect original list
```



### CopyRange( int, int )
#### *creates a copy of a range of elements*
```go
CopyRange(fromPos int, toPos int) (*XList[T], error)
```
Creates a new container with a shallow copy of elements from position `fromPos` to `toPos` (inclusive). Returns an error if the range is invalid.

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



### DeepCopy( func(T) T )
#### *creates a deep copy of the list*
```go
DeepCopy(deepCopyFn func(T) T) *XList[T]
```
Creates a new container with deep copies of all elements. You must provide a closure `deepCopyFn` that knows how to make a deep copy of type T.

Example:
```go
type Person struct {
    Name string
    Age  int
}

list := xlist.New[*Person](
    &Person{"Alice", 30},
    &Person{"Bob", 25},
)

// Deep copy with custom copy function
copyList := list.DeepCopy(func(p *Person) *Person {
    return &Person{Name: p.Name, Age: p.Age}
})

// Modifying the copy doesn't affect the original
copyList.AtPtr(0).Name = "Charlie"
fmt.Println(list.AtPtr(0).Name)     // Output: Alice
fmt.Println(copyList.AtPtr(0).Name) // Output: Charlie
```



### DeepCopyRange( int, int, func(T) T )
#### *creates a deep copy of a range of elements*
```go
DeepCopyRange(fromPos int, toPos int, deepCopyFn func(T) T) (*XList[T], error)
```
Creates a new container with deep copies of elements from position `fromPos` to `toPos` (inclusive). You must provide a closure `deepCopyFn` that knows how to make a deep copy of type T. Returns an error if the range is invalid or if `deepCopyFn` is nil.

Example:
```go
type Person struct {
    Name string
    Age  int
}

list := xlist.New[*Person](
    &Person{"Alice", 30},
    &Person{"Bob", 25},
    &Person{"Charlie", 35},
)

sublist, err := list.DeepCopyRange(0, 1, func(p *Person) *Person {
    return &Person{Name: p.Name, Age: p.Age}
})
if err != nil {
    fmt.Println("DeepCopyRange error:", err)
} else {
    fmt.Println(sublist.Size()) // Output: 2
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


## Range Iterators (Go 1.23+)

These methods leverage Go's new `iter` package (available from Go 1.23) to provide powerful and flexible ways to iterate over the list. They support functional programming patterns like filtering, mapping, and chaining.

### All
#### *returns a forward iterator over the list with index and value*
```go
All(opt ...func(*RangeOptions)) iter.Seq2[int, T]
```
Returns a forward iterator (`iter.Seq2`) that yields both the index and the value for each element. The iteration can be constrained using `WithPos(int)` to specify a starting index and `WithCount(int)` to limit the number of elements.

Example:
```go
list := xlist.New[string]("a", "b", "c", "d", "e")

// Iterate over all elements
for i, v := range list.All() {
    fmt.Printf("Index: %d, Value: %s\n", i, v)
}
// Output:
// Index: 0, Value: a
// Index: 1, Value: b
// ...

// Iterate starting from index 2 with a count of 2
for i, v := range list.All(xlist.WithPos(2), xlist.WithCount(2)) {
    fmt.Printf("Index: %d, Value: %s\n", i, v)
}
// Output:
// Index: 2, Value: c
// Index: 3, Value: d
```

### Backward
#### *returns a reverse iterator over the list with index and value*
```go
Backward(opt ...func(*RangeOptions)) iter.Seq2[int, T]
```
Returns a reverse iterator that yields the index and value for each element, starting from the end of the list and moving to the beginning. It also supports `WithPos` and `WithCount` options.

Example:
```go
list := xlist.New[string]("a", "b", "c", "d", "e")

// Iterate backward over all elements
for i, v := range list.Backward() {
    fmt.Printf("Index: %d, Value: %s\n", i, v)
}
// Output:
// Index: 4, Value: e
// Index: 3, Value: d
// ...

// Iterate backward starting from index 3 with a count of 2
for i, v := range list.Backward(xlist.WithPos(3), xlist.WithCount(2)) {
	fmt.Printf("Index: %d, Value: %s\n", i, v)
}
// Output:
// Index: 3, Value: d
// Index: 2, Value: c
```

### Values
#### *returns a forward iterator of values only*
```go
Values(opt ...func(*RangeOptions)) iter.Seq[T]
```
A convenience method that returns a forward iterator (`iter.Seq`) yielding only the values, without indices. It's equivalent to `ToValues(list.All(...))`.

Example:
```go
list := xlist.New[int](10, 20, 30)

for v := range list.Values() {
    fmt.Println(v)
}
// Output:
// 10
// 20
// 30
```

### ValuesBackward
#### *returns a reverse iterator of values only*
```go
ValuesBackward(opt ...func(*RangeOptions)) iter.Seq[T]
```
A convenience method that returns a reverse iterator yielding only values. It's equivalent to `ToValues(list.Backward(...))`.

Example:
```go
list := xlist.New[int](10, 20, 30)

for v := range list.ValuesBackward() {
    fmt.Println(v)
}
// Output:
// 30
// 20
// 10
```

### ToValues
#### *converts an iterator with index and value to one with values only*
```go
ToValues[T any](seq2 iter.Seq2[int, T]) iter.Seq[T]
```
A helper function that transforms an `iter.Seq2[int, T]` (yielding index and value) into an `iter.Seq[T]` (yielding only values).

Example:
```go
list := xlist.New[int](1, 2, 3)
valuesIterator := xlist.ToValues(list.All())

for v := range valuesIterator {
    fmt.Println(v)
}
// Output:
// 1
// 2
// 3
```

### Filter
#### *creates an iterator that yields only elements matching a predicate*
```go
Filter[T any](it2 iter.Seq2[int, T], is_ok func(int, T) bool) iter.Seq2[int, T]
```
An iterator adapter that takes an `iter.Seq2` and a predicate function. It returns a new `iter.Seq2` that only yields elements for which the predicate returns `true`.

Example:
```go
list := xlist.New[int](1, 2, 3, 4, 5, 6)
evenNumbers := xlist.Filter(list.All(), func(i int, v int) bool {
    return v%2 == 0
})

for i, v := range evenNumbers {
    fmt.Printf("Found even number at index %d: %d\n", i, v)
}
// Output:
// Found even number at index 1: 2
// Found even number at index 3: 4
// Found even number at index 5: 6
```

### TakeWhile
#### *creates an iterator that yields elements while a predicate is true*
```go
TakeWhile[T any](it2 iter.Seq2[int, T], is_ok func(int, T) bool) iter.Seq2[int, T]
```
An iterator adapter that yields elements from the input sequence as long as the predicate returns `true`. The iteration stops permanently after the first time the predicate returns `false`.

Example:
```go
list := xlist.New[int](2, 4, 5, 6, 7)
initialEvens := xlist.TakeWhile(list.All(), func(_ int, v int) bool {
    return v%2 == 0
})

for _, v := range initialEvens {
    fmt.Println(v)
}
// Output:
// 2
// 4
// (Stops at 5, because it's not even)
```

### SkipWhile
#### *creates an iterator that skips elements while a predicate is true*
```go
SkipWhile[T any](it2 iter.Seq2[int, T], is_ok func(int, T) bool) iter.Seq2[int, T]
```
An iterator adapter that skips elements from the input sequence as long as the predicate returns `true`. Once the predicate returns `false`, it starts yielding all subsequent elements.

Example:
```go
list := xlist.New[int](2, 4, 5, 6, 7)
tail := xlist.SkipWhile(list.All(), func(_ int, v int) bool {
    return v%2 == 0
})

for _, v := range tail {
    fmt.Println(v)
}
// Output:
// 5
// 6
// 7
// (Skips 2 and 4, starts yielding from 5)
```

### Map
#### *creates an iterator that transforms each element in a sequence*
```go
Map[T, V any](seq iter.Seq[T], transform func(T) V) iter.Seq[V]
```
An iterator adapter that applies a `transform` function to each element of an input `iter.Seq`, producing a new `iter.Seq` with the transformed values. Note this works on value-only iterators.

Example:
```go
list := xlist.New[int](1, 2, 3)
asStrings := xlist.Map(list.Values(), func(v int) string {
    return fmt.Sprintf("Value: %d", v)
})

for s := range asStrings {
    fmt.Println(s)
}
// Output:
// Value: 1
// Value: 2
// Value: 3
```

### AnyMatch
#### *checks if any element in a sequence matches a predicate*
```go
AnyMatch[T any](it2 iter.Seq2[int, T], is_ok func(int, T) bool) bool
```
A terminal operation that consumes an iterator and returns `true` if at least one element satisfies the predicate `is_ok`. The iteration stops as soon as a match is found.

Example:
```go
list := xlist.New[int](1, 3, 4, 5)
hasEvenNumber := xlist.AnyMatch(list.All(), func(_ int, v int) bool {
    return v%2 == 0
})

if hasEvenNumber {
    fmt.Println("The list contains at least one even number.") // This will print
}
```

### AllMatch
#### *checks if all elements in a sequence match a predicate*
```go
AllMatch[T any](it2 iter.Seq2[int, T], is_ok func(int, T) bool) bool
```
A terminal operation that consumes an iterator and returns `true` if all elements satisfy the predicate `is_ok`. The iteration stops as soon as an element fails the predicate.

Example:
```go
list := xlist.New[int](2, 4, 6)
allAreEven := xlist.AllMatch(list.All(), func(_ int, v int) bool {
    return v%2 == 0
})

if allAreEven {
    fmt.Println("All numbers in the list are even.") // This will print
}
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
Modify(change func(index int, object T) T) *XList[T]
```
Modifies each element in the collection by applying the provided function. This is useful when the XList is used in a highly concurrent environment, since each 'change' function's logic performs under the list's internal mutex.
Returns self for method chaining; return value can be ignored.

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
ModifyRev(change func(index int, object T) T) *XList[T]
```
Modifies each element in the collection by going in reverse order. The function applies the changes starting from the end of the list and moving towards the beginning.
Useful when current result depends on previous modified values.
Returns self for method chaining; return value can be ignored.

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
