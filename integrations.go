// integrations.go
// Integrations with standard  Go-functional
// Created by Vokhmin D.A. 01.2025

package xlist

// Slice : get all collection objects as a slice
func (p *XList[T]) Slice() []T {
	result := make([]T, 0, p.size)

	p.mtx.RLock()
	defer p.mtx.RUnlock()

	xobj := p.home
	for xobj != nil {
		result = append(result, *xobj.obj)
	}

	return result
}
