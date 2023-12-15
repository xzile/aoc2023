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
	lineChars := make([][]string, len(lines))
	for i, line := range lines {
		lineChars[i] = strings.Split(line, "")
	}
	sum := 0

	if part2 {
		// Part 2
		loadFoundInARow := 0
		loadCheck := map[int]int{}
		loadPerCycle := []int{}
		preVal := 0
		for iter := 1; iter <= 1000000000; iter++ {
			// Roll "up"
			Roll(lineChars)
			// West is now "up", North is "right"
			lineChars = TransposeRight(lineChars)
			// Roll "up"
			Roll(lineChars)
			// South is now "up", North is "down"
			lineChars = TransposeRight(lineChars)
			// Roll "up"
			Roll(lineChars)
			// East is now "up", North is "left"
			lineChars = TransposeRight(lineChars)
			// Roll "up"
			Roll(lineChars)
			// North is now "up"
			lineChars = TransposeRight(lineChars)

			loadVal := LoadNoMove(lineChars)
			loadPerCycle = append(loadPerCycle, loadVal)
			if loadCheck[loadVal] > 0 && loadCheck[preVal] == iter-1 {
				loadFoundInARow++
				if loadFoundInARow >= 1000 {
					// Cycle found
					endOfCycle := iter - 1000
					startOfCycle := loadCheck[loadVal] - 1000
					cycleLen := endOfCycle - startOfCycle
					desiredCycle := (1000000000-startOfCycle)%cycleLen - 1
					sum = loadPerCycle[startOfCycle+desiredCycle]
					break
				}
			} else {
				loadFoundInARow = 0
			}
			loadCheck[loadVal] = iter
			preVal = loadVal

		}
	} else {
		// Part 1
		sum = Load(lineChars)
	}

	return sum
}

func Load(lines [][]string) int {
	colStartLoad := lo.RepeatBy(len(lines[0]), func(index int) int {
		return len(lines)
	})
	colLoad := make([]int, len(lines[0]))
	for i, line := range lines {
		for j, colChar := range line {
			switch {
			case colChar == "#":
				colStartLoad[j] = len(lines) - i - 1
			case colChar == "O":
				colLoad[j] += colStartLoad[j]
				colStartLoad[j]--
			}
		}
	}
	return ez.Sum(colLoad)
}

func LoadNoMove(lines [][]string) int {
	sum := 0
	for i, line := range lines {
		for _, colChar := range line {
			switch {
			case colChar == "O":
				sum += len(lines) - i
			}
		}
	}
	return sum
}

func Roll(lines [][]string) {
	colFallTo := lo.RepeatBy(len(lines[0]), func(index int) int {
		return 0
	})
	for i, line := range lines {
		for j, colChar := range line {
			switch {
			case colChar == "#":
				colFallTo[j] = i + 1
			case colChar == "O":
				lines[colFallTo[j]][j] = "O"
				if colFallTo[j] != i {
					lines[i][j] = "."
				}

				colFallTo[j]++
			}
		}
	}
}

func TransposeRight(slice [][]string) [][]string {
	colCount := len(slice[0])
	rowCount := len(slice)
	result := make([][]string, colCount)
	for i := range result {
		result[i] = make([]string, rowCount)
	}
	for col := 0; col < colCount; col++ {
		for row := 0; row < rowCount; row++ {
			result[col][colCount-row-1] = slice[row][col]
		}
	}
	return result
}

func TransposeLeft(slice [][]string) [][]string {
	colCount := len(slice[0])
	rowCount := len(slice)
	result := make([][]string, colCount)
	for i := range result {
		result[i] = make([]string, rowCount)
	}
	for col := 0; col < colCount; col++ {
		for row := 0; row < rowCount; row++ {
			result[rowCount-row-1][col] = slice[row][col]
		}
	}
	return result
}
