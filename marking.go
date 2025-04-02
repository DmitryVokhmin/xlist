package xlist

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
