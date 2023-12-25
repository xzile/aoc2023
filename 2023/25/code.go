package main

import (
	"container/list"
	"math/rand"
	"slices"
	"strings"

	dij "github.com/RyanCarrier/dijkstra"
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

type Comp struct {
	ID    int
	Label string
}

type Edge struct {
	Max int
	Min int
}

func run(part2 bool, input string) any {
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")

	// Just parsing to various maps and a graph for dijkstra
	g := dij.NewGraph()
	ids := make(map[string]int)
	comps := make(map[int]Comp)
	conns := make(map[Comp]map[Comp]bool)
	id := 1
	for _, line := range lines {
		parts := strings.Split(line, ": ")
		a := Comp{
			Label: parts[0],
		}
		if anId, ok := ids[a.Label]; !ok {
			a.ID = id
			g.AddVertex(id)
			id++
		} else {
			a.ID = anId
		}
		ids[a.Label] = a.ID
		comps[a.ID] = a
		if _, ok := conns[a]; !ok {
			conns[a] = make(map[Comp]bool)
		}

		for _, con := range strings.Split(parts[1], " ") {
			con := con
			b := Comp{
				Label: con,
			}
			if anId, ok := ids[b.Label]; !ok {
				b.ID = id
				g.AddVertex(id)
				id++
			} else {
				b.ID = anId
			}
			ids[b.Label] = b.ID
			comps[b.ID] = b
			if _, ok := conns[b]; !ok {
				conns[b] = make(map[Comp]bool)
			}

			// Connections, both ways
			conns[a][b] = true
			conns[b][a] = true
			g.AddArc(a.ID, b.ID, 1)
			g.AddArc(b.ID, a.ID, 1)
		}
	}

	// No part 2 for day 25. Merry Christmas!
	if part2 {
		return "not implemented"
	}

	// Do random checks and record the edges that are crossed between the two points
	// The most common 3 edges are what we're after, these are most likely going to be the three to cut
	seen := make(map[Edge]int)
	// n could be lower, but, 20,000 gives a very high confidence
	// I was getting the same results at 1,000, but given the total node count > 1000, a few extra itterations isn't too bad
	for n := 0; n < 20000; n++ {
		aID := rand.Intn(len(ids)-2) + 1
		bID := rand.Intn(len(ids)-2) + 1
		for bID == aID {
			// Just making sure a != b
			bID = rand.Intn(len(ids)-2) + 1
		}
		short := lo.Must(g.Shortest(aID, bID))
		// Store the edge
		// Using min/max to make sure the direction of the path doesn't skew the results
		for i := 0; i < len(short.Path)-1; i++ {
			seen[Edge{
				Max: max(short.Path[i], short.Path[i+1]),
				Min: min(short.Path[i], short.Path[i+1]),
			}]++
		}
	}

	// Quick map[K]V -> []{K, V} so we can sort
	entries := lo.Entries(seen)
	slices.SortFunc(entries, func(a, b lo.Entry[Edge, int]) int {
		return b.Value - a.Value
	})

	// Take the top 3
	cuts := entries[:3]
	//Break the connection of the offending nodes, in both directions
	for _, cut := range cuts {
		a := comps[cut.Key.Max]
		b := comps[cut.Key.Min]
		conns[a][b] = false
		conns[b][a] = false
	}
	// Get the connected nodes within each set
	a := ConnectedNodes(conns, comps[cuts[0].Key.Max])
	b := ConnectedNodes(conns, comps[cuts[0].Key.Min])

	// solve
	return len(a) * len(b)
}

// ConnectedNodes uses the connection map and a given start to determine which components are connected
func ConnectedNodes(conns map[Comp]map[Comp]bool, start Comp) []Comp {
	seen := []Comp{}

	q := list.New()
	q.PushBack(start)

	for q.Len() > 0 {
		aAny := q.Front()
		a := aAny.Value.(Comp)
		seen = append(seen, a)

		// Get the connected nodes, and make sure the connection hasn't been disabled (t == false)
		for b, t := range conns[a] {
			if t && !lo.Contains(seen, b) {
				q.PushBack(b)
			}
		}

		q.Remove(aAny)
	}

	return lo.Uniq(seen)
}
