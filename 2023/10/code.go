package main

import (
	"fmt"
	"strconv"
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

var rowCount, colCount int

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
	rowCount = len(lines)
	colCount = len(lines[0])
	// Locs will represent a stored `row,col` location on the path
	locs := map[string]int{}

	// Find start
	var startRow, startCol int
	for i, line := range lines {
		if startCol = strings.Index(line, "S"); startCol != -1 {
			startRow = i
			break
		}
	}
	RecordLocation(locs, startRow, startCol)

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
	startMove := move

	// canComeFrom is a simple map for used to prevent the pipe navigator from tracking backwards
	canComeFrom := map[string]string{
		"W": "E",
		"E": "W",
		"N": "S",
		"S": "N",
	}

	// Begin moving
	row, col := startRow, startCol
	steps := 0
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
		RecordLocation(locs, row, col)

		steps++
		if row == startRow && col == startCol {
			break
		}
	}

	// Part 2
	if part2 {
		// For part 2, we're going to imagine we're re-walking the path
		// If the pipe is adjacent to an unused path pipe, we can count that and other adjacent (in a straight line)
		// until we hit a pipe that's used in the path, or the boundary of the grid
		// If we hit the boundary, we can consider that "side" to be outside and mark it as failed
		inner1Outside := false
		inner2Outside := false
		inner1Locs := map[string]int{}
		inner2Locs := map[string]int{}
		// Begin moving
		row, col := startRow, startCol
		move := startMove

		// Track the "sides" of the path
		// Depending on whether we're moving NS or EW, we'll begin by checking the opposite set of cardinal directions
		var inner1, inner2 string
		if startMove == "N" || startMove == "S" {
			inner1, inner2 = "W", "E"
		} else {
			inner1, inner2 = "N", "S"
		}

		// turner is a helper map for changing the "sides"
		// Given a particular pipe bend, the sides will shift directions based on which "side" it was previously
		turner := map[string]map[string]string{
			// "L": Pipe{Dir: []string{"N", "E"}},
			"L": {
				"N": "E",
				"S": "W",
				"E": "N",
				"W": "S",
			},
			// "J": Pipe{Dir: []string{"N", "W"}},
			"J": {
				"N": "W",
				"S": "E",
				"E": "S",
				"W": "N",
			},
			// "7": Pipe{Dir: []string{"S", "W"}},
			"7": {
				"N": "E",
				"S": "W",
				"E": "N",
				"W": "S",
			},
			// "F": Pipe{Dir: []string{"S", "E"}},
			"F": {
				"N": "W",
				"S": "E",
				"E": "S",
				"W": "N",
			},
		}

		for {
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

			// Check for unused in the direction of inner
			if !inner1Outside {
				err := RecordInside(row, col, inner1, inner1Locs, locs)
				if err != nil {
					inner1Outside = true
				}
			}
			if !inner2Outside {
				err := RecordInside(row, col, inner2, inner2Locs, locs)
				if err != nil {
					inner2Outside = true
				}
			}

			// Adjust inner based on turns
			if fromTo, exists := turner[pipeVal]; exists {
				inner1 = fromTo[inner1]
				inner2 = fromTo[inner2]

				// Recheck after the turn
				// Check for unused in the direction of inner
				if !inner1Outside {
					err := RecordInside(row, col, inner1, inner1Locs, locs)
					if err != nil {
						inner1Outside = true
					}
				}
				if !inner2Outside {
					err := RecordInside(row, col, inner2, inner2Locs, locs)
					if err != nil {
						inner2Outside = true
					}
				}
			}

			if row == startRow && col == startCol {
				break
			}
		}

		winner := inner1Locs
		if inner1Outside {
			winner = inner2Locs
		}
		return len(winner)
	}

	// Part 1 asks for the furthest from the start, which is just half of the total number of locations along the path
	return len(locs) / 2
}

func RecordLocation(locs map[string]int, startRow, startCol int) {
	locs[strconv.Itoa(startRow)+","+strconv.Itoa(startCol)] = 1
}

func LocExists(locs map[string]int, startRow, startCol int) bool {
	_, exists := locs[strconv.Itoa(startRow)+","+strconv.Itoa(startCol)]
	return exists
}

func RecordInside(row, col int, insideDir string, insideLocs, pathLocs map[string]int) error {
	for {
		switch insideDir {
		case "N":
			row--
		case "S":
			row++
		case "E":
			col++
		case "W":
			col--
		}

		if !LocExists(pathLocs, row, col) && row >= 0 && row < rowCount && col >= 0 && col < colCount {
			// If we ever hit a grid boundary, we know we're using the wrong inner
			if row == 0 || col == 0 || row == rowCount || col == colCount {
				return fmt.Errorf("direction out of bounds")
			}
			RecordLocation(insideLocs, row, col)
		} else {
			break
		}
	}

	return nil
}
