# Pathfinding Test Maps

This directory holds ASCII grid maps used to test pathfinding algorithms
(BFS, Dijkstra, A*, etc.). Each map is a self-contained text file.

## Tile legend

| Char | Meaning              |
|------|----------------------|
| `.`  | walkable space       |
| `#`  | non-walkable wall    |
| `$`  | start point          |
| `@`  | goal / end point     |

## Format rules

- **Rectangular**: every row has the same number of columns. No ragged rows.
- **Exactly one `$`** and **exactly one `@`** per map.
- Only the four characters above appear in the grid. No padding spaces.
- Lines end with `\n`; the file ends with a single trailing newline.
- Filenames are zero-based integers with a `.map` extension: `0.map`, `1.map`, ...
- Any reasonable size is allowed. Keep width/height between ~6 and ~40 so maps
  stay human-readable and quick to load.

## Movement assumptions

Maps are authored for **4-directional movement** (up/down/left/right). `$` and
`@` are themselves walkable. A solver treats `.`, `$`, and `@` as traversable
and `#` as blocked. Maps that are only solvable with diagonal movement should be
avoided unless explicitly noted, so the same fixtures work across solvers.

## Solvability

**Every map must have at least one valid path from `$` to `@`.** These are
positive fixtures — they exist to confirm a solver *finds* a path and, ideally,
the shortest one. Validate new maps with a BFS reachability check before adding
them.

To intentionally test failure handling, create a separate clearly-named file
(e.g. `unsolvable-*.map`) so it is never mistaken for a positive fixture.

## Design goals for a good test suite

Vary the maps so they exercise different code paths:

- **Open fields** — many equal-cost shortest paths (tie-breaking).
- **Single corridors** — exactly one route (correctness).
- **Mazes** — long winding solutions, lots of dead ends (frontier size).
- **Rooms + doors** — bottlenecks through narrow gaps.
- **Spirals / concentric walls** — long paths in a small footprint.
- **Start/goal placement** — corners, centers, adjacent, far apart.

## Adding a map

1. Author the grid following the format rules above.
2. Confirm exactly one `$` and one `@`.
3. Run a BFS check to confirm `$` can reach `@`.
4. Save as the next integer filename, e.g. `11.map`.

A quick validator (Python, 4-directional):

```python
import sys
from collections import deque

grid = [line.rstrip("\n") for line in open(sys.argv[1])]
start = goal = None
for r, row in enumerate(grid):
    for c, ch in enumerate(row):
        if ch == "$": start = (r, c)
        if ch == "@": goal = (r, c)

q, seen = deque([start]), {start}
while q:
    r, c = q.popleft()
    if (r, c) == goal:
        print("solvable"); break
    for dr, dc in ((1,0),(-1,0),(0,1),(0,-1)):
        nr, nc = r+dr, c+dc
        if 0 <= nr < len(grid) and 0 <= nc < len(grid[nr]):
            if grid[nr][nc] != "#" and (nr, nc) not in seen:
                seen.add((nr, nc)); q.append((nr, nc))
else:
    print("UNSOLVABLE")
```
