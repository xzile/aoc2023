package main

import (
	"aoc-in-go/ez"
	"fmt"
	"math"
	"regexp"
	"strconv"
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

type Dig struct {
	ez.Point
	Hex string
}

var re = regexp.MustCompile(`(\w+) (\d+) \(#(\w+)\)`)

func run(part2 bool, input string) any {
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	var digs []Dig

	// Start on an imaginary graph at the center: 0,0
	startPoint := ez.Point{X: 0.0, Y: 0.0}
	digs = append(digs, Dig{
		Point: startPoint,
		Hex:   "",
	})

	prevPoint := startPoint
	for i := range lines {
		var dist int
		var dir string
		var hex string

		parts := re.FindStringSubmatch(lines[i])
		hex = parts[3]

		if !part2 {
			// Part 1
			dist = ez.Atoi(parts[2])
			dir = parts[1]
		} else {
			// Part 2
			distParsed, _ := strconv.ParseInt(hex[0:5], 16, 0)
			dist = int(distParsed)

			switch hex[5:] {
			case "0":
				dir = "R"
			case "1":
				dir = "D"
			case "2":
				dir = "L"
			case "3":
				dir = "U"
			}
		}

		// Add individual dig increments
		// This is not performant and brute-force, but necessary for correctness as the trench may overlap itself
		for j := 1; j <= dist; j++ {
			nextPoint := NextPoint(dir, 1, prevPoint)
			digs = append(digs, Dig{
				Point: nextPoint,
				Hex:   hex,
			})
			prevPoint = nextPoint
		}
	}

	// Get unique trench points, this is the perimeter
	// This felt necessary to account for when the trench digging overlaps with a previously dug trench
	uniqTrenches := lo.UniqBy(digs, func(item Dig) string {
		return fmt.Sprintf("%v,%v", item.Point.X, item.Point.Y)
	})
	// Calculate the inner area using shoelace
	locs := lo.Map(digs, func(item Dig, _ int) ez.Point {
		return item.Point
	})
	inner := ez.Picks(math.Abs(ez.Shoelace(locs)), len(locs))

	// combine trench points with inner shoelace area
	return int64(float64(len(uniqTrenches)) + inner)
}

func NextPoint(dir string, dist float64, prevPoint ez.Point) ez.Point {
	xMod := .0
	yMod := .0

	switch dir {
	case "U":
		yMod += dist
	case "D":
		yMod -= dist
	case "L":
		xMod -= dist
	case "R":
		xMod += dist
	}

	return ez.Point{X: prevPoint.X + xMod, Y: prevPoint.Y + yMod}
}
