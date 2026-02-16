// integrations.go
// Integrations with standard  Go-functional
// Created by Vokhmin D.A. 01.2025

package xlist

// Slice : get all collection objects as a slice
func (p *XList[T]) Slice() []T {
	result := make([]T, 0, p.Size())

	p.mtx.RLock()
	defer p.mtx.RUnlock()

	for xobj := p.home; xobj != nil; xobj = xobj.next {
		result = append(result, *xobj.obj)
	}

	return result
}
