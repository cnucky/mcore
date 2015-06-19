// CopyPasted from golang-source and adjusted for int64
// https://golang.org/src/sort/sort.go?s=6672:6690#L260
//
// TODO: Add unittest as prove it works
package sort

import (
	"sort"
)

// IntSlice attaches the methods of Interface to []int, sorting in increasing order.
type IntSlice []int64

func (p IntSlice) Len() int           { return len(p) }
func (p IntSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p IntSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Sort is a convenience method.
func (p IntSlice) Sort() { sort.Sort(p) }

func Ints(a []int64) { sort.Sort(IntSlice(a)) }