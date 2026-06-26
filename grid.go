package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	startTile = '$'
	goalTile  = '@'
	wallTile  = '#'
)

type Coord struct {
	X, Y int
}

func readGrid(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimRight(string(data), "\n"), "\n"), nil
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
