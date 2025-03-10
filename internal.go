// internal.go
// Internal container methods
// Created by Vokhmin D.A. 01.2025

package xlist

//  ----------------

// goToPosition : go to object at 'pos' position
// returns internal 'xlistObj' struct
func (p *XList[T]) goToPosition(pos int) *xlistObj[T] {
	if pos < 0 || pos > p.size-1 {
		return nil
	}

	xobj := p.home
	i := 0

	for i != pos && xobj.next != nil {
		xobj = xobj.next
		i++
	}

	if i != pos {
		return nil
	}

	return xobj
}
