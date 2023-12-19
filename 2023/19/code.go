package main

import (
	"aoc-in-go/ez"
	"regexp"
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

type Condition struct {
	Piece    string
	Comp     string
	Val      int
	Terminal string
}

type Rule struct {
	Conditions []Condition
	Terminal   string
}

type Part struct {
	Rating map[string]int
	Total  int
}

type Workflow map[string]Rule

var reRule = regexp.MustCompile(`([xmas]+)([<>])(\d+):(\w+)`)
var rePart = regexp.MustCompile(`{([xmas])=(\d+),([xmas])=(\d+),([xmas])=(\d+),([xmas])=(\d+)}`)

func run(part2 bool, input string) any {
	rulesAndParts := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n\n")

	// Parse Rules within the workflows
	workflows := make(Workflow)
	for _, rawRule := range strings.Split(rulesAndParts[0], "\n") {
		rawRule = strings.TrimRight(rawRule, "}")
		labelAndRules := strings.Split(rawRule, "{")
		label := labelAndRules[0]
		rulePieces := strings.Split(labelAndRules[1], ",")
		ruleTerminal := rulePieces[len(rulePieces)-1:]

		conditions := make([]Condition, 0)
		for _, rulePiece := range rulePieces[:len(rulePieces)-1] {
			rulePieceParts := reRule.FindStringSubmatch(rulePiece)

			// 1 = xmas, 2 = < OR >, 3 = N, 4 = A/R/someOtherWorkflow
			conditions = append(conditions, Condition{
				Piece:    rulePieceParts[1],
				Comp:     rulePieceParts[2],
				Val:      ez.Atoi(rulePieceParts[3]),
				Terminal: rulePieceParts[4],
			})
		}

		rule := Rule{
			Conditions: conditions,
			Terminal:   ruleTerminal[0],
		}
		workflows[label] = rule
	}

	// Part 2
	if part2 {
		minVal := 1
		maxVal := 4000

		pr := PartRange{
			"x": Range{Min: minVal, Max: maxVal},
			"m": Range{Min: minVal, Max: maxVal},
			"a": Range{Min: minVal, Max: maxVal},
			"s": Range{Min: minVal, Max: maxVal},
		}

		return workflows.EvalRange(pr, "in")
	}

	// Parse parts
	parts := make([]Part, 0)
	for _, rawPart := range strings.Split(rulesAndParts[1], "\n") {
		partPieces := rePart.FindStringSubmatch(rawPart)
		ratings := make(map[string]int)
		ratings[partPieces[1]] = ez.Atoi(partPieces[2])
		ratings[partPieces[3]] = ez.Atoi(partPieces[4])
		ratings[partPieces[5]] = ez.Atoi(partPieces[6])
		ratings[partPieces[7]] = ez.Atoi(partPieces[8])
		parts = append(parts, Part{
			Rating: ratings,
			Total:  ez.Atoi(partPieces[2]) + ez.Atoi(partPieces[4]) + ez.Atoi(partPieces[6]) + ez.Atoi(partPieces[8]),
		})
	}

	sum := 0
	for i := range parts {
		if workflows.Eval(parts[i], "in") == "A" {
			sum += parts[i].Total
		}
	}
	return sum
}

// Eval is used in Part 1 to traverse the workflows and return A or R
func (w Workflow) Eval(part Part, label string) string {
	rule := w[label]
	for _, cond := range rule.Conditions {
		cond := cond
		passes := false
		switch cond.Comp {
		case ">":
			if part.Rating[cond.Piece] > cond.Val {
				passes = true
			}
		case "<":
			if part.Rating[cond.Piece] < cond.Val {
				passes = true
			}
		}

		if passes {
			switch cond.Terminal {
			case "A":
				return "A"
			case "R":
				return "R"
			default:
				return w.Eval(part, cond.Terminal)
			}
		}
	}

	// If no conditions are met, evaluate the terminal condition of the rule
	switch rule.Terminal {
	case "A":
		return "A"
	case "R":
		return "R"
	default:
		return w.Eval(part, rule.Terminal)
	}
}

type PartRange map[string]Range

type Range struct {
	Min int
	Max int
}

// EvalRange utilizes dynamic programming to evaluate a rang of possible solutions and ultimately return all possible combinations
func (w Workflow) EvalRange(pr PartRange, label string) int64 {
	rule := w[label]
	accepted := int64(0)
	for _, cond := range rule.Conditions {
		cond := cond
		// Make a copy of the part range
		// We'll be modifying this new copy and passing to recursive calls
		newPr := PartRange{
			"x": Range{Min: pr["x"].Min, Max: pr["x"].Max},
			"m": Range{Min: pr["m"].Min, Max: pr["m"].Max},
			"a": Range{Min: pr["a"].Min, Max: pr["a"].Max},
			"s": Range{Min: pr["s"].Min, Max: pr["s"].Max},
		}
		switch {
		case cond.Comp == ">" && pr[cond.Piece].Max > cond.Val:
			if newPr[string(cond.Piece)].Min < cond.Val {
				// Adjust the range to meet the condition
				// For greater than, we need to increase the min
				newRange := newPr[string(cond.Piece)]
				newRange.Min = cond.Val + 1
				newPr[string(cond.Piece)] = newRange

				// We need to adjust the original range to "fail" this condition
				// For greater than, we need to decrease the max
				newRange = pr[string(cond.Piece)]
				newRange.Max = cond.Val
				pr[string(cond.Piece)] = newRange
			}
		case cond.Comp == "<" && pr[cond.Piece].Min < cond.Val:
			if newPr[string(cond.Piece)].Max > cond.Val {
				// Adjust the range to meet the condition
				// For less than, we need to decrease the max
				newRange := newPr[string(cond.Piece)]
				newRange.Max = cond.Val - 1
				newPr[string(cond.Piece)] = newRange

				// We need to adjust the original range to "fail" this condition
				// For less than, we need to increase the min
				newRange = pr[string(cond.Piece)]
				newRange.Min = cond.Val
				pr[string(cond.Piece)] = newRange
			}
		}

		// Determine what to do next
		switch cond.Terminal {
		case "A":
			accepted += newPr.Within()
		case "R":
			// Do nothing
		default:
			accepted += w.EvalRange(newPr, cond.Terminal)
		}

	}

	// Also test the terminal on the "failing" conditions
	switch rule.Terminal {
	case "A":
		accepted += pr.Within()
	case "R":
		// Do Nothing
	default:
		accepted += w.EvalRange(pr, rule.Terminal)
	}

	return accepted
}

// Within evaluates each range and multiples each to produce a all possible solutions within the range
func (pr PartRange) Within() int64 {
	x := pr["x"].Max - pr["x"].Min + 1
	m := pr["m"].Max - pr["m"].Min + 1
	a := pr["a"].Max - pr["a"].Min + 1
	s := pr["s"].Max - pr["s"].Min + 1

	return int64(x * m * a * s)
}
