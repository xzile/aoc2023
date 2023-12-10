package main

import (
	"aoc-in-go/ez"
	"math"
	"strings"

	"github.com/jpillora/puzzler/harness/aoc"
	"github.com/samber/lo"
)

func main() {
	aoc.Harness(run)
}

type Pipe struct {
	Dir []string
}

// on code change, run will be executed 4 times:
// 1. with: false (part1), and example input
// 2. with: true (part2), and example input
// 3. with: false (part1), and user input
// 4. with: true (part2), and user input
// the return value of each run is printed to stdout
func run(part2 bool, input string) any {
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	// Pipe mapping, denoting the two directions a given pip connects
	pipes := map[string]Pipe{
		"|": {Dir: []string{"N", "S"}},
		"-": {Dir: []string{"E", "W"}},
		"L": {Dir: []string{"N", "E"}},
		"J": {Dir: []string{"N", "W"}},
		"7": {Dir: []string{"S", "W"}},
		"F": {Dir: []string{"S", "E"}},
	}
	// Locs will represent a stored point location on the path
	var locs []ez.Point

	// Find start
	var startRow, startCol int
	for i, line := range lines {
		if startCol = strings.Index(line, "S"); startCol != -1 {
			startRow = i
			break
		}
	}
	locs = append(locs, ez.Point{X: float64(startCol), Y: float64(startRow)})

	// Find a connected piece
	move := ""
	switch {
	// Look East
	case lo.Contains([]string{"7", "J", "-"}, string(lines[startRow][startCol+1])):
		move = "E"
	// Look West
	case lo.Contains([]string{"F", "L", "-"}, string(lines[startRow][startCol-1])):
		move = "W"
	// Look North
	case lo.Contains([]string{"7", "F", "|"}, string(lines[startRow-1][startCol])):
		move = "N"
	// Look South
	case lo.Contains([]string{"L", "J", "|"}, string(lines[startRow+1][startCol])):
		move = "N"
	}

	// canComeFrom is a simple map for used to prevent the pipe navigator from tracking backwards
	canComeFrom := map[string]string{
		"W": "E",
		"E": "W",
		"N": "S",
		"S": "N",
	}

	// Begin moving
	row, col := startRow, startCol
	for {
		// Prevent moving backwards
		cameFrom := canComeFrom[move]
		switch move {
		case "N":
			row--
		case "S":
			row++
		case "E":
			col++
		case "W":
			col--
		}
		pipeVal := string(lines[row][col])
		if pipeVal != "S" {
			nextMove := lo.Without(pipes[pipeVal].Dir, cameFrom)
			move = nextMove[0]
		}
		locs = append(locs, ez.Point{X: float64(col), Y: float64(row)})

		if row == startRow && col == startCol {
			break
		}
	}

	// Part 2
	if part2 {
		// Use Shoelace formula and Pick's theorem to generate the inner boundary points
		return ez.Picks(math.Abs(ez.Shoelace(locs)), len(locs))
	}

	// Part 1 asks for the furthest from the start, which is just half of the total number of locations along the path
	return len(locs) / 2
}
