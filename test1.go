package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Queue[T any] struct {
	elements []T
}

func (q *Queue[T]) Size() int {
	return len(q.elements)
}

func (q *Queue[T]) Enqueue(value T) {
	q.elements = append(q.elements, value)
}

func (q *Queue[T]) Dequeue() (T, error) {
	if len(q.elements) == 0 {
		var zero T
		return zero, errors.New("queue is empty")
	}

	value := q.elements[0]

	var zero T
	q.elements[0] = zero
	q.elements = q.elements[1:]

	return value, nil
}

type Coord struct {
	X int
	Y int
}

func find_start(grid []string) (Coord, error) {
	for i := range grid {
		for j := 0; j < len(grid[i]); j++ {
			if grid[i][j] == '$' {
				return Coord{X: j, Y: i}, nil
			}
		}
	}
	return Coord{X: -1, Y: -1}, errors.New("No starting point found")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func visited(m *map[Coord]Coord, c *Coord) bool {
	_, ok := (*m)[Coord{c.X, c.Y}]
	return ok
}

func valid_move(grid []string, c *Coord) bool {
	return c.X >= 0 && c.X < len(grid[0]) && c.Y >= 0 && c.Y < len(grid) && grid[c.Y][c.X] != '#'
}

func print_grid(grid []string) {
	for i := 0; i < len(grid); i++ {
		fmt.Println(grid[i])
	}
}

func main() {
	fpath := filepath.Join("maps", "10.map")
	dat, err := os.ReadFile(fpath)
	check(err)

	lines := strings.Split(string(dat), "\n")
	start_idx, err := find_start(lines)
	check(err)

	parent := make(map[Coord]Coord)
	var q Queue[Coord]

	q.Enqueue(start_idx)
	parent[start_idx] = Coord{X: -1, Y: -1}
	looking := true
	var current, end Coord

	print_grid(lines)

	for looking {
		current, err := q.Dequeue()
		check(err)

		if current != start_idx {
			bytes := []byte(lines[current.Y])
			bytes[current.X] = 'x'
			lines[current.Y] = string(bytes)
		}

		neighbors := [4]Coord{
			{current.X + 1, current.Y},
			{current.X - 1, current.Y},
			{current.X, current.Y + 1},
			{current.X, current.Y - 1},
		}

		for i := 0; i < 4; i++ {
			if valid_move(lines, &neighbors[i]) && !visited(&parent, &neighbors[i]) {
				q.Enqueue(neighbors[i])
				parent[neighbors[i]] = current
				if lines[neighbors[i].Y][neighbors[i].X] == '@' {
					looking = false
					end = neighbors[i]
				}
			}
		}
	}

	current = end
	path := []Coord{current}

	for current != start_idx {
		path = append(path, parent[current])
		current = parent[current]
		if current != start_idx {
			bytes := []byte(lines[current.Y])
			bytes[current.X] = '^'
			lines[current.Y] = string(bytes)
		}
	}

	for i := 0; i < len(path); i++ {
		fmt.Printf("(%d, %d)\n", path[i].X, path[i].Y)
	}

	print_grid(lines)

}
