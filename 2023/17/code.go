package main

import (
	"aoc-in-go/ez"
	"strings"

	pq "github.com/emirpasic/gods/queues/priorityqueue"
	"github.com/jpillora/puzzler/harness/aoc"
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

// Cell contains information about the cell as a Path passes over it
// It is primarily used as a cache key for optimization
type Cell struct {
	Direction string
	Row       int
	Col       int
	DirCount  int
}

// Path tracks the path's accumulated heat, and points (for debugging)
// Path is composed of Cell
type Path struct {
	Cell
	Heat   int
	Points []Point
}

// Point is just a tracking Row Column collection for debugging
type Point struct {
	R int
	C int
}

func run(part2 bool, input string) any {
	// Use a PriorityQueue to compare heat values, using a - b ensures we're only ever testing the path with the least heat
	// ref: https://github.com/emirpasic/gods#priorityqueue
	q := pq.NewWith(func(a, b interface{}) int {
		return a.(Path).Heat - b.(Path).Heat
	})

	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	grid := make([][]int, len(lines))
	for i, line := range lines {
		grid[i] = make([]int, len(line))
		for j, char := range strings.Split(line, "") {
			grid[i][j] = ez.Atoi(char)
		}
	}

	targetRow := len(grid) - 1
	targetCol := len(grid[0]) - 1
	cache := make(map[Cell]int)
	// Part 1
	maxStraight := 3
	minStraight := 0
	if part2 {
		// Part 2
		maxStraight = 10
		minStraight = 3
	}

	// Kick off going East and South on the grid
	q.Enqueue(Path{
		Cell: Cell{
			Direction: "E",
			Row:       0,
			Col:       1,
			DirCount:  1,
		},
		Heat:   0,
		Points: []Point{{R: 0, C: 0}},
	})
	q.Enqueue(Path{
		Cell: Cell{
			Direction: "S",
			Row:       1,
			Col:       0,
			DirCount:  1,
		},
		Heat:   0,
		Points: []Point{{R: 0, C: 0}},
	})

	for {
		pathI, pathExists := q.Dequeue()
		if !pathExists {
			break
		}
		path := pathI.(Path)

		// Make sure path is within the grid
		if path.Row < 0 || path.Row > targetRow || path.Col < 0 || path.Col > targetCol {
			continue
		}

		// Calculate heat
		heat := path.Heat + grid[path.Row][path.Col]

		// Exit condition
		if path.Row == targetRow && path.Col == targetCol {
			if path.DirCount < minStraight {
				continue
			}

			return heat
		}

		// Check cache
		if cacheHeat, exists := cache[path.Cell]; exists {
			if cacheHeat <= heat {
				continue
			}
		}
		cache[path.Cell] = heat

		// Turns
		if path.DirCount > minStraight {
			newDirs := []string{"E", "W"}
			if path.Direction == "W" || path.Direction == "E" {
				newDirs = []string{"N", "S"}
			}

			for _, newDir := range newDirs {
				newDir := newDir
				nextRow, nextCol := NextPoint(path.Row, path.Col, newDir)
				q.Enqueue(Path{
					Cell: Cell{
						Direction: newDir,
						Row:       nextRow,
						Col:       nextCol,
						DirCount:  1,
					},
					Heat:   heat,
					Points: append(path.Points, Point{R: path.Row, C: path.Col}),
				})
			}
		}

		// Go straight
		if path.DirCount < maxStraight {
			nextRow, nextCol := NextPoint(path.Row, path.Col, path.Direction)
			q.Enqueue(Path{
				Cell: Cell{
					Direction: path.Direction,
					Row:       nextRow,
					Col:       nextCol,
					DirCount:  path.DirCount + 1,
				},
				Heat:   heat,
				Points: append(path.Points, Point{R: path.Row, C: path.Col}),
			})
		}
	}

	ez.Log("Failed to solve the problem!")
	return 0
}

// NextPoint returns the row and column value given a direction (NSEW) to go in
func NextPoint(row, col int, newDir string) (int, int) {
	rowMod := 0
	colMod := 0
	switch newDir {
	case "N":
		rowMod = -1
	case "S":
		rowMod = 1
	case "E":
		colMod = 1
	case "W":
		colMod = -1
	}

	return row + rowMod, col + colMod
}
