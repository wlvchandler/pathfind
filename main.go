package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

var algorithms = map[string]solver{
	"bfs": bfs,
}

func algoNames() string {
	names := make([]string, 0, len(algorithms))
	for name := range algorithms {
		names = append(names, name)
	}
	slices.Sort(names)
	return strings.Join(names, ", ")
}

func main() {
	algo := flag.String("algo", "bfs", "search algorithm: "+algoNames())
	delay := flag.Duration("delay", 40*time.Millisecond, "pause between animation frames")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: %s [-algo name] [-delay d] <map-file>\n", os.Args[0])
		os.Exit(2)
	}

	solve, ok := algorithms[*algo]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown algo %q (have: %s)\n", *algo, algoNames())
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
