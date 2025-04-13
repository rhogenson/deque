package heap_test

import (
	"cmp"
	"fmt"

	"github.com/rhogenson/container/heap"
)

func Example_dijkstra() {
	const (
		maze = `
---------------------
    |   |           |
| --- | | --- ----- |
|   | |   |   |     |
|-- |-----| ----- --|
|   |     |     | | |
| --- --- | --- | | |
| |   | |     | | | |
| | --- | | --- | | |
| | |   | | |   |   |
| | --- | |-- --|-- |
| |   | | |   | |   |
| |-- | | | --- | --|
| |   | |   |     | |
| | --| ----- ----- |
| |   |       |     |
| | | | --------- | |
| | | |     |   | | |
| --- |---- | | | | |
|     |       |      
---------------------
`
		width  = 21
		height = 21
	)

	type point struct{ x, y int }
	start := point{0, 1}
	goal := point{20, 19}
	neighbors := func(p point) []point {
		return []point{
			{p.x - 1, p.y},
			{p.x + 1, p.y},
			{p.x, p.y - 1},
			{p.x, p.y + 1},
		}
	}
	walkable := func(p point) bool {
		return 0 <= p.y && p.y < height &&
			0 <= p.x && p.x < width &&
			maze[p.y*(width+1)+p.x+1] == ' '
	}

	visited := map[point]int{start: 1}
	q := heap.New(func(x, y point) int { return cmp.Compare(visited[x], visited[y]) })
	q.Push(start)
Dijkstra:
	for {
		p, ok := q.Pop()
		if !ok {
			fmt.Println("Giving up!")
			return
		}
		for _, neighbor := range neighbors(p) {
			if !walkable(neighbor) || visited[neighbor] > 0 {
				continue
			}
			visited[neighbor] = visited[p] + 1
			if neighbor == goal {
				break Dijkstra
			}
			q.Push(neighbor)
		}
	}

	completedMaze := []byte(maze)
	fillIn := func(p point) {
		completedMaze[p.y*(width+1)+p.x+1] = '*'
	}
	for p := goal; p != start; {
		fillIn(p)
		closestPoint := p
		for _, neighbor := range neighbors(p) {
			if walkable(neighbor) && visited[neighbor] > 0 && visited[neighbor] < visited[closestPoint] {
				closestPoint = neighbor
			}
		}
		p = closestPoint
	}
	fillIn(start)

	fmt.Printf("%s\n", completedMaze)

	// Output:
	// ---------------------
	// **  |   |    *******|
	// |*--- | | ---*-----*|
	// |***| |   |***|  ***|
	// |--*|-----|*-----*--|
	// |***|*****|*    |*| |
	// |*---*---*|*--- |*| |
	// |*|***| |***  | |*| |
	// |*|*--- | | --- |*| |
	// |*|*|   | | |   |***|
	// |*|*--- | |-- --|--*|
	// |*|***| | |   | |***|
	// |*|--*| | | --- |*--|
	// |*|***| |   |*****| |
	// |*|*--| -----*----- |
	// |*|***|*******|     |
	// |*| |*|*--------- | |
	// |*| |*|*****|***| | |
	// |*---*|----*|*|*| | |
	// |*****|    ***|******
	// ---------------------
}

func ExampleNew() {
	priority := map[string]int{
		"job1": 10,
		"job2": 30,
		"job3": 100,
		"job4": 20,
	}
	h := heap.New(func(j1, j2 string) int { return cmp.Compare(priority[j1], priority[j2]) })
	h.Push("job1")
	h.Push("job2")
	h.Push("job3")
	h.Push("job4")
	if highestPriorityJob, ok := h.Pop(); ok {
		fmt.Println(highestPriorityJob)
	}

	// Output:
	// job1
}

func ExampleHeap_Len() {
	h := heap.New(cmp.Compare[int])
	h.Push(1)
	h.Push(2)
	h.Push(3)
	fmt.Println(h.Len())

	// Output:
	// 3
}

func ExampleHeap_Grow() {
	h := heap.New(cmp.Compare[int])
	h.Grow(3)
	// Push without allocating:
	h.Push(1)
	h.Push(2)
	h.Push(3)
}

func ExampleHeap_Push() {
	h := heap.New(cmp.Compare[int])
	h.Push(1)
	h.Push(2)
	h.Push(3)
}

func ExampleHeap_Pop() {
	h := heap.New(cmp.Compare[int])
	h.Push(1)
	h.Push(2)
	h.Push(3)
	if n, ok := h.Pop(); ok {
		fmt.Println(n)
	}

	// Output:
	// 1
}
