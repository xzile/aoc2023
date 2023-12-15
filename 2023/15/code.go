package main

import (
	"aoc-in-go/ez"
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
	lines = strings.Split(lines[0], ",")
	sum := 0

	// Part 2
	if part2 {
		type Lens struct {
			Label string
			Focal int
		}
		type Box struct {
			Lenses []Lens
		}
		boxes := make([]Box, 256)

		for _, line := range lines {
			// Determine where - or = is in the string
			opIdx := max(strings.Index(line, "-"), strings.Index(line, "="))
			label := line[0:opIdx]
			boxNo := Hash(label)
			// Pointer so we can operate on the box lenses directly
			box := &boxes[boxNo]
			if box.Lenses == nil {
				box.Lenses = []Lens{}
			}
			op := line[opIdx : opIdx+1]
			switch op {
			case "-":
				// Remove matching label if found
				box.Lenses = lo.Reject(box.Lenses, func(item Lens, _ int) bool {
					return item.Label == label
				})
			case "=":
				// Add or replace based on label
				focal := ez.Atoi(line[opIdx+1:])
				found := false
				for i, lens := range box.Lenses {
					if lens.Label == label {
						found = true
						box.Lenses[i].Focal = focal
						break
					}
				}
				if !found {
					box.Lenses = append(box.Lenses, Lens{
						Label: label,
						Focal: focal,
					})
				}
			}
		}

		for i, box := range boxes {
			for j, lens := range box.Lenses {
				// The focusing power of a single lens is the result of multiplying together:
				// - One plus the box number of the lens in question.
				// - The slot number of the lens within the box: 1 for the first lens, 2 for the second lens, and so on.
				// - The focal length of the lens.
				out := (1 + i) * (j + 1) * lens.Focal
				sum += out
			}
		}

		return sum
	}

	// Part 1
	for _, line := range lines {
		sum += Hash(line)
	}

	return sum
}

func Hash(line string) int {
	seq := 0
	for _, char := range line {
		seq += int(char)
		seq *= 17
		seq = seq % 256
	}
	return seq
}
