package GSStringHeap

import (
	"sync"
)

// To be used with the package "container/heap"
type GSStringHeap struct {
	list []string
	cond *sync.Cond
}

func NewGSStringHeap() *GSStringHeap {
	return &GSStringHeap{
		list: make([]string, 0),
		cond: &sync.Cond{L: &sync.Mutex{}},
	}
}

func (h GSStringHeap) Len() int {
	return len(h.list)
}

func (h GSStringHeap) Less(i, j int) bool {
	return h.list[i] < h.list[j]
}

func (h GSStringHeap) Swap(i, j int) {
	h.list[i], h.list[j] = h.list[j], h.list[i]
}

func (h *GSStringHeap) Push(x interface{}) {
	h.cond.L.Lock()
	defer h.cond.L.Unlock()
	h.list = append(h.list, x.(string))
}

func (h *GSStringHeap) Pop() interface{} {
	h.cond.L.Lock()
	defer h.cond.L.Unlock()
	old := h.list
	n := len(old)
	x := old[n-1]
	h.list = old[0 : n-1]
	return x
}
