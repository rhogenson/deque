// Package deque implements a double-ended queue (deque) implemented with a
// slice-backed ring buffer.
//
// This queue has O(1) amortized inserts and removals from both ends of the
// container. It also has O(1) indexing like a vector.
//
// The "default" usage of this type as a queue is to use [Deque.PushBack] to add
// to the queue, and [Deque.PopFront] to remove from the queue. Iterating over
// Deque goes front to back.
//
// The core implementation is "ported" (stolen) from Rust's VecDeque.
package deque

import (
	"fmt"
	"iter"
	"slices"
	"strings"
)

// Deque is a double-ended queue. The zero value is ready for use.
type Deque[T any] struct {
	head int
	buf  []T
}

// WithCapacity allocates a deque with the given capacity.
func WithCapacity[T any](cap int) *Deque[T] {
	return &Deque[T]{buf: make([]T, 0, cap)}
}

// From creates a new queue using the given slice as the backing buffer.
func From[S ~[]T, T any](slice S) *Deque[T] {
	return &Deque[T]{buf: slice}
}

func (q *Deque[T]) wrapAdd(i, addend int) int {
	i += addend
	if i >= cap(q.buf) {
		return i - cap(q.buf)
	}
	return i
}

func (q *Deque[T]) toPhysicalIdx(i int) int {
	return q.wrapAdd(q.head, i)
}

// At returns the item at position i. At panics if i < 0 or i >= q.Len().
func (q *Deque[T]) At(i int) T {
	if !(0 <= i && i < len(q.buf)) {
		panic(fmt.Sprintf("index out of range [%d] with length %d", i, len(q.buf)))
	}
	return q.buf[:cap(q.buf)][q.toPhysicalIdx(i)]
}

// Cap returns the number of elements the deque can hold without reallocating.
func (q *Deque[T]) Cap() int {
	return cap(q.buf)
}

// Len returns the number of elements in the deque.
func (q *Deque[T]) Len() int {
	return len(q.buf)
}

// PopFront removes and returns the item at index 0 if the deque is non-empty.
func (q *Deque[T]) PopFront() (T, bool) {
	if len(q.buf) == 0 {
		var zero T
		return zero, false
	}
	oldHead := q.head
	q.head = q.toPhysicalIdx(1)
	q.buf = q.buf[:len(q.buf)-1]
	return q.buf[:cap(q.buf)][oldHead], true
}

// PopBack removes and returns the last item in the deque if it is non-empty.
func (q *Deque[T]) PopBack() (T, bool) {
	if len(q.buf) == 0 {
		var zero T
		return zero, false
	}
	q.buf = q.buf[:len(q.buf)-1]
	return q.buf[:cap(q.buf)][q.toPhysicalIdx(len(q.buf))], true
}

// PushFront prepends the given items to the front of the deque.
func (q *Deque[T]) PushFront(values ...T) {
	q.Grow(len(values))
	q.buf = q.buf[:len(q.buf)+len(values)]
	if q.head >= len(values) {
		newHead := q.head - len(values)
		copy(q.buf[newHead:q.head], values)
		q.head = newHead
	} else {
		tailLen := len(values) - q.head
		copy(q.buf[:q.head], values[tailLen:])
		copy(q.buf[cap(q.buf)-tailLen:cap(q.buf)], values[:tailLen])
		q.head = cap(q.buf) - tailLen
	}
}

// PushBack appends the given items to the back of the deque.
func (q *Deque[T]) PushBack(values ...T) {
	q.Grow(len(values))
	endIdx := q.wrapAdd(q.head, len(q.buf))
	if len(values) <= cap(q.buf)-endIdx {
		copy(q.buf[endIdx:endIdx+len(values)], values)
	} else {
		headLen := cap(q.buf) - endIdx
		copy(q.buf[endIdx:cap(q.buf)], values[:headLen])
		copy(q.buf[:len(values)-headLen], values[headLen:])
	}
	q.buf = q.buf[:len(q.buf)+len(values)]
}

// Reset empties the deque, retaining the underlying storage for use by
// future pushes.
func (q *Deque[T]) Reset() {
	q.buf = q.buf[:0]
}

// Grow makes space for at least n more elements to be inserted in the given
// deque without reallocation.
func (q *Deque[T]) Grow(n int) {
	if n <= cap(q.buf)-len(q.buf) {
		return
	}

	oldCap := cap(q.buf)
	q.buf = slices.Grow(q.buf, n)
	newCap := cap(q.buf)

	// Move the shortest contiguous section of the ring buffer
	//
	// H := head
	// L := last element (`self.to_physical_idx(self.len - 1)`)
	//
	//    H             L
	//   [o o o o o o o o ]
	//    H             L
	// A [o o o o o o o o . . . . . . . . ]
	//        L H
	//   [o o o o o o o o ]
	//          H             L
	// B [. . . o o o o o o o o . . . . . ]
	//              L H
	//   [o o o o o o o o ]
	//              L                 H
	// C [o o o o o o . . . . . . . . o o ]

	if q.head <= oldCap-len(q.buf) {
		// A
		return
	}
	headLen := oldCap - q.head
	tailLen := len(q.buf) - headLen
	if headLen > tailLen && newCap-oldCap >= tailLen {
		// B
		copy(q.buf[oldCap:oldCap+tailLen], q.buf[:tailLen])
		return
	}
	// C
	newHead := newCap - headLen
	copy(q.buf[newHead:newHead+headLen], q.buf[q.head:q.head+headLen])
	q.head = newHead
}

// All returns an iterator over the elements in the deque. It does not pop
// any elements.
func (q *Deque[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		// Don't use range over int in case the length changes while
		// we're iterating
		for i := 0; i < len(q.buf); i++ {
			if !yield(i, q.buf[:cap(q.buf)][q.toPhysicalIdx(i)]) {
				return
			}
		}
	}
}

// PopAll empties the deque and returns an iterator over the popped elements.
// It's not safe to modify the deque while iterating using PopAll.
func (q *Deque[T]) PopAll() iter.Seq[T] {
	n := len(q.buf)
	q.buf = q.buf[:0]
	return func(yield func(T) bool) {
		for i := range n {
			if !yield(q.buf[:cap(q.buf)][q.toPhysicalIdx(i)]) {
				return
			}
		}
	}
}

// String displays the deque as a string, using fmt.Sprint to show each element.
func (q *Deque[T]) String() string {
	buf := new(strings.Builder)
	buf.WriteString("[")
	for i := range len(q.buf) {
		if i > 0 {
			buf.WriteString(" ")
		}
		fmt.Fprint(buf, q.buf[:cap(q.buf)][q.toPhysicalIdx(i)])
	}
	buf.WriteString("]")
	return buf.String()
}
