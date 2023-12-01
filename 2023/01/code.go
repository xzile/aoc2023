package main

import (
	"bufio"
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
	sum := 0
	sc := bufio.NewScanner(strings.NewReader(input))
	for sc.Scan() {
		line := sc.Text()
		newLine := line
		if len(line) == 0 {
			continue
		}
		if part2 {
			numWordsRx := regexp.MustCompile(`(one|two|three|four|five|six|seven|eight|nine)`)
			words := map[string]string{
				// Account for word-reuse, e.g. eightwo where the end result should be 82
				"one":   "o1e",
				"two":   "t2o",
				"three": "t3e",
				"four":  "4",
				"five":  "5e",
				"six":   "6",
				"seven": "7n",
				"eight": "e8t",
				"nine":  "n9e",
			}
			for {
				foundWord := numWordsRx.FindString(newLine)
				newLine = strings.Replace(newLine, foundWord, words[foundWord], 1)
				if foundWord == "" {
					break
				}
			}
		}

		nonDigits := regexp.MustCompile(`\D`)
		onlyNums := nonDigits.ReplaceAllString(newLine, "")

		firstNum := onlyNums[0:1]
		lastNum := onlyNums[len(onlyNums)-1:]
		newNum, _ := strconv.Atoi(fmt.Sprintf("%s%s", firstNum, lastNum))
		//fmt.Println(fmt.Sprintf("line: %s; newLine: %s; onlyNums: %s; newNum: %d", line, newLine, onlyNums, newNum))

		sum += newNum
	}

	return sum
}
