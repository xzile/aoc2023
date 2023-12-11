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

// on code change, run will be executed 4 times:
// 1. with: false (part1), and example input
// 2. with: true (part2), and example input
// 3. with: false (part1), and user input
// 4. with: true (part2), and user input
// the return value of each run is printed to stdout
func run(part2 bool, input string) any {
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	repeat := 1
	if part2 {
		// Since the "repeated" row replaces the existing row, we still end up counting the existing row
		// So, we subtract the desired repeat by 1
		repeat = 999999
	}

	// Mark rows and columns that need to be considered repeated
	repeatRows := []int{}
	for i, line := range lines {
		if !strings.Contains(line, "#") {
			repeatRows = append(repeatRows, i)
		}
	}

	repeatCols := []int{}
	for i := len(lines[0]) - 1; i >= 0; i-- {
		expandCol := true
		for _, line := range lines {
			if string(line[i]) == "#" {
				expandCol = false
			}
		}
		if expandCol {
			repeatCols = append(repeatCols, i)
		}
	}

	// Chart the points, using column/j/X and row/i/Y
	var points []ez.Point
	for i, line := range lines {
		for j, c := range line {
			if c == '#' {
				points = append(points, ez.Point{X: float64(j), Y: float64(i)})
			}
		}
	}

	sum := .0
	for i, p1 := range points {
		for j, p2 := range points {
			// We only need unique pairs, which will only be points that are greater than the existing point
			if j > i {
				// Repeats are those rows or columns that pass over one of the rows or columns that are marked as repeating
				xDist := math.Abs(p1.X - p2.X)
				xRepeats := lo.Filter(repeatCols, func(item, index int) bool {
					a := max(p1.X, p2.X)
					b := min(p1.X, p2.X)
					return float64(item) < a && float64(item) > b
				})
				yDist := math.Abs(p1.Y - p2.Y)
				yRepeats := lo.Filter(repeatRows, func(item, index int) bool {
					a := max(p1.Y, p2.Y)
					b := min(p1.Y, p2.Y)
					return float64(item) < a && float64(item) > b
				})
				// Sum up the distances and add on any repeating rows/columns
				sum += xDist + yDist + float64((len(xRepeats)+len(yRepeats))*repeat)
			}
		}
	}

	return int64(sum)
}
