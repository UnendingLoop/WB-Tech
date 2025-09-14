// Package heaper provides a heap-sort for lines from multiple tmp-files
package heaper

import (
	"sortClone/internal/sorter"
)

type FileLine struct {
	Value  string
	FileID int
}

type StrHeap []*FileLine

func (h StrHeap) Len() int { return len(h) }
func (h StrHeap) Less(i, j int) bool { // Переиспользуем UniversalComparator для кучи
	lines := []string{h[i].Value, h[j].Value}
	return sorter.UniversalComparator(lines)(0, 1)
}
func (h StrHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *StrHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*FileLine))
}

func (h *StrHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
