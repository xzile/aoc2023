package main

import (
	"aoc-in-go/ez"
	"github.com/jpillora/puzzler/harness/aoc"
	"golang.org/x/exp/maps"
	"regexp"
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
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	steps := strings.TrimSpace(lines[0])
	network := map[string]Node{}
	networkParser := regexp.MustCompile(`(\w+) = \((\w+), (\w+)\)`)
	for _, v := range lines[2:] {
		v := v
		parts := networkParser.FindAllStringSubmatch(v, -1)
		if len(parts) == 0 {
			continue
		}
		network[parts[0][1]] = Node{
			L: parts[0][2],
			R: parts[0][3],
		}
	}

	// Part 2
	if part2 {
		// Find all paths starting with A
		var paths []string
		for _, v := range maps.Keys(network) {
			v := v
			if string(v[len(v)-1:]) == "A" {
				paths = append(paths, v)
			}
		}

		// Since the challenge stated that:
		// the number of nodes with names ending in `A` is equal to the number ending in `Z`
		// We can assume each A path will only ever lead to a single Z path, after some number of steps
		// Calculate the number of steps needed for each A to reach it's Z
		var pathSteps []int
		for _, v := range paths {
			node := network[v]
			stepsTaken := 0
			for {
				stepToTake := stepsTaken % len(steps)
				which := steps[stepToTake]
				stepsTaken++

				var nextStep string
				if string(which) == "L" {
					nextStep = node.L
				} else {
					nextStep = node.R
				}

				// Exit condition, is that our next step ends in Z
				if string(nextStep[len(nextStep)-1:]) == "Z" {
					pathSteps = append(pathSteps, stepsTaken)
					break
				}
				node = network[nextStep]
			}
		}

		// The least common multiple of all steps will be when all paths step will end in Z
		return ez.LCM(pathSteps[0], pathSteps[1], pathSteps[2:]...)
	}

	// Part 1
	stepsTaken := 0
	node := network["AAA"]
	for {
		stepToTake := stepsTaken % len(steps)
		which := steps[stepToTake]
		stepsTaken++

		var nextStep string
		if string(which) == "L" {
			nextStep = node.L
		} else {
			nextStep = node.R
		}

		// Exit condition, is that our next step is ZZZ
		if nextStep == "ZZZ" {
			break
		}

		node = network[nextStep]
	}

	return stepsTaken
}

type Node struct {
	L string
	R string
}
