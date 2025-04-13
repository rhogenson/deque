// Package heap implements a priority queue as a min heap backed by a slice.
//
// This can be seen as a replacement for the standard library [container/heap]
// package which was created before generics were a thing.
package heap

import (
	"slices"
)

// Heap is a binary heap backed by a slice.
type Heap[T any] struct {
	buf     []T
	compare func(T, T) int
}

// New creates a new heap with the given comparison function.
func New[T any](compare func(T, T) int) *Heap[T] {
	return &Heap[T]{compare: compare}
}

// Len returns the number of elements in the Heap.
func (h *Heap[T]) Len() int {
	return len(h.buf)
}

// Grow makes space for at least n more elements to be pushed onto the heap
// without reallocating.
func (h *Heap[T]) Grow(n int) {
	h.buf = slices.Grow(h.buf, n)
}

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = h.Len().
func (h *Heap[T]) Push(x T) {
	n := len(h.buf)
	h.buf = append(h.buf, x)
	h.up(n)
}

// Pop removes and returns the minimum element (according to Less) from
// the heap. The complexity is O(log n) where n = h.Len().
func (h *Heap[T]) Pop() (T, bool) {
	if len(h.buf) == 0 {
		var zero T
		return zero, false
	}
	x := h.buf[0]
	last := h.buf[len(h.buf)-1]
	h.buf = h.buf[:len(h.buf)-1]
	if len(h.buf) > 0 {
		h.down(0, last)
	}
	return x, true
}

func (h *Heap[T]) up(j int) {
	x := h.buf[j]
	for {
		i := (j - 1) / 2 // parent
		if i == j || h.compare(x, h.buf[i]) >= 0 {
			break
		}
		h.buf[j] = h.buf[i]
		j = i
	}
	h.buf[j] = x
}

func (h *Heap[T]) down(i int, x T) {
	for {
		j1 := 2*i + 1
		if j1 >= len(h.buf) || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < len(h.buf) && h.compare(h.buf[j2], h.buf[j1]) < 0 {
			j = j2 // = 2*i + 2  // right child
		}
		if h.compare(x, h.buf[j]) <= 0 {
			break
		}
		h.buf[i] = h.buf[j]
		i = j
	}
	h.buf[i] = x
}
