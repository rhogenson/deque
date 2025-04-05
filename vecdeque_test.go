package vecdeque

import (
	"slices"
	"testing"
)

func TestWithCapacity(t *testing.T) {
	t.Parallel()

	const cap = 10
	q := WithCapacity[int](cap)
	for i := range cap {
		q.PushBack(i)
	}
	if got := q.Cap(); got != cap {
		t.Errorf("Cap() = %d, want %d", got, cap)
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	q := new(DQ[int])
	for i := range 10 {
		q.PushBack(i)
	}
	for i := range 3 {
		if got := q.Get(i); got != i {
			t.Errorf("Get(%d) = %d, want %d", i, got, i)
		}
	}
}

func TestPopFront(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		desc         string
		in           []int
		wantOk       bool
		wantVal      int
		wantContents []int
	}{{
		desc:         "PopVal",
		in:           []int{1, 2, 3},
		wantOk:       true,
		wantVal:      1,
		wantContents: []int{2, 3},
	}, {
		desc:   "PopNone",
		in:     nil,
		wantOk: false,
	}} {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			q := From(tc.in)
			got, ok := q.PopFront()
			if ok != tc.wantOk {
				t.Errorf("%d: PopFront() returned ok = %t, want %t", tc.in, ok, tc.wantOk)
			}
			if got != tc.wantVal {
				t.Errorf("%d: PopFront() = %d, want %d", tc.in, got, tc.wantVal)
			}
			gotContents := make([]int, q.Len())
			for i, x := range q.All() {
				gotContents[i] = x
			}
			if !slices.Equal(gotContents, tc.wantContents) {
				t.Errorf("%d: Contents after PopFront are %d, want %d", tc.in, gotContents, tc.wantContents)
			}
		})
	}
}

func TestPopBack(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		desc         string
		in           []int
		wantOk       bool
		wantVal      int
		wantContents []int
	}{{
		desc:         "PopVal",
		in:           []int{1, 2, 3},
		wantOk:       true,
		wantVal:      3,
		wantContents: []int{1, 2},
	}, {
		desc:   "PopNone",
		in:     nil,
		wantOk: false,
	}} {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			q := From(tc.in)
			got, ok := q.PopBack()
			if ok != tc.wantOk {
				t.Errorf("%d: PopBack() returned ok = %t, want %t", tc.in, ok, tc.wantOk)
			}
			if got != tc.wantVal {
				t.Errorf("%d: PopBack() = %d, want %d", tc.in, got, tc.wantVal)
			}
			gotContents := make([]int, q.Len())
			for i, x := range q.All() {
				gotContents[i] = x
			}
			if !slices.Equal(gotContents, tc.wantContents) {
				t.Errorf("%d: Contents after PopBack are %d, want %d", tc.in, gotContents, tc.wantContents)
			}
		})
	}
}

func TestPushFront(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		desc        string
		prevContent []int
		push        []int
		want        []int
	}{{
		desc:        "PushNil",
		prevContent: nil,
		push:        []int{1},
		want:        []int{1},
	}, {
		desc:        "PushExisting",
		prevContent: []int{1, 2, 3},
		push:        []int{4, 5, 6},
		want:        []int{4, 5, 6, 1, 2, 3},
	}} {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			q := From(tc.prevContent)
			q.PushFront(tc.push...)
			got := make([]int, q.Len())
			for i, x := range q.All() {
				got[i] = x
			}
			if !slices.Equal(got, tc.want) {
				t.Errorf("%d: PushFront(%d) = %d, want %d", tc.prevContent, tc.push, got, tc.want)
			}
		})
	}
}

func TestPushBack(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		desc        string
		prevContent []int
		push        []int
		want        []int
	}{{
		desc:        "PushNil",
		prevContent: nil,
		push:        []int{1},
		want:        []int{1},
	}, {
		desc:        "PushExisting",
		prevContent: []int{1, 2, 3},
		push:        []int{4, 5, 6},
		want:        []int{1, 2, 3, 4, 5, 6},
	}} {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			q := From(tc.prevContent)
			q.PushBack(tc.push...)
			got := make([]int, q.Len())
			for i, x := range q.All() {
				got[i] = x
			}
			if !slices.Equal(got, tc.want) {
				t.Errorf("%d: PushBack(%d) = %d, want %d", tc.prevContent, tc.push, got, tc.want)
			}
		})
	}
}

func TestPopFrontPushBackB(t *testing.T) {
	t.Parallel()

	q := From([]int{1, 2, 3})
	q.PopFront()
	q.PushBack(4)
	q.PushBack(5)
	got := make([]int, q.Len())
	for i, x := range q.All() {
		got[i] = x
	}
	want := []int{2, 3, 4, 5}
	if !slices.Equal(got, want) {
		t.Errorf("Contents = %d, want %d", got, want)
	}
}

func TestPopFrontPushBackC(t *testing.T) {
	t.Parallel()

	q := From([]int{1, 2, 3})
	q.PopFront()
	q.PopFront()
	q.PushBack(4)
	q.PushBack(5)
	q.PushBack(6)
	got := make([]int, q.Len())
	for i, x := range q.All() {
		got[i] = x
	}
	want := []int{3, 4, 5, 6}
	if !slices.Equal(got, want) {
		t.Errorf("Contents = %d, want %d", got, want)
	}
}

func TestPopFrontPushFront(t *testing.T) {
	t.Parallel()

	q := From([]int{1, 2, 3})
	q.PopFront()
	q.PopFront()
	q.PushFront(4, 5)
	if got, want := q.Cap(), 3; got != want {
		t.Errorf("Cap() = %d, want %d", got, want)
	}
	got := make([]int, q.Len())
	for i, x := range q.All() {
		got[i] = x
	}
	want := []int{4, 5, 3}
	if !slices.Equal(got, want) {
		t.Errorf("Contents = %d, want %d", got, want)
	}
}

func TestPopAll(t *testing.T) {
	q := From([]int{1, 2, 3})
	got := make([]int, 0, q.Len())
	for x := range q.PopAll() {
		got = append(got, x)
	}
	if got, want := q.Len(), 0; got != want {
		t.Errorf("Len() = %d, want %d", got, want)
	}
	want := []int{1, 2, 3}
	if !slices.Equal(got, want) {
		t.Errorf("PopAll() returned values %d, want %d", got, want)
	}
}
