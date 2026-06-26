# pathfind

A small grid pathfinder in Go. It reads an ASCII map, runs a breadth-first
search from the start tile to the goal, and prints the map with the shortest
path drawn on it.

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
go run . maps/10.map
```

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

- Add A\* and Dijkstra to compare against BFS on the same fixtures.
- Add `unsolvable-*.map` fixtures to exercise the no-path branch.
