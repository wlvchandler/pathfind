package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"strings"
	"time"
)

const (
	startTile    = '$'
	goalTile     = '@'
	wallTile     = '#'
	pathTile     = '*'
	exploredTile = 'o'
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
	delay := flag.Duration("delay", 40*time.Millisecond, "pause between animation frames")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: %s [-delay d] <map-file>\n", os.Args[0])
		os.Exit(2)
	}

	grid, err := readGrid(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	path, explored, err := solve(grid)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	animate(grid, explored, path, *delay)
}

func readGrid(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimRight(string(data), "\n"), "\n"), nil
}

// BFS, so the path it finds is a shortest one. explored is the order cells were
// visited, which the animation replays as the expanding search frontier.
func solve(grid []string) (path, explored []Coord, err error) {
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

// animate redraws the grid in place once per step: first the search frontier
// flooding outward, then the shortest path being traced back.
func animate(grid []string, explored, path []Coord, delay time.Duration) {
	frame := make([][]byte, len(grid))
	for i, row := range grid {
		frame[i] = []byte(row)
	}

	hideCursor()
	restoreOnInterrupt()
	defer showCursor()

	clearScreen()
	draw(frame)

	for _, c := range explored {
		mark(frame, c, exploredTile)
		draw(frame)
		time.Sleep(delay)
	}
	for _, c := range path {
		mark(frame, c, pathTile)
		draw(frame)
		time.Sleep(delay)
	}

	fmt.Printf("\npath length: %d steps\n", len(path)-1)
}

func mark(frame [][]byte, c Coord, tile byte) {
	if cur := frame[c.Y][c.X]; cur == startTile || cur == goalTile {
		return
	}
	frame[c.Y][c.X] = tile
}

func draw(frame [][]byte) {
	var b strings.Builder
	b.WriteString("\033[H")
	for _, row := range frame {
		for _, tile := range row {
			b.WriteString(colorize(tile))
		}
		b.WriteByte('\n')
	}
	fmt.Print(b.String())
}

func colorize(tile byte) string {
	switch tile {
	case wallTile:
		return "\033[90m#\033[0m"
	case startTile:
		return "\033[1;32m$\033[0m"
	case goalTile:
		return "\033[1;31m@\033[0m"
	case exploredTile:
		return "\033[34mo\033[0m"
	case pathTile:
		return "\033[1;33m*\033[0m"
	default:
		return string(tile)
	}
}

func clearScreen() { fmt.Print("\033[2J\033[H") }
func hideCursor()  { fmt.Print("\033[?25l") }
func showCursor()  { fmt.Print("\033[?25h") }

func restoreOnInterrupt() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		showCursor()
		os.Exit(130)
	}()
}
