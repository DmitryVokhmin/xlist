// sorting.go
// Sorts the items in the container
// Created by Vokhmin D.A. 03.2025

package xlist

// Sort sorts the list according to the provided comparison function.
//
// Parameters:
//   - compare: A function that compares two elements. Returns true when `a` have to be before `b`, otherwise false.
func (p *XList[T]) Sort(compare func(a, b T) bool) {
	// Performance benchmarks:
	//   - Single thread:  10,000 items ~0.18s, 100,000 items ~28.00s
	//   - Two threads:    10,000 items ~0.15s, 100,000 items ~28.00s
	//   - 16 threads:     10,000 items ~0.06s, 100,000 items ~3.53s, 200,000 items ~12-16s
	//   - 16 threads with indexes: 10,000 items ~0.05s, 100,000 items ~1.04-1.47s,
	//     200,000 items ~2.67-3s, 1,000,000 items ~44.69-55.46s
	// p.ScanSort(compare)
	p.PDQSort(compare)
}
