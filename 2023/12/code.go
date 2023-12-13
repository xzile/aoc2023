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
	sum := int64(0)
	for _, line := range lines {
		parts := strings.Split(line, " ")
		pattern := parts[0]
		setsStr := parts[1]

		if part2 {
			// Part 2
			pattern = strings.Join(Dupe(pattern, 5), "?")
			setsStr = strings.Join(Dupe(setsStr, 5), ",")
		}

		setsStrs := strings.Split(setsStr, ",")
		sets := lo.Map(setsStrs, func(item string, _ int) int {
			return ez.Atoi(item)
		})

		sum += int64(ProcessPattern(pattern, sets))
	}

	return sum
}

func Dupe(parts string, dupes int) []string {
	out := []string{}
	for i := 0; i < dupes; i++ {
		out = append(out, parts)
	}
	return out
}

func ProcessPattern(pattern string, set []int) int {
	// Pre-build a cache to track
	// The first slice represents the characters in the pattern
	// The second slice represents whether or not the character can fulfill a given set
	var cache [][]int
	for i := 0; i < len(pattern); i++ {
		// Append a extra set that'll never get matched to provide an exit from our recursion
		cache = append(cache, make([]int, len(set)+1))
		for j := 0; j < len(set)+1; j++ {
			cache[i][j] = -1
		}
	}

	// Kick off the call to Iter, which will eventually return a result
	return Iter(0, 0, pattern, set, cache)
}

func Iter(i, j int, pattern string, set []int, cache [][]int) int {
	// Exit conditions for recursion
	// We've made it to the end of the pattern
	if i >= len(pattern) {
		// If we've not fulfilled all of our sets, do not consider the iteration a success
		if j < len(set) {
			return 0
		}
		return 1
	}

	// Check the cache for a given character/set having already been evaluated
	// If so, return that count
	if cache[i][j] != -1 {
		return cache[i][j]
	}

	result := 0
	if pattern[i] == '.' {
		// If '.', move on to the next character, result doesn't change
		result = Iter(i+1, j, pattern, set, cache)
	} else {
		// If '?', sum the result, as the ? is a wildcard and may fulfill a set
		if pattern[i] == '?' {
			result += Iter(i+1, j, pattern, set, cache)
		}
		// Recursion exit condition
		if j < len(set) {
			count := 0
			// Starting at the character, examine the remainder of the pattern
			for k := i; k < len(pattern); k++ {
				// Stop counting when
				// 1: we've counted higher than the set's size
				// OR
				// 2: we hit a '.'
				// OR
				// 3: (We've matched the set size AND the last character is a wildcard)
				if count > set[j] || pattern[k] == '.' || count == set[j] && pattern[k] == '?' {
					break
				}
				count += 1
			}

			// If our count matches, it can fulfill the set
			if count == set[j] {
				// If we're not at the end of the pattern
				// And the next character is a ?
				if i+count < len(pattern) && pattern[i+count] == '?' {
					// We can safely skip a character (treat it as a .) the set + an extra character and move to the next set
					result += Iter(i+count+1, j+1, pattern, set, cache)
				} else {
					// Otherwise, continue on evaluating the next range of character
					result += Iter(i+count, j+1, pattern, set, cache)
				}
			}
		}
	}

	// Store the result
	// Which is the number of sets a given character can fulfill
	cache[i][j] = result

	return result
}
