package main

import (
	"strings"
	"sync"

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

type Cell struct {
	Char      string
	Energized int
	Tracked   []string
	tMu       sync.RWMutex
}

func (c *Cell) GetTracked() []string {
	c.tMu.RLock()
	defer c.tMu.RUnlock()
	return c.Tracked
}

func (c *Cell) AddTracked(s string) {
	c.tMu.Lock()
	defer c.tMu.Unlock()
	c.Tracked = append(c.Tracked, s)
}

// var cells [][]Cell
var rowCount int
var colCount int

func run(part2 bool, input string) any {
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	rowCount = len(lines)
	colCount = len(lines[0])

	// Part 2
	if part2 {
		sumMax := 0
		// Left and Rights
		for i := 0; i <= rowCount; i++ {
			cellCopyA := MakeCells(lines)
			wgA := &sync.WaitGroup{}
			wgA.Add(1)
			go HandleBeam(wgA, cellCopyA, i, 0, "E")

			cellCopyB := MakeCells(lines)
			wgB := &sync.WaitGroup{}
			wgB.Add(1)
			go HandleBeam(wgB, cellCopyB, i, colCount-1, "W")

			wgA.Wait()
			wgB.Wait()

			sumMax = max(sumMax, SumCells(cellCopyA), SumCells(cellCopyB))
		}

		// Tops and Bottoms
		for i := 0; i <= colCount; i++ {
			cellCopyA := MakeCells(lines)
			wgA := &sync.WaitGroup{}
			wgA.Add(1)
			go HandleBeam(wgA, cellCopyA, 0, i, "S")

			cellCopyB := MakeCells(lines)
			wgB := &sync.WaitGroup{}
			wgB.Add(1)
			go HandleBeam(wgB, cellCopyB, rowCount-1, i, "N")

			wgA.Wait()
			wgB.Wait()

			sumMax = max(sumMax, SumCells(cellCopyA), SumCells(cellCopyB))
		}

		return sumMax
	}

	// Part 1
	cellCopy := MakeCells(lines)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go HandleBeam(wg, cellCopy, 0, 0, "E")
	wg.Wait()

	return SumCells(cellCopy)
}

// Make cells generates a fresh, 0 initialized matrix of Cells
func MakeCells(lines []string) [][]Cell {
	cells := make([][]Cell, len(lines))
	for i, line := range lines {
		cells[i] = make([]Cell, len(line))

		for j, char := range strings.Split(line, "") {
			cells[i][j] = Cell{
				Char:      char,
				Energized: 0,
				Tracked:   []string{},
			}
		}
	}
	return cells
}

// SumCells returns the count of Energized cells
func SumCells(cells [][]Cell) int {
	sum := 0
	for i := range cells {
		for j := range cells[i] {
			if cells[i][j].Energized > 0 {
				sum++
			}
		}
	}
	return sum
}

// HandleBeam terminates, moves, redirects, or splits a beam given a point on the grid
func HandleBeam(wg *sync.WaitGroup, cells [][]Cell, i, j int, dir string) {
	// Exit if past the edge of the grid
	if i < 0 || i >= rowCount || j < 0 || j >= colCount {
		wg.Done()
		return
	}

	// Exit if the path has entered the cell going in the same direction before
	if lo.Contains(cells[i][j].GetTracked(), dir) {
		wg.Done()
		return
	}

	// Mark cell energized
	cells[i][j].Energized++
	// Track the direction the path is going for this cell to assist with exiting if the beam gets in a loop
	cells[i][j].AddTracked(dir)

	// Handle what to do
	switch cells[i][j].Char {
	case ".":
		// Continue in current direction
		switch dir {
		case "N":
			HandleBeam(wg, cells, i-1, j, "N")
		case "E":
			HandleBeam(wg, cells, i, j+1, "E")
		case "S":
			HandleBeam(wg, cells, i+1, j, "S")
		case "W":
			HandleBeam(wg, cells, i, j-1, "W")
		}
	case "/":
		// Change Directions SW, NE
		switch dir {
		case "N":
			HandleBeam(wg, cells, i, j+1, "E")
		case "E":
			HandleBeam(wg, cells, i-1, j, "N")
		case "S":
			HandleBeam(wg, cells, i, j-1, "W")
		case "W":
			HandleBeam(wg, cells, i+1, j, "S")
		}
	case "\\":
		// Change Directions SE, NW
		switch dir {
		case "N":
			HandleBeam(wg, cells, i, j-1, "W")
		case "E":
			HandleBeam(wg, cells, i+1, j, "S")
		case "S":
			HandleBeam(wg, cells, i, j+1, "E")
		case "W":
			HandleBeam(wg, cells, i-1, j, "N")
		}
	case "|":
		// Split (EW) or keep going (NS)
		switch dir {
		case "N":
			// Continue
			HandleBeam(wg, cells, i-1, j, "N")
		case "E":
			// Split the current beam into two new beams
			wg.Add(1)
			go HandleBeam(wg, cells, i-1, j, "N")
			wg.Add(1)
			go HandleBeam(wg, cells, i+1, j, "S")

			wg.Done()
		case "S":
			// Continue
			HandleBeam(wg, cells, i+1, j, "S")
		case "W":
			// Split the current beam into two new beams
			wg.Add(1)
			go HandleBeam(wg, cells, i-1, j, "N")
			wg.Add(1)
			go HandleBeam(wg, cells, i+1, j, "S")

			wg.Done()
		}
	case "-":
		// Split (NS) or keep going (EW)
		switch dir {
		case "N":
			// Split the current beam into two new beams
			wg.Add(1)
			go HandleBeam(wg, cells, i, j-1, "W")
			wg.Add(1)
			go HandleBeam(wg, cells, i, j+1, "E")

			wg.Done()
		case "E":
			// Continue
			HandleBeam(wg, cells, i, j+1, "E")
		case "S":
			// Split the current beam into two new beams
			wg.Add(1)
			go HandleBeam(wg, cells, i, j-1, "W")
			wg.Add(1)
			go HandleBeam(wg, cells, i, j+1, "E")

			wg.Done()
		case "W":
			// Continue
			HandleBeam(wg, cells, i, j-1, "W")
		}
	}
}
