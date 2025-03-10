# XLIST is a classic two-linked list

## Introduction

---

Xlist is a container that represents a classic doubly linked list, where each element is connected to the previous and next ones. 
This kind of container is efficient for storing and sequential elements processing.

Creation of this container was inspired by rich functionality of NSArray from Apple dev-library.  


## Installation and usage

---

_TODO_

## API description

---

### At - returns value at specified position.
```Go
At(index int) (T, bool)
```
Returns Value and Ok flag: true - value is valid, false - no value


### AtPtr - returns pointer to a value at specified position.
_(experimental future)_


```Go
AtPtr(index int) T
```
Designed specifically to work with pointers in container.
AtPtr(...) can return 'nil', so no need to return additional validity flag like `At(...)`.
That makes the code easier.

### IsEmpty - returns 'true' if container is empty
```Go
IsEmpty() bool
```

### Size() - returns number of elements inside container
```Go
Size() int 
```

