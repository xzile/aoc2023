package main

import (
	"errors"
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

type Stone struct {
	LineNo int
	Point  Point
	Line   Line
	VX     float64
	VY     float64
	VZ     float64
}

type Point struct {
	X float64
	Y float64
	Z float64
}

type Line struct {
	Slope float64
	Yint  float64
}

func CreateLine(a, b Point) Line {
	slope := (b.Y - a.Y) / (b.X - a.X)
	yint := a.Y - slope*a.X
	return Line{slope, yint}
}

func EvalX(l Line, x float64) float64 {
	return l.Slope*x + l.Yint
}

func Intersection(l1, l2 Line) (Point, error) {
	if l1.Slope == l2.Slope {
		return Point{}, errors.New("The lines do not intersect")
	}
	x := (l2.Yint - l1.Yint) / (l1.Slope - l2.Slope)
	y := EvalX(l1, x)
	return Point{X: x, Y: y}, nil
}

var reStone = regexp.MustCompile(`(\d+),\s+(\d+),\s+(\d+)\s+@\s+([-]?\d+),\s+([-]?\d+),\s+([-]?\d+)`)

func run(part2 bool, input string) any {
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	stones := []Stone{}
	for i, line := range lines {
		parts := reStone.FindStringSubmatch(line)
		stone := Stone{
			LineNo: i + 1,
			Point: Point{
				X: lo.Must(strconv.ParseFloat(parts[1], 64)),
				Y: lo.Must(strconv.ParseFloat(parts[2], 64)),
				Z: lo.Must(strconv.ParseFloat(parts[3], 64)),
			},
			VX: lo.Must(strconv.ParseFloat(parts[4], 64)),
			VY: lo.Must(strconv.ParseFloat(parts[5], 64)),
			VZ: lo.Must(strconv.ParseFloat(parts[6], 64)),
		}
		stone.Line = CreateLine(stone.Point, Point{
			X: stone.Point.X + stone.VX,
			Y: stone.Point.Y + stone.VY,
		})
		stones = append(stones, stone)
	}

	// Part 2
	// ref: https://github.com/rumkugel13/AdventOfCode2023/blob/main/day24.go
	// Originally solved via Z3, but used the Go / Math based approach in the above to have a Go-only solution
	if part2 {
		if len(lines) < 20 {
			// Example not supported
			return 0
		}

		// Effectively, we're looking for stones that have matching velocities in X, Y, and Z, each evaluated independently
		// These stones are travelling in parallel in that given direction and will help solve the coordinates
		maybeX, maybeY, maybeZ := []int{}, []int{}, []int{}
		for i := 0; i < len(stones)-1; i++ {
			for j := i + 1; j < len(stones); j++ {
				a, b := stones[i], stones[j]
				if a.VX == b.VX {
					nextMaybe := findMatchingVel(int(b.Point.X-a.Point.X), int(a.VX))
					if len(maybeX) == 0 {
						maybeX = nextMaybe
					} else {
						maybeX = lo.Intersect(maybeX, nextMaybe)
					}
				}
				if a.VY == b.VY {
					nextMaybe := findMatchingVel(int(b.Point.Y-a.Point.Y), int(a.VY))
					if len(maybeY) == 0 {
						maybeY = nextMaybe
					} else {
						maybeY = lo.Intersect(maybeY, nextMaybe)
					}
				}
				if a.VZ == b.VZ {
					nextMaybe := findMatchingVel(int(b.Point.Z-a.Point.Z), int(a.VZ))
					if len(maybeZ) == 0 {
						maybeZ = nextMaybe
					} else {
						maybeZ = lo.Intersect(maybeZ, nextMaybe)
					}
				}
			}
		}

		var out = 0
		if len(maybeX) == len(maybeY) && len(maybeY) == len(maybeZ) && len(maybeZ) == 1 {
			// only one possible velocity in all dimensions
			// rockVel is a Point, just to capture the X, Y, and Z into a struct, it's not an actual point
			rockVel := Point{X: float64(maybeX[0]), Y: float64(maybeY[0]), Z: float64(maybeZ[0])}
			// The stones don't really matter, we just need two to calculate
			A, B := stones[0], stones[1]
			// Lots of math
			// ref: https://www.reddit.com/r/adventofcode/comments/18pnycy/comment/keqf8uq/
			mA := (A.VY - rockVel.Y) / (A.VX - rockVel.X)
			mB := (B.VY - rockVel.Y) / (B.VX - rockVel.X)
			cA := A.Point.Y - (mA * A.Point.X)
			cB := B.Point.Y - (mB * B.Point.X)
			xPos := (cB - cA) / (mA - mB)
			yPos := mA*xPos + cA
			t := (xPos - A.Point.X) / (A.VX - rockVel.X)
			zPos := A.Point.Z + (A.VZ-rockVel.Z)*t
			out = int(xPos + yPos + zPos)
		}

		return out
	}

	boxMin := float64(200000000000000)
	boxMax := float64(400000000000000)
	if len(lines) < 20 {
		// Example data
		boxMin = float64(7)
		boxMax = float64(27)
	}

	count := 0

	for i := range stones {
		a := stones[i]
		for j := i + 1; j < len(stones); j++ {
			b := stones[j]

			p, err := Intersection(a.Line, b.Line)
			if err != nil {
				// Do not intersect
				continue
			}

			if !(p.X >= boxMin && p.X <= boxMax && p.Y >= boxMin && p.Y <= boxMax) {
				// Not within test area
				continue
			}

			// Check for "past" at intersection
			if a.PointInFuture(p) && b.PointInFuture(p) {
				count++
			}
		}
	}

	// solve part 1 here
	return count
}

func (s Stone) PointInFuture(p Point) bool {
	inFuture := true
	if s.VX < 0 && p.X > s.Point.X {
		inFuture = false
	} else if s.VX > 0 && p.X < s.Point.X {
		inFuture = false
	}
	if s.VY < 0 && p.Y > s.Point.Y {
		inFuture = false
	} else if s.VY > 0 && p.Y < s.Point.Y {
		inFuture = false
	}
	return inFuture
}

func findMatchingVel(dvel, pv int) []int {
	match := []int{}
	for v := -1000; v < 1000; v++ {
		if v != pv && dvel%(v-pv) == 0 {
			match = append(match, v)
		}
	}
	return match
}
