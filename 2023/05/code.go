package main

import (
	"fmt"
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
	var seedToSoil, soilToFert, fertToWater, waterToLight, lightToTemp, tempToHum, humToLoc []Boundry
	if len(lines) < 50 {
		// Example file
		seedToSoil = Boundries(lines[3:5])
		soilToFert = Boundries(lines[7:10])
		fertToWater = Boundries(lines[12:16])
		waterToLight = Boundries(lines[18:20])
		lightToTemp = Boundries(lines[22:25])
		tempToHum = Boundries(lines[27:29])
		humToLoc = Boundries(lines[31:33])
	} else {
		seedToSoil = Boundries(lines[3:15])
		soilToFert = Boundries(lines[17:46])
		fertToWater = Boundries(lines[48:65])
		waterToLight = Boundries(lines[67:85])
		lightToTemp = Boundries(lines[87:121])
		tempToHum = Boundries(lines[123:158])
		humToLoc = Boundries(lines[160:172])
	}

	var seeds []string
	lowestLoc := int64(9999999999999999)
	if part2 {
		Log("Using brute force... please wait.")
		seedDefs := strings.Split(lines[0][7:], " ")
		for i := 0; i < len(seedDefs); i += 2 {
			// NOTE: There's probably a much more efficient data structure that computes the breakpoints of the entire "range" of seeds,
			// where the next seed ultimately changes the location by something non-incrementally.
			// This would allow you to only check the breakpoints.
			// To set this up, though, you'd likely need to start at location and work backwards, identifying in the previous step:
			// - it's own breakpoints
			// - where in those ranges location changes non-incrementally
			//
			// However, brute force only took ~1min 30sec to generate the correct answer, which is faster than I can come up the above.

			seedStart := Atoi(seedDefs[i])
			// Brute force, checking each seed
			for j := int64(0); j <= Atoi(seedDefs[i+1]); j++ {
				seed := seedStart + j
				soil := Next(seed, seedToSoil)
				fert := Next(soil, soilToFert)
				water := Next(fert, fertToWater)
				light := Next(water, waterToLight)
				temp := Next(light, lightToTemp)
				hum := Next(temp, tempToHum)
				loc := Next(hum, humToLoc)
				if loc < lowestLoc {
					lowestLoc = loc
				}
			}
		}
	} else {
		seeds = strings.Split(lines[0][7:], " ")
		for _, seedStr := range seeds {
			seed := Atoi(seedStr)
			soil := Next(seed, seedToSoil)
			fert := Next(soil, soilToFert)
			water := Next(fert, fertToWater)
			light := Next(water, waterToLight)
			temp := Next(light, lightToTemp)
			hum := Next(temp, tempToHum)
			loc := Next(hum, humToLoc)
			if loc < lowestLoc {
				lowestLoc = loc
			}
		}
	}

	return lowestLoc
}

type Boundry struct {
	LeftMin int64
	LeftMax int64
	Right   int64
}

func Boundries(lines []string) []Boundry {
	out := []Boundry{}
	for _, line := range lines {
		parts := strings.Split(line, " ")
		// Note: parts are destination/rgt THEN source/lft, and finally the distance/range
		rgt, lft, dis := parts[0], parts[1], parts[2]
		out = append(out, Boundry{
			LeftMin: Atoi(lft),
			LeftMax: Atoi(lft) + Atoi(dis),
			Right:   Atoi(rgt),
		})
	}
	return out
}

func Next(in int64, boundries []Boundry) int64 {
	for _, b := range boundries {
		// If within the boundries of the source/lft
		if in >= b.LeftMin && in < b.LeftMax {
			// Compute the difference, so we can add it to the destingation/rgt
			diff := in - b.LeftMin
			return b.Right + diff
		}
	}

	// If not found, the source and destination are the same
	return in
}

// Atoi ignores the errors in strconv.Atoi and returns the response
func Atoi(in string) int64 {
	out, _ := strconv.Atoi(in)
	return int64(out)
}

// Log is a simple wrapper around Println + Sprintf
func Log(format string, a ...any) {
	fmt.Println(fmt.Sprintf(format, a...))
}
