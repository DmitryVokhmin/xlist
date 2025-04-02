// xlist.go
// The container is a classic two-linked list
// Created by Vokhmin D.A. 11.2024

package xlist

import (
	"errors"
	"sync"
)

// Sort indexing starts when elements in sorted array are x2 of sortContext.grains
const sortIndexingFromSize = 2

var (
	ErrElementNotFound = errors.New("element not found")
	ErrInvalidIndex    = errors.New("invalid index")
	ErrIsNotAPointer   = errors.New("object is not a pointer")
)

type Compare[T any] interface {
	func(a, b T) bool
}

type XList[T any] struct {
	home *xlistObj[T] // first object
	end  *xlistObj[T] // last object

	size int // counts elements inside container

	mtx sync.RWMutex

	// Work params ----

	// false - make deep copy of containing objects
	// true - make shallow copy (use pointers)
	shallowCopy bool

	// Sort mutex
	sortContext *sortContext[T]
}

// element of bidirectional XList
type xlistObj[T any] struct {
	next *xlistObj[T] // pointer to next element in chain
	prev *xlistObj[T] // pointer to previous element element in chain
	mark bool         // mark element

	obj *T
}

type sortContext[T any] struct {
	changeMtx sync.RWMutex
	cond      *sync.Cond // wait for signal changes done
	canRead   bool

	grains  int
	indexes []*indexPair[T]
}

// indexPair : structure to store a pair of index and object in XList (need for sorting)
type indexPair[T any] struct {
	ix  int
	obj *xlistObj[T]
}

// Iterator : optimal for sequential element passes
type Iterator[T any] struct {
	parent *XList[T]    // parent structure
	index  int          // index
	lobj   *xlistObj[T] // pointer to XList object

	// Allowed range
	start  int
	finish int
}

// New : create new empty XList container
func New[T any](objects ...T) *XList[T] {
	newList := XList[T]{
		mtx: sync.RWMutex{},
	}

	if len(objects) == 0 {
		return &newList
	}

	newList.Set(objects...)

	return &newList
}
