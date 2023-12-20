package main

import (
	"aoc-in-go/ez"
	"container/list"
	"regexp"
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

var reTypeLabel = regexp.MustCompile(`([%&])(\w+)`)

type Module struct {
	// Type represents the % (flip-flop) or & (conjuction), or broadcaster
	Type  string
	Label string
	// Input is the list of input module labels leading to this module
	Input []string
	// Output is the list of modules this module will output to
	Output []string

	// State is on/off, on = true, off = false
	// Only applies to Type = % (flip-flop) modules
	State bool

	// PulseMem is a memory store of the Input module (key: string)
	// Where true/false represents high/low memory for the pulse
	PulseMem map[string]bool
}

// ModList is a map, keyed by the label of the Module
type ModList map[string]*Module

// Pulse stores the from -> to as module labels, and whether the pulse is high (true) or low (false)
type Pulse struct {
	From string
	To   string
	High bool
}

func run(part2 bool, input string) any {
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	// Build the modList
	modList := make(ModList)
	for _, line := range lines {
		parts := strings.Split(line, " -> ")
		outputs := strings.Split(parts[1], ", ")
		if parts[0] == "broadcaster" {
			modList[parts[0]] = &Module{
				Type:   "broadcaster",
				Label:  "broadcaster",
				Input:  []string{},
				Output: outputs,
			}
		} else {
			typeAndLabel := reTypeLabel.FindStringSubmatch(parts[0])

			modList[typeAndLabel[2]] = &Module{
				Type:     typeAndLabel[1],
				Label:    typeAndLabel[2],
				Input:    []string{},
				Output:   outputs,
				State:    false,
				PulseMem: make(map[string]bool),
			}
		}
	}

	// Populate inputs and pulse memory for modList
	for label, module := range modList {
		for _, output := range module.Output {
			if inModule, ok := modList[output]; ok {
				inModule.Input = append(inModule.Input, label)
				inModule.PulseMem[label] = false
				modList[output] = inModule
			}
		}
	}

	// Part 2
	if part2 {
		if len(lines) < 20 {
			// Example is not supported
			return 1
		}

		return modList.PushButton2()
	}

	// Part 1
	high := 0
	low := 0
	for i := 1; i <= 1000; i++ {
		addHigh, addLow := modList.PushButton()
		high += addHigh
		low += addLow
	}

	return high * low
}

// PushButton is for Part 1, returning how many high/low signals are sent when the button is pushed
func (l ModList) PushButton() (int, int) {
	high := 0
	low := 0

	q := list.New()

	// Broadcast
	q.PushBack(Pulse{
		From: "",
		To:   "broadcaster",
		High: false,
	})
	for q.Len() > 0 {
		e := q.Front()
		p := e.Value.(Pulse)

		// Increment high/low
		if p.High {
			high++
		} else {
			low++
		}

		// Handle pulse
		if _, exists := l[p.To]; exists {
			l[p.To].HandlePulse(q, p)
		}

		q.Remove(e) // Dequeue
	}

	return high, low
}

// HandlePusle determines what to do given a pulse is "sent" to the module
func (m *Module) HandlePulse(q *list.List, p Pulse) {
	switch m.Type {
	case "broadcaster":
		for _, out := range m.Output {
			q.PushBack(Pulse{
				From: m.Label,
				To:   out,
				High: p.High,
			})
		}
	case "%":
		// Flip-flop
		if p.High {
			return
		}

		if m.State == true {
			// On
			// Send low pulse
			for _, out := range m.Output {
				q.PushBack(Pulse{
					From: m.Label,
					To:   out,
					High: false,
				})
			}

			// Flip to off
			m.State = false
		} else {
			// Off
			// Send high pulse
			for _, out := range m.Output {
				q.PushBack(Pulse{
					From: m.Label,
					To:   out,
					High: true,
				})
			}

			// Flip to on
			m.State = true
		}
	case "&":
		// Conjuction
		// When a pulse is received, the conjunction module first updates its memory for that input.
		m.PulseMem[p.From] = p.High

		// Then, if it remembers high pulses for all inputs
		allHigh := true
		for _, mem := range m.PulseMem {
			if mem == false {
				allHigh = false
				break
			}
		}
		if allHigh {
			// it sends a low pulse;
			for _, out := range m.Output {
				q.PushBack(Pulse{
					From: m.Label,
					To:   out,
					High: false,
				})
			}
		} else {
			// otherwise, it sends a high pulse.
			for _, out := range m.Output {
				q.PushBack(Pulse{
					From: m.Label,
					To:   out,
					High: true,
				})
			}
		}
	}
}

// PushButton2 is for Part 2, and returns the LCM of cycles needed for rx to receive a single low signal
func (l ModList) PushButton2() int {
	// NOTE: We assume the puzzle is a cycle
	//
	// rx is bound to a single Conjunction (&) module
	// The input for rx then is bound to ~4 other Conjunction (&) modules
	// Detect those parent Conjuction modules
	// When they are "low", store the amount of itterations needed to become low
	// Use LCM once we detect the low cycle for all to be low
	parents := map[string]int{}
	rxParent := ""
	for label, module := range l {
		if module.Output[0] == "rx" {
			rxParent = label
		}
	}
	for label, module := range l {
		if lo.Contains(module.Output, rxParent) {
			parents[label] = 0
		}
	}

	i := 0
	for {
		i++

		q := list.New()

		// Broadcast
		q.PushBack(Pulse{
			From: "",
			To:   "broadcaster",
			High: false,
		})
		for q.Len() > 0 {
			e := q.Front()
			p := e.Value.(Pulse)

			// Capture when a low signal is sent to one of our parents
			if _, exists := parents[p.To]; exists && p.High == false {
				parents[p.To] = i
			}

			// Handle pulse
			if _, exists := l[p.To]; exists {
				l[p.To].HandlePulse(q, p)
			}

			q.Remove(e) // Dequeue
		}

		// Detect if we have any remaining parents that need a cycle set
		if !lo.Contains(lo.Values(parents), 0) {
			break
		}
	}

	vals := lo.Values(parents)
	return ez.LCM(vals[0], vals[1], vals[2:]...)
}
