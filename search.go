package main

import (
	"errors"
	"slices"
)

// A solver searches grid from start to goal. It returns the path and explored,
// the order cells were visited, which the animation replays as the search
// frontier.
type solver func(grid []string) (path, explored []Coord, err error)

type Queue[T any] struct {
	elements []T
}

func (q *Queue[T]) Len() int {
	return len(q.elements)
}

func (q *Queue[T]) Enqueue(value T) {
	q.elements = append(q.elements, value)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	var value T
	if len(q.elements) == 0 {
		return value, false
	}
	value = q.elements[0]
	q.elements = q.elements[1:]
	return value, true
}

// bfs explores in expanding rings, so the first path it reaches is a shortest one.
func bfs(grid []string) (path, explored []Coord, err error) {
	start, err := find(grid, startTile)
	if err != nil {
		return nil, nil, err
	}

	parent := map[Coord]Coord{start: {-1, -1}}
	var queue Queue[Coord]
	queue.Enqueue(start)

	for queue.Len() > 0 {
		current, _ := queue.Dequeue()
		explored = append(explored, current)
		if grid[current.Y][current.X] == goalTile {
			return reconstruct(parent, start, current), explored, nil
		}
		for _, next := range neighbors(current) {
			if walkable(grid, next) && !seen(parent, next) {
				parent[next] = current
				queue.Enqueue(next)
			}
		}
	}
	return nil, explored, errors.New("no path from start to goal")
}

func seen(parent map[Coord]Coord, c Coord) bool {
	_, ok := parent[c]
	return ok
}

func reconstruct(parent map[Coord]Coord, start, goal Coord) []Coord {
	var path []Coord
	for c := goal; c != start; c = parent[c] {
		path = append(path, c)
	}
	path = append(path, start)
	slices.Reverse(path)
	return path
}
