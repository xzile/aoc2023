package main

import (
	"aoc-in-go/ez"
	"fmt"
	"strings"

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

type Cell struct {
	R    int
	C    int
	Type string
}

type Step struct {
	Cell Cell
}

type Grid [][]Cell

type Steps []Step

func run(part2 bool, input string) any {
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")

	// Map the grid to Cells
	grid := make(Grid, len(lines))
	var startCell Cell
	for i, line := range lines {
		grid[i] = make([]Cell, len(line))
		for j, char := range strings.Split(line, "") {
			grid[i][j] = Cell{
				R: i,
				C: j,
				Type: func() string {
					if char == "#" {
						return char
					}
					return "."
				}(),
			}

			if char == "S" {
				startCell = grid[i][j]
			}
		}
	}

	// Part 2
	if part2 {
		if len(lines) < 20 {
			// Example, skip
			return 1
		}

		// I'll be honest, the math here is a bit out of my league and I kinda get it but not to the point I know what's going
		// I converted other solutions that were posted to Go and it worked
		// Ultimately, we're using 3 points along a quadratic curve of equal distance to determine the quadratic coefficients

		// Grid is square with all edges and the starting column/row are not blocked by a #
		sqLen := len(lines)
		allStepsToTake := 26501365
		// Our "x", for the quadratic solution at the end
		x := allStepsToTake / sqLen
		// Capture points for our a, b, and c "step counts"/points
		cCapture := allStepsToTake % sqLen
		bCapture := cCapture + sqLen
		aCapture := cCapture + 2*sqLen

		// We'll expand the grid to support at least 2*sqLen + cCapture
		grid = ExpandGrid(grid, 4)
		// Adjust our startCell so we don't hit the edges
		startCell.C = startCell.C + (2 * sqLen)
		startCell.R = startCell.R + (2 * sqLen)

		steps := make(Steps, 0)
		steps = append(steps, Step{Cell: startCell})
		var aStepCount, bStepCount, cStepCount int
		for i := 0; i < aCapture; i++ {
			steps = steps.TakeSteps(grid)
			// Once we hit one of the capture points, store those values
			if i+1 == aCapture {
				aStepCount = len(steps)
			}
			if i+1 == bCapture {
				bStepCount = len(steps)
			}
			if i+1 == cCapture {
				cStepCount = len(steps)
			}
		}

		// Math the coefficients, I'm a web developer, not a mathematician, not sure how these actually work
		c := cStepCount
		a := (aStepCount + cStepCount - 2*bStepCount) / 2
		b := bStepCount - cStepCount - a

		// The quadratic equation to produce our result
		return a*x*x + b*x + c
	}

	// Part 1
	stepsToTake := 64
	if len(lines) < 20 {
		// Example
		stepsToTake = 6
	}

	steps := make(Steps, 0)
	steps = append(steps, Step{Cell: startCell})
	for i := 0; i < stepsToTake; i++ {
		steps = steps.TakeSteps(grid)
	}

	return len(steps)
}

// ExpandGrid takes a Grid and replicates it given a factor count, inclusive of the original grid
func ExpandGrid(grid Grid, factor int) Grid {
	initLen := len(grid)
	// Expand the columns
	for i := 0; i < initLen; i++ {
		line := grid[i]
		for j := 0; j < factor; j++ {
			grid[i] = append(grid[i], line...)
		}
	}
	// Expand the rows
	for j := 0; j < factor; j++ {
		for i := 0; i < initLen; i++ {
			grid = append(grid, grid[i])
		}
	}
	// Return a new Grid with all new cells
	rowLen := len(grid)
	colLen := len(grid[0])
	newGrid := make(Grid, rowLen)
	for i := 0; i < len(grid); i++ {
		newGrid[i] = make([]Cell, colLen)
		for j := 0; j < len(grid[i]); j++ {
			newGrid[i][j] = Cell{
				R:    i,
				C:    j,
				Type: grid[i][j].Type,
			}
		}
	}

	return newGrid
}

// TakeSteps iterates over the current steps and produces a new set of "steps", evaluating valid moves N/E/S/W
func (s Steps) TakeSteps(grid Grid) []Step {
	newSteps := make([]Step, 0)

	for _, step := range s {
		i := step.Cell.R
		j := step.Cell.C
		if grid[i][j].C != j || grid[i][j].R != i {
			ez.Log("bad grid", i, j, grid[i][j].R, grid[i][j].C)
		}

		// North
		if i-1 >= 0 {
			nextCell := grid[i-1][j]
			if nextCell.Type != "#" {
				newSteps = append(newSteps, Step{Cell: nextCell})
			}
		}
		// South
		if i+1 < len(grid) {
			nextCell := grid[i+1][j]
			if nextCell.Type != "#" {
				newSteps = append(newSteps, Step{Cell: nextCell})
			}
		}
		// East
		if j+1 < len(grid[0]) {
			nextCell := grid[i][j+1]
			if nextCell.Type != "#" {
				newSteps = append(newSteps, Step{Cell: nextCell})
			}
		}
		// West
		if j-1 >= 0 {
			nextCell := grid[i][j-1]
			if nextCell.Type != "#" {
				newSteps = append(newSteps, Step{Cell: nextCell})
			}
		}
	}

	// Remove duplicates, this could probably be more optimized to be a cache
	newSteps = lo.UniqBy(newSteps, func(item Step) string {
		return fmt.Sprintf("%d,%d", item.Cell.C, item.Cell.R)
	})

	return newSteps
}
