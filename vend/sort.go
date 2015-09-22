package vend

import (
	"sort"
)

// By is the type of a "less" function that defines the ordering of its arguments.
type By func(p1, p2 *Sale) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(sales []Sale) {
	ss := &saleSorterDesc{
		sales: sales,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ss)
}

// salesSorter joins a By function and a slice of Sales to be sorted.
type saleSorterDesc struct {
	sales []Sale
	by    func(p1, p2 *Sale) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *saleSorterDesc) Len() int {
	return len(s.sales)
}

// Swap is part of sort.Interface.
func (s *saleSorterDesc) Swap(i, j int) {
	s.sales[i], s.sales[j] = s.sales[j], s.sales[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *saleSorterDesc) Less(i, j int) bool {
	return s.by(&s.sales[j], &s.sales[i])
}
