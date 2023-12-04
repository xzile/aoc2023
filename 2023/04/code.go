package main

import (
	"bufio"
	"regexp"
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
	// Regex for input-example
	// re := regexp.MustCompile(`Card\s+\d+:\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+\|\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)`)
	// Regex for input-user
	re := regexp.MustCompile(`Card\s+\d+:\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+\|\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)`)

	// Part 2
	if part2 {
		sc := bufio.NewScanner(strings.NewReader(input))
		lineNo := 0
		// Copies is indexed by line number, and will hold the count of copies that line number has
		copies := map[int]int{}
		sum := 0
		for sc.Scan() {
			lineNo++
			line := sc.Text()
			newLine := line
			if len(line) == 0 {
				continue
			}

			matches := re.FindAllStringSubmatch(newLine, -1)
			for _, match := range matches {
				winners := match[1:11]
				cards := match[11:]
				// Intersect is a copy of all cards that exist within the winners set
				intersect := lo.Intersect(winners, cards)
				// Increment the copies of future card sets given the original winner count
				for i := 1; i <= len(intersect); i++ {
					copies[lineNo+i] += 1
				}
				// Using any copies of the current original, increment the copies of future card sets given the original winner count
				for j := 1; j <= copies[lineNo]; j++ {
					for i := 1; i <= len(intersect); i++ {
						copies[lineNo+i] += 1
					}
				}
			}
			// Add the original + any copies to the total count
			sum += 1 + copies[lineNo]
		}

		return sum
	}

	// Part 1
	sc := bufio.NewScanner(strings.NewReader(input))
	sum := 0
	for sc.Scan() {
		line := sc.Text()
		newLine := line
		if len(line) == 0 {
			continue
		}

		matches := re.FindAllStringSubmatch(newLine, -1)
		for _, match := range matches {
			winners := match[1:11]
			cards := match[11:]
			// Intersect is a copy of all cards that exist within the winners set
			intersect := lo.Intersect(winners, cards)
			if len(intersect) > 0 {
				count := 1
				// Double for all winning cards after the first
				for i := 2; i <= len(intersect); i++ {
					count *= 2
				}
				sum += count
			}
		}
	}

	return sum
}
