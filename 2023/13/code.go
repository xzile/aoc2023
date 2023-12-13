package main

import (
	"strings"

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
func run(part2 bool, input string) any {
	grids := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n\n")

	// Part 2
	if part2 {
		sum := 0
		for _, grid := range grids {
			lineLen, lineCount := GridDim(grid)
		gridLoop:
			for i := 0; i < lineCount; i++ {
				for j := 0; j < lineLen; j++ {
					// smudgeGrid has the character at i,j coordinates flipped
					smudgeGrid := Smudge(i, j, grid)

					lines := strings.Split(strings.ReplaceAll(smudgeGrid, "\r\n", "\n"), "\n")
					cols := func() []string {
						out := make([]string, len(lines[0]))
						for i := 0; i < len(lines[0]); i++ {
							for j := 0; j < len(lines); j++ {
								out[i] += string(lines[j][i])
							}
						}
						return out
					}()

					top, bottom := Mirror(lines, i)
					left, right := Mirror(cols, j)
					// checkRange is the number of lines that were used to evaluate the mirror
					checkRange := min(top, bottom)
					checkRangeLR := min(left, right)
					// If we have a mirrored set and the smudge is part of the mirrored set
					if top != 0 && bottom != 0 && i+1 > top-checkRange && i+1 <= top+checkRange {
						sum += 100 * top
						break gridLoop
					} else if left != 0 && right != 0 && j+1 > left-checkRangeLR && j+1 <= left+checkRangeLR {
						sum += left
						break gridLoop
					}
				}
			}
		}
		return sum
	}

	// Part 1
	sum := 0
	for _, grid := range grids {
		lines := strings.Split(strings.ReplaceAll(grid, "\r\n", "\n"), "\n")
		cols := func() []string {
			out := make([]string, len(lines[0]))
			for i := 0; i < len(lines[0]); i++ {
				for j := 0; j < len(lines); j++ {
					out[i] += string(lines[j][i])
				}
			}
			return out
		}()
		top, bottom := Mirror(lines, 0)
		if top == 0 || bottom == 0 {
			left, _ := Mirror(cols, 0)
			sum += left
		} else {
			sum += 100 * top
		}
	}
	return sum
}

// Smudge flips the character at row i, in character j (0 indexed), on the grid and returns the smudged grid
func Smudge(i, j int, grid string) string {
	lines := strings.Split(strings.ReplaceAll(grid, "\r\n", "\n"), "\n")
	currentChar := string(lines[i][j])
	replace := "#"
	if currentChar == "#" {
		replace = "."
	}
	lines[i] = lines[i][:j] + replace + lines[i][j+1:]
	return strings.Join(lines, "\n")
}

// GridDim just returns the dimenions of the grid as: Line Number, Character Count
func GridDim(grid string) (int, int) {
	lines := strings.Split(strings.ReplaceAll(grid, "\r\n", "\n"), "\n")
	return len(lines[0]), len(lines)
}

// Mirror checks the lines for mirrored set, returning the top and bottom count of rows
// Provide smudge to skip lines that may cause the Mirror to return a match that doesn't contain the smuded row
func Mirror(lines []string, smudge int) (int, int) {
	for i := range lines {
		// Smudge checker for part 2
		// skip evaluating lines that are before the smudge is introduced
		if i < smudge {
			continue
		}
		isMirror := true
		for j := i + 1; j < len(lines); j++ {
			// Prevent looking back beyond the first line
			if i-(j-i)+1 < 0 {
				break
			}
			if lines[j] != lines[i-(j-i)+1] {
				isMirror = false
			}
		}
		if isMirror {
			return i + 1, len(lines) - (i + 1)
		}
	}
	// We should never get here, but pass 0 in case we do...
	return 0, len(lines)
}
