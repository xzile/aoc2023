package main

import (
	"aoc-in-go/ez"
	"container/list"
	"regexp"
	"slices"
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

type Cube struct {
	X int
	Y int
	Z int
}

type Brick struct {
	LineNo   int  // Used as a distinct label for a brick
	Fallen   int  // Negative number of how much depth Z fell
	Start    Cube // Only used for calculating "All", and is not updated during fall
	End      Cube // Only used for calculating "All", and is not updated during fall
	All      []Cube
	Supports []*Brick
	RestsOn  []*Brick
}

var reCube = regexp.MustCompile(`(\d+),(\d+),(\d+)~(\d+),(\d+),(\d+)`)

func run(part2 bool, input string) any {
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	// Convert lines to Bricks
	bricks := make([]*Brick, 0)
	for i, line := range lines {
		parts := reCube.FindStringSubmatch(line)
		brick := Brick{
			Fallen: 0,
			LineNo: i + 1,
			Start: Cube{
				X: ez.Atoi(parts[1]),
				Y: ez.Atoi(parts[2]),
				Z: ez.Atoi(parts[3]),
			},
			End: Cube{
				X: ez.Atoi(parts[4]),
				Y: ez.Atoi(parts[5]),
				Z: ez.Atoi(parts[6]),
			},
			All:      []Cube{},
			Supports: make([]*Brick, 0),
			RestsOn:  make([]*Brick, 0),
		}
		brick.SetAll()
		bricks = append(bricks, &brick)
	}

	// Sort bricks
	slices.SortStableFunc(bricks, func(a, b *Brick) int {
		// Use Z to determine the depth of the bricks
		// b - a to get a reverse sort, the lowest bricks will be at the end of the slice
		if a.Start.Z != b.Start.Z {
			return b.Start.Z - a.Start.Z
		}
		if a.End.Z != b.End.Z {
			return b.End.Z - a.End.Z
		}
		return 0
	})

	// Make them fall
	FallBricks(bricks)

	// Map of cube to brick will help us build the Supports and RestsOn slices later
	cubeMap := make(map[Cube]*Brick, 0)
	for i := 0; i < len(bricks); i++ {
		b := *bricks[i]
		for _, c := range b.All {
			cubeMap[c] = bricks[i]
		}
	}

	// Set supported and rests on per brick
	for i := 0; i < len(bricks); i++ {
		b := *bricks[i]
		for _, c := range b.All {
			// Look up the brick stack and determine which bricks support bricks above them
			c := c
			c.Z++
			if bs, ok := cubeMap[c]; ok {
				if bricks[i].LineNo != bs.LineNo { // Check to make sure we're not looking at ourself for vertical bricks
					bricks[i].Supports = append(bricks[i].Supports, bs)
				}
			}
			// Look down the brick stack and determine which bricks rest on bricks below them
			// Using the same "c", so, we need to reduce it by 2, since we increased it by 1 above
			c.Z -= 2
			if bs, ok := cubeMap[c]; ok {
				if bricks[i].LineNo != bs.LineNo { // Check to make sure we're not looking at ourself for vertical bricks
					bricks[i].RestsOn = append(bricks[i].RestsOn, bs)
				}
			}
		}
		// Make sure lists are unique, since support and rests on may have multiple touch points
		bricks[i].Supports = lo.UniqBy(bricks[i].Supports, func(item *Brick) int { return item.LineNo })
		bricks[i].RestsOn = lo.UniqBy(bricks[i].RestsOn, func(item *Brick) int { return item.LineNo })
	}

	// Part 2
	if part2 {
		sum := 0
		for i := 0; i < len(bricks); i++ {
			// Use a queue per brick to process "falling"
			q := list.New()
			// Falling is a list of bricks that are considered falling
			falling := make([]*Brick, 0)
			falling = append(falling, bricks[i])
			// Start with the desired brick to destroy
			q.PushBack(bricks[i])

			for q.Len() > 0 {
				bAny := q.Front()
				b := bAny.Value.(*Brick)
				// Check all bricks that current brick is supporting
				for j := 0; j < len(b.Supports); j++ {
					// Every brick the supported brick rests on must be falling to also be considered falling
					if lo.Every(falling, b.Supports[j].RestsOn) {
						q.PushBack(b.Supports[j])
						falling = append(falling, b.Supports[j])
					}
				}
				q.Remove(bAny)
			}

			// Uniq to account for multiple bricks supporting many bricks that each are considered as "falling"
			// - 1 due to the destroyed brick doesn't actually fall
			sum += len(lo.Uniq(falling)) - 1
		}

		return sum
	}

	sum := 0
	for i := 0; i < len(bricks); i++ {
		b := *bricks[i]
		// Any bricks that do not support other bricks are safe to destroy
		if len(b.Supports) == 0 {
			sum++
			continue
		} else {
			// We can only elimiate bricks that support other bricks if the supported bricks are supported by 2 or more bricks
			if lo.EveryBy(b.Supports, func(item *Brick) bool {
				return len(item.RestsOn) >= 2
			}) {
				sum++
				continue
			}
		}
	}

	return sum
}

// FallBricks moves the All slice cubes down on the Z axis until they can no longer move
func FallBricks(bricks []*Brick) {
	// Taken is a cache of brick space that is occupied
	taken := []Cube{}

	// Start with the lowest brick
	for i := len(bricks) - 1; i >= 0; i-- {
		// Get the cubes of that brick
		b := bricks[i]
		cubes := make([]Cube, len(b.All))
		copy(cubes, b.All)
		for {
			atLowest := false
			//Check if the Z of any of the cubes are at 1, if so, it's already at it's lowest point
			if !atLowest {
				for _, cube := range cubes {
					if cube.Z == 1 {
						atLowest = true
					}
				}
			}
			// Attempt to lower all cubes Z by -1
			if !atLowest {
				testCubes := make([]Cube, len(cubes))
				copy(testCubes, cubes)
				for j := range testCubes {
					testCubes[j].Z -= 1
				}
				if lo.Some(taken, testCubes) {
					atLowest = true
				} else {
					cubes = testCubes
				}
			}
			// If we're not at our lowest, we can update our cubes within the bricks
			if !atLowest {
				// Update the cubes of the brick
				b.All = cubes
				// Fallen could be used for adjusting the Z on the Start/End if it's needed
				b.Fallen--
			} else {
				break
			}
		}

		// Cache taken cubes
		for _, c := range cubes {
			c := c
			taken = append(taken, c)
		}
	}
}

// SetAll populates individual "cubes" in the b.All slice
func (b *Brick) SetAll() {
	// Always append the start
	b.All = append(b.All, b.Start)

	if b.Start.X != b.End.X {
		// X increments
		for x := b.Start.X + 1; x <= b.End.X; x++ {
			b.All = append(b.All, Cube{
				X: x,
				Y: b.Start.Y,
				Z: b.Start.Z,
			})
		}
	}
	if b.Start.Y != b.End.Y {
		// Y increments
		for y := b.Start.Y + 1; y <= b.End.Y; y++ {
			b.All = append(b.All, Cube{
				X: b.Start.X,
				Y: y,
				Z: b.Start.Z,
			})
		}
	}
	if b.Start.Z != b.End.Z {
		// Z increments
		for z := b.Start.Z + 1; z <= b.End.Z; z++ {
			b.All = append(b.All, Cube{
				X: b.Start.X,
				Y: b.Start.Y,
				Z: z,
			})
		}
	}
}
