package main

import (
	"fmt"
	"github.com/jpillora/puzzler/harness/aoc"
	"regexp"
	"strconv"
	"strings"
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
	onlyNums := regexp.MustCompile(`(\d+)`)
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	times := onlyNums.FindAllString(lines[0], -1)
	dists := onlyNums.FindAllString(lines[1], -1)

	if part2 {
		// For part 2, just collapse all the times & distances down to a single value
		times = []string{strings.Join(times, "")}
		dists = []string{strings.Join(dists, "")}
	}

	out := 1
	for raceNo, timeStr := range times {
		beats := 0
		time := Atoi(timeStr)
		dist := Atoi(dists[raceNo])

		// holdTime is how the button is held, but is also the "speed" of the boat
		for holdTime := 1; holdTime <= time; holdTime++ {
			// moveTime represents how much time is left in the race
			moveTime := time - holdTime
			// holdTime*moveTime is the total distance covered by the boat for the race
			if holdTime*moveTime > dist {
				beats++
			}
		}
		out *= beats
	}

	return out
}

// Atoi ignores the errors in strconv.Atoi and returns the response
func Atoi(in string) int {
	out, _ := strconv.Atoi(in)
	return out
}

// Log is a simple wrapper around Println + Sprintf
func Log(format string, a ...any) {
	fmt.Println(fmt.Sprintf(format, a...))
}
