package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

const (
	startTile = '$'
	goalTile  = '@'
	wallTile  = '#'
	pathTile  = '*'
)

type Coord struct {
	X, Y int
}

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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <map-file>\n", os.Args[0])
		os.Exit(2)
	}

	grid, err := readGrid(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	path, err := solve(grid)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	render(grid, path)
}

func readGrid(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimRight(string(data), "\n"), "\n"), nil
}

// BFS, so the path it finds is a shortest one.
func solve(grid []string) ([]Coord, error) {
	start, err := find(grid, startTile)
	if err != nil {
		return nil, err
	}

	parent := map[Coord]Coord{start: {-1, -1}}
	var queue Queue[Coord]
	queue.Enqueue(start)

	for queue.Len() > 0 {
		current, _ := queue.Dequeue()
		if grid[current.Y][current.X] == goalTile {
			return reconstruct(parent, start, current), nil
		}
		for _, next := range neighbors(current) {
			if walkable(grid, next) && !seen(parent, next) {
				parent[next] = current
				queue.Enqueue(next)
			}
		}
	}
	return nil, errors.New("no path from start to goal")
}

func find(grid []string, tile byte) (Coord, error) {
	for y, row := range grid {
		for x := range len(row) {
			if row[x] == tile {
				return Coord{X: x, Y: y}, nil
			}
		}
	}
	return Coord{}, fmt.Errorf("tile %q not found in map", tile)
}

func neighbors(c Coord) [4]Coord {
	return [4]Coord{
		{c.X + 1, c.Y},
		{c.X - 1, c.Y},
		{c.X, c.Y + 1},
		{c.X, c.Y - 1},
	}
}

func walkable(grid []string, c Coord) bool {
	return c.Y >= 0 && c.Y < len(grid) &&
		c.X >= 0 && c.X < len(grid[c.Y]) &&
		grid[c.Y][c.X] != wallTile
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

func render(grid []string, path []Coord) {
	overlay := make([][]byte, len(grid))
	for i, row := range grid {
		overlay[i] = []byte(row)
	}
	for _, c := range path {
		if tile := overlay[c.Y][c.X]; tile != startTile && tile != goalTile {
			overlay[c.Y][c.X] = pathTile
		}
	}
	for _, row := range overlay {
		fmt.Println(string(row))
	}
	fmt.Printf("\npath length: %d steps\n", len(path)-1)
}
