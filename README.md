# XLIST is a classic two-linked list

## Introduction

---

Xlist is a container that represents a classic doubly linked list, where each element is connected to the previous and next ones. 
This kind of container is efficient for storing and sequential elements processing.

Creation of this container was inspired by rich functionality of NSArray from Apple dev-library.  


## Installation and usage

---
Install go module:
```shell
go get github.com/DmitryVokhmin/xlist
```

# API description



### At(int)  
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

---

### AtPtr(int) 
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

---

### IsEmpty() 
#### *Checks wheither container is empty or not*
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

---

### Size() 
#### *returns number of elements inside container*
```go
Size() int 
```
Example:
```Go
fmt.Println("List size is", list.Size())
```

---

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

---

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

---

### Clear()
#### *Clears container content*

```Go
Clear()
```

Example:
```Go
list.Clear()
```

---

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

---

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

---

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

---

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

---

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

