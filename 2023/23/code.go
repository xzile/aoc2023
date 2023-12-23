package main

import (
	"container/list"
	"strings"

	dij "github.com/RyanCarrier/dijkstra"
	"github.com/jpillora/puzzler/harness/aoc"
	"github.com/samber/lo"
)

func main() {
	aoc.Harness(run)
}

// on code change, run will be executed 4 times:
// 1. with: false (part1), and example input
// 2. with: true (part2), and example input
// 3. with: false (part1), and user input
// 4. with: true (part2), and user input
// the return value of each run is printed to stdout

type Path struct {
	Cell Cell
	From Cell
	Dist int
}

type Cell struct {
	ID   int
	R    int
	C    int
	Type string
}

type Grid [][]Cell

func run(part2 bool, input string) any {
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")

	var startCell *Cell
	var endCell *Cell
	g := dij.NewGraph()
	// Map the grid to Cells
	grid := make(Grid, len(lines))
	id := 1
	for i, line := range lines {
		grid[i] = make([]Cell, len(line))
		for j, char := range strings.Split(line, "") {
			grid[i][j] = Cell{
				ID:   id,
				R:    i,
				C:    j,
				Type: char,
			}

			// Start is the only . in the first row
			if char == "." && i == 0 {
				startCell = &grid[i][j]
			}
			// End is the only . in the last row
			if char == "." && i == len(lines)-1 {
				endCell = &grid[i][j]
			}

			// Add an ID'd vertex to the graph for part 1
			g.AddVertex(id)
			id++
		}
	}

	// Part 2
	if part2 {
		// Collect all junctions
		junctions := make([]Cell, 0)
		junctions = append(junctions, *startCell)
		junctions = append(junctions, *endCell)

		for i := 0; i < len(grid); i++ {
			for j := 0; j < len(grid[i]); j++ {
				if grid[i][j].IsJunction(grid) {
					junctions = append(junctions, grid[i][j])
				}
			}
		}

		// Get distance from junction to directly connected junctions
		dg := make(map[Cell]map[Cell]int)
		for n := 0; n < len(junctions); n++ {
			dg[junctions[n]] = make(map[Cell]int)

			q := list.New()
			q.PushBack(Path{
				Cell: junctions[n],
				From: junctions[n],
				Dist: 0,
			})
			for q.Len() > 0 {
				pAny := q.Front()
				p := pAny.Value.(Path)
				c := p.Cell
				i := c.R
				j := c.C

				// Exit condition, we've hit a junction, add to map
				if c.IsJunction(grid) && c != p.From {
					dg[junctions[n]][c] = p.Dist
					q.Remove(pAny)
					continue
				}
				// Exit condition, we've hit the desired end cell, add to map
				if c == *endCell {
					dg[junctions[n]][c] = p.Dist
					q.Remove(pAny)
					continue
				}

				// Variables that will be used later, xV = is that cell the one you came from?
				var n, e, s, w Cell
				var nV, eV, sV, wV bool
				if i-1 >= 0 {
					n = grid[i-1][j]
					nV = (n.ID == p.From.ID)
				}
				if i+1 < len(grid) {
					s = grid[i+1][j]
					sV = (s.ID == p.From.ID)
				}
				if j-1 >= 0 {
					w = grid[i][j-1]
					wV = (w.ID == p.From.ID)
				}
				if j+1 < len(grid[i]) {
					e = grid[i][j+1]
					eV = (e.ID == p.From.ID)
				}

				// Part 2 doesn't care about slopes, just make sure it's on the grid and it's not a forest (#)
				newCells := make([]Cell, 0)
				if !nV && n.ID != 0 && n.Type != "#" {
					newCells = append(newCells, n)
				}
				if !eV && e.ID != 0 && e.Type != "#" {
					newCells = append(newCells, e)
				}
				if !sV && s.ID != 0 && s.Type != "#" {
					newCells = append(newCells, s)
				}
				if !wV && w.ID != 0 && w.Type != "#" {
					newCells = append(newCells, w)
				}

				for i := range newCells {
					q.PushBack(Path{
						Cell: newCells[i],
						From: c,
						Dist: p.Dist + 1,
					})
				}

				q.Remove(pAny)
			}
		}

		return LongestRoute(*startCell, *endCell, dg)
	}

	// Part 1
	q := list.New()
	visited := make([]Cell, 0)
	q.PushBack(Path{
		Cell: *startCell,
		From: *startCell,
		Dist: 0,
	})
	for q.Len() > 0 {
		pAny := q.Front()
		p := pAny.Value.(Path)
		c := p.Cell
		visited = append(visited, c)
		i := c.R
		j := c.C

		// Variables that will be used later, xV = is that cell the one you came from?
		var n, e, s, w Cell
		var nV, eV, sV, wV bool
		if i-1 >= 0 {
			n = grid[i-1][j]
			nV = (n.ID == p.From.ID)
		}
		if i+1 < len(grid) {
			s = grid[i+1][j]
			sV = (s.ID == p.From.ID)
		}
		if j-1 >= 0 {
			w = grid[i][j-1]
			wV = (w.ID == p.From.ID)
		}
		if j+1 < len(grid[i]) {
			e = grid[i][j+1]
			eV = (e.ID == p.From.ID)
		}

		newCells := make([]Cell, 0)
		switch c.Type {
		case ".": // Any non forest direction that's not a slope pointing to it
			if !nV && n.ID != 0 && n.Type != "#" && n.Type != "v" {
				newCells = append(newCells, n)
			}
			if !eV && e.ID != 0 && e.Type != "#" && e.Type != "<" {
				newCells = append(newCells, e)
			}
			if !sV && s.ID != 0 && s.Type != "#" && s.Type != "^" {
				newCells = append(newCells, s)
			}
			if !wV && w.ID != 0 && w.Type != "#" && w.Type != ">" {
				newCells = append(newCells, w)
			}
		case ">": // East only
			if !eV && e.ID != 0 && e.Type != "#" && e.Type != "<" {
				newCells = append(newCells, e)
			}
		case "^": // North only
			if !nV && n.ID != 0 && n.Type != "#" && n.Type != "v" {
				newCells = append(newCells, n)
			}
		case "v": // South only
			if !sV && s.ID != 0 && s.Type != "#" && s.Type != "^" {
				newCells = append(newCells, s)
			}
		case "<": // West only
			if !wV && w.ID != 0 && w.Type != "#" && w.Type != ">" {
				newCells = append(newCells, w)
			}
		}

		for i := range newCells {
			// Add path for any new cells to the graph
			g.AddArc(c.ID, newCells[i].ID, 1)

			// Continue on if we've not visited the newCell
			if !lo.Contains(visited, newCells[i]) {
				q.PushBack(Path{
					Cell: newCells[i],
					From: c,
					Dist: p.Dist + 1,
				})
			}
		}

		q.Remove(pAny)
	}

	// Get the longest
	// ref: https://github.com/RyanCarrier/dijkstra
	bests := lo.Must(g.Longest(startCell.ID, endCell.ID))
	return bests.Distance
}

type Route struct {
	Start   Cell
	Visited []Cell
	Dist    int
}

// LongestRoute returns the maximum distance to without repeating a cell, from start to end
// dg should be a map keyed by FromCell then ToCell, and the distance stored as the value
func LongestRoute(start, end Cell, dg map[Cell]map[Cell]int) int {
	dist := 0
	visited := make([]Cell, 0)
	visited = append(visited, start)

	q := list.New()
	q.PushBack(Route{
		Start:   start,
		Dist:    0,
		Visited: visited,
	})

	for q.Len() > 0 {
		rAny := q.Front()
		r := rAny.Value.(Route)

		for c, d := range dg[r.Start] {
			if c == end {
				// Exit condition, we've hit the desired end cell, set the max
				dist = max(dist, r.Dist+d)
				continue
			}

			if lo.Contains(r.Visited, c) {
				// Exit condition, we've already visited the cell in our traversal
				continue
			}

			// Copy the visited slice for each new route step we take
			newVisited := make([]Cell, 0)
			newVisited = append(newVisited, r.Visited...)
			newVisited = append(newVisited, c)
			q.PushBack(Route{
				Start:   c,
				Dist:    r.Dist + d,
				Visited: newVisited,
			})
		}

		q.Remove(rAny)
	}

	return dist
}

// IsJunction evaluates a cell on the grid to determine if it has more than 2 exits
func (c Cell) IsJunction(g Grid) bool {
	paths := 0
	i := c.R
	j := c.C
	if i-1 >= 0 && g[i-1][j].Type != "#" {
		paths++
	}
	if i+1 < len(g) && g[i+1][j].Type != "#" {
		paths++
	}
	if j-1 >= 0 && g[i][j-1].Type != "#" {
		paths++
	}
	if j+1 < len(g[i]) && g[i][j+1].Type != "#" {
		paths++
	}
	return paths > 2
}
