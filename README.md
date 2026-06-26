# pathfind

A small grid pathfinder in Go. It reads an ASCII map, runs a breadth-first
search from the start tile to the goal, and animates the search in the terminal:
the frontier floods outward (`o`), then the shortest path is traced back (`*`).

BFS on an unweighted grid is guaranteed to find a shortest path, which is why
it's the starting point here.

## Map format

Maps are rectangular ASCII grids. See [`maps/CLAUDE.md`](maps/CLAUDE.md) for the
full spec.

| Char | Meaning        |
|------|----------------|
| `.`  | walkable space |
| `#`  | wall           |
| `$`  | start          |
| `@`  | goal           |

Each map has exactly one `$` and one `@`, and movement is 4-directional.

## Usage

```sh
go run . maps/10.map              # animate BFS on a map
go run . -algo bfs maps/10.map    # pick the algorithm
go run . -delay 0 maps/10.map     # no frame pause
```

`-algo` selects the search algorithm (currently `bfs`); more can be registered
in `algorithms` without touching the CLI or the animation.

The output is the original map with the path marked using `*`, followed by the
path length:

```
###################
#$****#.....#.....#
#####*#.#.###.#.###
...
path length: 48 steps
```

## Next steps

- Register A\* and Dijkstra as additional `-algo` choices to compare against BFS.
- Add `unsolvable-*.map` fixtures to exercise the no-path branch.
