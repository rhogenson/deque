// Package vecdeque implements a double-ended queue (deque) implemented with a
// growable ring buffer.
//
// This queue has O(1) amortized inserts and removals from both ends of the
// container. It also has O(1) indexing like a vector.
package vecdeque

// DQ is a double-ended queue. The zero value is ready for use.
type DQ[T any] struct {
	head int
	buf  []T
}

// WithCapacity allocates a deque with the given capacity.
func WithCapacity[T any](cap int) *DQ[T] {
	return &DQ[T]{buf: make([]T, 0, cap)}
}

// From creates a new queue using the given slice as the backing buffer.
func From[S ~[]T, T any](slice S) *DQ[T] {
	return &DQ[T]{buf: slice}
}

func (q *DQ[T]) wrapAdd(i, addend int) int {
	i += addend
	if i >= cap(q.buf) {
		return i - cap(q.buf)
	}
	return i
}

func (q *DQ[T]) toPhysicalIdx(i int) int {
	return q.wrapAdd(q.head, i)
}

// Get returns the item at position i.
func (q *DQ[T]) Get(i int) T {
	return q.buf[:cap(q.buf)][q.toPhysicalIdx(i)]
}

// Cap returns the number of elements the deque can hold without reallocating.
func (q *DQ[T]) Cap() int {
	return cap(q.buf)
}

// Len returns the number of elements in the deque.
func (q *DQ[T]) Len() int {
	return len(q.buf)
}

// PopFront removes and returns the item at index 0 if the deque is non-empty.
func (q *DQ[T]) PopFront() (T, bool) {
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
func (q *DQ[T]) PopBack() (T, bool) {
	if len(q.buf) == 0 {
		var zero T
		return zero, false
	}
	q.buf = q.buf[:len(q.buf)-1]
	return q.Get(len(q.buf)), true
}

// PushFront prepends the given items to the front of the deque.
func (q *DQ[T]) PushFront(values ...T) {
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
func (q *DQ[T]) PushBack(values ...T) {
	q.Grow(len(values))
	endIdx := q.wrapAdd(q.head, len(q.buf))
	if len(values) < cap(q.buf)-endIdx {
		copy(q.buf[endIdx:endIdx+len(values)], values)
	} else {
		headLen := cap(q.buf) - endIdx
		copy(q.buf[endIdx:cap(q.buf)], values[:headLen])
		copy(q.buf[:len(values)-headLen], values[headLen:])
	}
	q.buf = q.buf[:len(q.buf)+len(values)]
}

// Grow makes space for at least n more elements to be inserted in the given
// deque without reallocation.
func (q *DQ[T]) Grow(n int) {
	n -= cap(q.buf) - len(q.buf)
	if n <= 0 {
		return
	}

	oldCap := cap(q.buf)
	q.buf = append(q.buf[:cap(q.buf)], make([]T, n)...)[:len(q.buf)]
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
func (q *DQ[T]) All() func(func(int, T) bool) {
	return func(yield func(int, T) bool) {
		for i := range q.Len() {
			if !yield(i, q.Get(i)) {
				return
			}
		}
	}
}

// PopAll returns an iterator that consumes all the values in the deque, leaving
// it empty.
func (q *DQ[T]) PopAll() func(func(T) bool) {
	return func(yield func(T) bool) {
		for val, ok := q.PopFront(); ok; val, ok = q.PopFront() {
			if !yield(val) {
				return
			}
		}
	}
}
