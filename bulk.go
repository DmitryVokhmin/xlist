// bulk.go
// Bulk processing functions
// Created by Vokhmin D.A. 03.2025

package xlist

// Find : looking for objects in list according to criteria defined in 'is' function and
// returns new list with objects that were found.
func (p *XList[T]) Find(is func(index int, object T) bool) *XList[T] {
	lobj := p.home
	newList := &XList[T]{}
	i := 0

	p.mtx.RLock()
	defer p.mtx.RUnlock()

	for lobj != nil {
		if is(i, *lobj.obj) {
			newList.Append(*lobj.obj)
		}
		i++
	}

	return newList
}

// Modify : modifies each element in collection.
// Useful when XList works in highly concurrency mode, since each 'change' func logic performs under internal mutex.
func (p *XList[T]) Modify(change func(index int, object T) T) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	lobj := p.home
	i := 0

	for lobj != nil {
		*lobj.obj = change(i, *lobj.obj)
		lobj = lobj.next
		i++
	}
}

// ModifyRev : modify each element in collection (go in reverse order)
func (p *XList[T]) ModifyRev(change func(index int, object T) T) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	lobj := p.end
	i := p.Size() - 1

	for lobj != nil {
		*lobj.obj = change(i, *lobj.obj)
		lobj = lobj.prev
		i--
	}
}
