// internal.go
// Internal container methods
// Created by Vokhmin D.A. 01.2025

package xlist

import "sort"

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

// TODO: Проверить!
// getObjectsAt : returns the xlistObj objects at the specified positions.
func (p *XList[T]) getObjectsAt(pos ...int) []*xlistObj[T] {

	lenpos := len(pos)
	if lenpos == 0 {
		return nil
	}

	if lenpos > 1 {
		sort.Ints(pos) // Sorting positions in ascending order
	}

	var objects []*xlistObj[T]
	xobj := p.home
	i := 0
	ip := 0

	for xobj != nil {
		position := pos[ip]
		if position < 0 || position > p.size-1 { // in case of position outside the range
			ip++
			if ip >= lenpos {
				break
			}
			continue
		}

		if i == position {
			objects = append(objects, xobj)
			ip++
			if ip >= lenpos {
				return objects
			}
		}

		xobj = xobj.next
		i++
	}

	return objects
}
