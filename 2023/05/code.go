package main

import (
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
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	var seedToSoil, soilToFert, fertToWater, waterToLight, lightToTemp, tempToHum, humToLoc []Boundary
	if len(lines) < 50 {
		// Example file
		seedToSoil = Boundaries(lines[3:5])
		soilToFert = Boundaries(lines[7:10])
		fertToWater = Boundaries(lines[12:16])
		waterToLight = Boundaries(lines[18:20])
		lightToTemp = Boundaries(lines[22:25])
		tempToHum = Boundaries(lines[27:29])
		humToLoc = Boundaries(lines[31:33])
	} else {
		seedToSoil = Boundaries(lines[3:15])
		soilToFert = Boundaries(lines[17:46])
		fertToWater = Boundaries(lines[48:65])
		waterToLight = Boundaries(lines[67:85])
		lightToTemp = Boundaries(lines[87:121])
		tempToHum = Boundaries(lines[123:158])
		humToLoc = Boundaries(lines[160:172])
	}

	var seeds []string
	lowestLoc := int64(9999999999999999)
	if part2 {
		seedDefs := strings.Split(lines[0][7:], " ")
		for i := 0; i < len(seedDefs); i += 2 {
			seedStart := Atoi(seedDefs[i])
			for j := int64(0); j <= Atoi(seedDefs[i+1]); j++ {
				seed := seedStart + j
				soil, soilSkip := Next(seed, seedToSoil)
				fert, fertSkip := Next(soil, soilToFert)
				water, waterSkip := Next(fert, fertToWater)
				light, lightSkip := Next(water, waterToLight)
				temp, tempSkip := Next(light, lightToTemp)
				hum, humSkip := Next(temp, tempToHum)
				loc, locSkip := Next(hum, humToLoc)
				if loc < lowestLoc {
					lowestLoc = loc
				}
				// Skip allows us to "jump" ahead an amount of seeds where incrementing by 1 would just produce a location that is 1 more
				// When j++ is evaluated, the location should be significantly different as the increment will change one or more of the Boundaries that are matched in Next
				j += j + max(min(soilSkip, fertSkip, waterSkip, lightSkip, tempSkip, humSkip, locSkip)-1, 0)
			}
		}
	} else {
		seeds = strings.Split(lines[0][7:], " ")
		for _, seedStr := range seeds {
			seed := Atoi(seedStr)
			soil, _ := Next(seed, seedToSoil)
			fert, _ := Next(soil, soilToFert)
			water, _ := Next(fert, fertToWater)
			light, _ := Next(water, waterToLight)
			temp, _ := Next(light, lightToTemp)
			hum, _ := Next(temp, tempToHum)
			loc, _ := Next(hum, humToLoc)
			if loc < lowestLoc {
				lowestLoc = loc
			}
		}
	}

	return lowestLoc
}

type Boundary struct {
	LeftMin int64
	LeftMax int64
	Right   int64
}

func Boundaries(lines []string) []Boundary {
	out := []Boundary{}
	for _, line := range lines {
		parts := strings.Split(line, " ")
		// Note: parts are destination/rgt THEN source/lft, and finally the distance/range
		rgt, lft, dis := parts[0], parts[1], parts[2]
		out = append(out, Boundary{
			LeftMin: Atoi(lft),
			LeftMax: Atoi(lft) + Atoi(dis),
			Right:   Atoi(rgt),
		})
	}
	return out
}

// Next provides the location of the destination/rgt given an in/source/rgt
// Additionally, a "skip" is returned, which is the remaining range within the boundary
func Next(in int64, boundaries []Boundary) (int64, int64) {
	for _, b := range boundaries {
		// If within the boundaries of the source/lft
		if in >= b.LeftMin && in < b.LeftMax {
			// Compute the difference, so we can add it to the destination/rgt
			diff := in - b.LeftMin
			// Skip is the remaining "range" that is just incrementally increasing by 1
			skip := b.LeftMax - in
			return b.Right + diff, skip
		}
	}

	// If not found, the source and destination are the same
	return in, 0
}

// Atoi ignores the errors in strconv.Atoi and returns the response
func Atoi(in string) int64 {
	out, _ := strconv.Atoi(in)
	return int64(out)
}
