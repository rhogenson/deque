package deque_test

import (
	"fmt"

	"github.com/rhogenson/deque"
)

func ExampleDeque() {
	q := new(deque.Deque[int])
	for i := range 10 {
		q.PushBack(i)
	}
	for range 3 {
		q.PopFront()
	}
	fmt.Println(q)

	// Output:
	// [3 4 5 6 7 8 9]
}

func ExampleWithCapacity() {
	q := deque.WithCapacity[int](10)
	for i := range 100 {
		if q.Len() == q.Cap() {
			q.PopFront()
		}
		q.PushBack(i)
	}
	fmt.Println(q)

	// Output:
	// [90 91 92 93 94 95 96 97 98 99]
}

func ExampleFrom() {
	q := deque.From([]int{1, 2, 3, 4, 5})
	fmt.Println(q.PopFront())

	// Output:
	// 1 true
}

func ExampleDeque_At() {
	q := deque.From([]int{1, 2, 3, 4, 5})
	fmt.Println(q.At(3))

	// Output:
	// 4
}

func ExampleDeque_Cap() {
	q := deque.WithCapacity[int](10)
	q.PushBack(1, 2, 3, 4, 5)
	fmt.Println(q.Cap())

	// Output:
	// 10
}

func ExampleDeque_Len() {
	q := new(deque.Deque[int])
	q.PushBack(1, 2, 3, 4, 5)
	fmt.Println(q.Len())

	// Output:
	// 5
}

func ExampleDeque_PopFront() {
	q := deque.From([]int{1, 2, 3, 4, 5})
	for range 3 {
		q.PopFront()
	}
	fmt.Println(q)

	// Output:
	// [4 5]
}

func ExampleDeque_PopBack() {
	q := deque.From([]int{1, 2, 3, 4, 5})
	for range 3 {
		q.PopBack()
	}
	fmt.Println(q)

	// Output:
	// [1 2]
}

func ExampleDeque_PushFront() {
	q := deque.From([]int{6, 7, 8, 9, 10})
	q.PushFront(1, 2, 3, 4, 5)
	fmt.Println(q)

	// Output:
	// [1 2 3 4 5 6 7 8 9 10]
}

func ExampleDeque_PushBack() {
	q := deque.From([]int{1, 2, 3, 4, 5})
	q.PushBack(6, 7, 8, 9, 10)
	fmt.Println(q)

	// Output:
	// [1 2 3 4 5 6 7 8 9 10]
}

func ExampleDeque_Reset() {
	q := deque.From([]int{1, 2, 3, 4, 5})
	q.Reset()
	fmt.Println(q.Cap())

	// Output:
	// 5
}

func ExampleDeque_Grow() {
	q := new(deque.Deque[int])
	q.Grow(5)
	// PushBack will not allocate:
	q.PushBack(1, 2, 3, 4, 5)
}

func ExampleDeque_All() {
	q := new(deque.Deque[int])
	q.PushBack(1, 2, 3, 4, 5)
	q.PopFront()
	for _, x := range q.All() {
		fmt.Println(x)
	}

	// Output:
	// 2
	// 3
	// 4
	// 5
}

func ExampleDeque_PopAll() {
	q := deque.From([]int{1, 2, 3, 4, 5})
	for x := range q.PopAll() {
		fmt.Println(x)
	}
	fmt.Println(q)

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// []
}
