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

	predictions := []int{}
	for _, line := range lines {
		// Convert line parts to a slice of ints
		vals := lo.Map(strings.Split(line, " "), func(item string, _ int) int {
			return ez.Atoi(item)
		})
		diffs := make([][]int, 0)
		// Line is our base [0], slice in our diffs
		diffs = append(diffs, vals)

		i := 1
		for {
			results := []int{}
			for j := 0; j < len(diffs[i-1])-1; j++ {
				// Each result in the new diff set is the difference between the next item (j+1) and the current item (j) from the previous row
				// Example, where "i" is row 2 being generated from the previous row (row 1)
				// row    j    j+1
				// 1      10   15
				// 2      5
				res := diffs[i-1][j+1] - diffs[i-1][j]
				results = append(results, res)
			}
			diffs = append(diffs, results)

			// Break if newest diff set is all 0s
			if lo.Every([]int{0}, diffs[i]) {
				break
			}

			i++
		}

		newPrediction := 0
		// Work from the next-to-last (before the all 0s line), to the original line
		for n := len(diffs) - 1; n >= 0; n-- {
			if part2 {
				// Part 2 uses the first item, and must compute the difference
				first := diffs[n][0]
				newPrediction = first - newPrediction
			} else {
				// Part 1 uses the last item and is computed using summation
				last, _ := lo.Last(diffs[n])
				newPrediction += last
			}
		}

		predictions = append(predictions, newPrediction)
	}

	return ez.Sum(predictions)
}
