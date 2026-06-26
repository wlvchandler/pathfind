package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"
)

const (
	pathTile     = '*'
	exploredTile = 'o'
)

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
