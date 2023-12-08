package main

import (
	"aoc-in-go/ez"
	"bufio"
	"regexp"
	"slices"
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
	// Part 2
	if part2 {
		sc := bufio.NewScanner(strings.NewReader(input))
		symMap := map[int]map[int][]int{}
		lineNo := 0
		for sc.Scan() {
			lineNo++
			line := sc.Text()
			newLine := line
			if len(line) == 0 {
				continue
			}

			// Find all * symbols
			symMatch := regexp.MustCompile(`[\*]`)
			for _, symMatches := range symMatch.FindAllStringIndex(newLine, -1) {
				if len(symMatches) > 0 {
					if symMap[lineNo] == nil {
						symMap[lineNo] = map[int][]int{}
					}
					// Store symbol location by line number and location (column)
					// Eventually will add adjacent numbers to the slice of ints
					symMap[lineNo][symMatches[0]] = []int{}
				}
			}
		}

		sc = bufio.NewScanner(strings.NewReader(input))
		lineNo = 0
		sum := 0
		for sc.Scan() {
			lineNo++
			line := sc.Text()
			newLine := line
			if len(line) == 0 {
				continue
			}

			// Find numbers
			numMatch := regexp.MustCompile(`[\d]+`)
			for _, numMatches := range numMatch.FindAllStringIndex(newLine, -1) {
				// Given the location of the digits "look around" to the adjacent locations
				for i := lineNo - 1; i <= lineNo+1; i++ {
					for j := numMatches[0] - 1; j <= numMatches[1]; j++ {
						// If the location appears in the symMap, it's adjacent to a symbol
						if _, rowExists := symMap[i]; rowExists {
							if _, colExists := symMap[i][j]; colExists {
								// Append the number to the symbol it's adjacent to
								num := ez.Atoi(newLine[numMatches[0]:numMatches[1]])
								symMap[i][j] = append(symMap[i][j], num)
							}
						}
					}
				}
			}
		}

		// Loop through the symMap
		for _, rows := range symMap {
			for _, vals := range rows {
				// When exactly 2 numbers are adjacent to a symbol, multiply them together and add to the sum
				if len(vals) == 2 {
					sum += (vals[0] * vals[1])
				}
			}
		}

		return sum
	}

	// Part 1
	sc := bufio.NewScanner(strings.NewReader(input))
	symMap := map[int][]int{}
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := sc.Text()
		newLine := line
		if len(line) == 0 {
			continue
		}

		// Grab all the symbols and store their locations in symMap
		symMatch := regexp.MustCompile(`[^\d\.\s]`)
		for _, symMatches := range symMatch.FindAllStringIndex(newLine, -1) {
			if len(symMatches) > 0 {
				symMap[lineNo] = append(symMap[lineNo], symMatches[0])
			}
		}
	}

	sc = bufio.NewScanner(strings.NewReader(input))
	lineNo = 0
	sum := 0
	for sc.Scan() {
		lineNo++
		line := sc.Text()
		newLine := line
		if len(line) == 0 {
			continue
		}

		// Match consecutive digis
		numMatch := regexp.MustCompile(`[\d]+`)
		for _, numMatches := range numMatch.FindAllStringIndex(newLine, -1) {
			numPasses := false
			// Given the location of the digits "look around" to the adjacent locations
			for i := lineNo - 1; i <= lineNo+1; i++ {
				for j := numMatches[0] - 1; j <= numMatches[1]; j++ {
					// If the location appears in the symMap, it's adjacent to a symbol
					if _, rowExists := symMap[i]; rowExists {
						if slices.Contains(symMap[i], j) {
							numPasses = true
						}
					}
				}
			}

			// The number has an adjacent symbol; convert the digits to an int and add to the sum
			if numPasses {
				num := ez.Atoi(newLine[numMatches[0]:numMatches[1]])
				sum += num
			}
		}
	}

	return sum
}
