package main

import (
	"fmt"
	"github.com/jpillora/puzzler/harness/aoc"
	"golang.org/x/exp/maps"
	"slices"
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
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	// Parse games
	var games []Game
	for i := range lines {
		parts := strings.Split(lines[i], " ")
		if len(parts) != 2 {
			continue
		}

		games = append(games, Game{
			Cards: parts[0],
			Bid:   Atoi(parts[1]),
		})
	}

	// Part 1
	CardValue := map[string]int{
		"A": 13,
		"K": 12,
		"Q": 11,
		"J": 10,
		"T": 9,
		"9": 8,
		"8": 7,
		"7": 6,
		"6": 5,
		"5": 4,
		"4": 3,
		"3": 2,
		"2": 1,
	}

	if part2 {
		CardValue = map[string]int{
			"A": 13,
			"K": 12,
			"Q": 11,
			"J": 0, // Jokers are lower than 2
			"T": 9,
			"9": 8,
			"8": 7,
			"7": 6,
			"6": 5,
			"5": 4,
			"4": 3,
			"3": 2,
			"2": 1,
		}

	}

	// Sort games
	slices.SortFunc(games, func(a, b Game) int {
		aType := a.Type(part2)
		bType := b.Type(part2)

		// Types are different, we can sort just on the type
		if aType != bType {
			return aType - bType
		}

		// Types are the same, compare each card
		for i := range a.Cards {
			aCard := CardValue[string(a.Cards[i])]
			bCard := CardValue[string(b.Cards[i])]

			// Same card, skip to the next
			if aCard == bCard {
				continue
			}

			// Sort on the card values
			return aCard - bCard
		}

		return 0
	})

	// Games are sorted from weakest to strongest, generate winnings
	sum := 0
	for i := range games {
		sum += (i + 1) * games[i].Bid
	}
	return sum
}

// Game contains a set of cards, as well as a bid
type Game struct {
	// Cards is a string of 5 characters, representing the cards in the game
	Cards string
	Bid   int
}

// Jokers returns the count of jokers, used in part 2
func (g Game) Jokers() int {
	cardCount := g.CardCount()
	return cardCount["J"]
}

// CardCount returns a map that is keyed by the Card, and a value of the count of that card within the game
func (g Game) CardCount() map[string]int {
	counts := map[string]int{}
	for i := range g.Cards {
		counts[string(g.Cards[i])]++
	}

	return counts
}

// Type returns an integer, 7 = strongest, 1 = weakest, representing the strength of a game's cards
func (g Game) Type(part2 bool) int {
	cardCount := g.CardCount()

	if part2 {
		jokers := g.Jokers()
		if jokers == 5 {
			// Five of a kind with just jokers
			return 7
		} else if jokers > 0 {
			// Remove jokers from cardCount
			delete(cardCount, "J")

			// Add the jokers to the most repeating card
			maxCardCount := slices.Max(maps.Values(cardCount))
			for k := range cardCount {
				if cardCount[k] == maxCardCount {
					cardCount[k] += jokers
					break
				}
			}
		}
	}

	switch {
	case len(cardCount) == 1: // Five of a kind
		return 7
	case len(cardCount) == 2: // Four of a kind or full house
		if slices.Max(maps.Values(cardCount)) == 4 { // four of a kind
			return 6
		}
		// full house
		return 5
	case len(cardCount) == 3: // three of a kind or two pair
		if slices.Max(maps.Values(cardCount)) == 3 { // three of a kind
			return 4
		}
		return 3
	case len(cardCount) == 4: // one pair
		return 2
	default: // High card
		return 1
	}
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
