package heap

import (
	"cmp"
	"testing"
)

func verify(t *testing.T, h *Heap[int], i int) {
	t.Helper()

	n := h.Len()
	j1 := 2*i + 1
	j2 := 2*i + 2
	if j1 < n {
		if h.compare(h.buf[j1], h.buf[i]) < 0 {
			t.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d", i, h.buf[i], j1, h.buf[j1])
			return
		}
		verify(t, h, j1)
	}
	if j2 < n {
		if h.compare(h.buf[j2], h.buf[i]) < 0 {
			t.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d", i, h.buf[i], j1, h.buf[j2])
			return
		}
		verify(t, h, j2)
	}
}

func Test(t *testing.T) {
	t.Parallel()

	h := New(cmp.Compare[int])
	verify(t, h, 0)
	h.Grow(20)
	verify(t, h, 0)

	for i := 20; i > 10; i-- {
		h.Push(i)
	}
	verify(t, h, 0)

	for i := 10; i > 0; i-- {
		h.Push(i)
		verify(t, h, 0)
	}

	for i := 1; h.Len() > 0; i++ {
		x, ok := h.Pop()
		if !ok {
			t.Errorf("Pop() = false, want %d", i)
		}
		if i < 20 {
			h.Push(20 + i)
		}
		verify(t, h, 0)
		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func TestPopEmpty(t *testing.T) {
	t.Parallel()

	h := New(cmp.Compare[int])
	_, gotOk := h.Pop()
	const want = false
	if gotOk != want {
		t.Errorf("Pop() on empty heap = %t, want %t", gotOk, want)
	}
}
