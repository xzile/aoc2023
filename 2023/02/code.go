package main

import (
	"bufio"
	"strconv"
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
	// when you're ready to do part 2, remove this "not implemented" block
	if part2 {
		sum := 0
		sc := bufio.NewScanner(strings.NewReader(input))
		for sc.Scan() {
			line := sc.Text()
			newLine := line
			if len(line) == 0 {
				continue
			}

			// Split on first :, which is Game XXX: ...
			gameParts := strings.SplitN(newLine, ":", 2)
			_, gamePulls := gameParts[0], gameParts[1]
			// Initialize got's to 1, so multiplication with 0 doesn't cause issues
			got := map[string]int{
				"red":   1,
				"green": 1,
				"blue":  1,
			}

			// Split each game's pulls into invidiual pulls
			pulls := strings.Split(gamePulls, ";")
			for _, pull := range pulls {
				for _, part := range strings.Split(pull, ",") {
					pullParts := strings.Split(strings.TrimSpace(part), " ")
					numPulled, _ := strconv.Atoi(pullParts[0])
					colorPulled := pullParts[1]

					// Only set the got if the numPulled is greater than the existing got
					if numPulled > got[colorPulled] {
						got[colorPulled] = numPulled
					}
				}
			}

			sum += (got["red"] * got["green"] * got["blue"])
		}

		return sum
	}

	// solve part 1 here
	desired := map[string]int{
		"red":   12,
		"green": 13,
		"blue":  14,
	}
	// games := map[string]Game{}
	sum := 0
	sc := bufio.NewScanner(strings.NewReader(input))
	for sc.Scan() {
		line := sc.Text()
		newLine := line
		if len(line) == 0 {
			continue
		}

		// Split on first :, which is Game XXX: ...
		gameParts := strings.SplitN(newLine, ":", 2)
		gameRef, gamePulls := gameParts[0], gameParts[1]
		// Get the game ID as an integer
		gameID, _ := strconv.Atoi(strings.Split(gameRef, " ")[1])
		gamePass := true

		// Split each game's pulls into invidiual pulls
		pulls := strings.Split(gamePulls, ";")
		for _, pull := range pulls {
			got := map[string]int{
				"red":   0,
				"green": 0,
				"blue":  0,
			}
			for _, part := range strings.Split(pull, ",") {
				// part example: 3 blue
				pullParts := strings.Split(strings.TrimSpace(part), " ")
				numPulled, _ := strconv.Atoi(pullParts[0])
				colorPulled := pullParts[1]

				got[colorPulled] += numPulled
			}
			if got["red"] > desired["red"] ||
				got["green"] > desired["green"] ||
				got["blue"] > desired["blue"] {
				// If any game got more than the desired max, fail the game
				gamePass = false
				// fmt.Println(fmt.Sprintf("Game %d failed: got %v", gameID, got))
			}
		}

		if gamePass {
			sum += gameID
		}
	}

	return sum
}
