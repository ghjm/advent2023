package main

import (
	"fmt"
	"github.com/ghjm/advent2023/pkg/utils"
	"os"
	"regexp"
	"strings"
)

type node struct {
	left  string
	right string
}

type state struct {
	pos     string
	instPos int64
}

type data struct {
	instructions string
	nodes        map[string]node
}

func run() error {
	d := data{
		instructions: "",
		nodes:        make(map[string]node),
	}
	nodeRe := regexp.MustCompile(`(...) = \((...), (...)\)`)
	err := utils.OpenAndReadLines("input8.txt", func(s string) error {
		if s == "" {
			return nil
		}
		if d.instructions == "" {
			d.instructions = s
			return nil
		}
		m := nodeRe.FindStringSubmatch(s)
		if m == nil {
			return fmt.Errorf("bad input")
		}
		d.nodes[m[1]] = node{
			left:  m[2],
			right: m[3],
		}
		return nil
	})
	if err != nil {
		return err
	}
	var i int64
	pos := "AAA"
	for {
		instr := d.instructions[i%int64(len(d.instructions))]
		switch instr {
		case 'L':
			pos = d.nodes[pos].left
		case 'R':
			pos = d.nodes[pos].right
		default:
			return fmt.Errorf("invalid instruction")
		}
		i += 1
		if pos == "ZZZ" {
			break
		}
	}
	fmt.Printf("Part 1: %d\n", i)
	i = 0
	posns := make([]string, 0)
	for n := range d.nodes {
		if strings.HasSuffix(n, "A") {
			posns = append(posns, n)
		}
	}
	pastPosns := make([]map[state]int64, len(posns))
	loopLens := make([]int64, len(posns))
	lastZ := make([]int64, len(posns))
	for j := range pastPosns {
		pastPosns[j] = make(map[state]int64)
		loopLens[j] = -1
		lastZ[j] = -1
	}
	for {
		instPos := i % int64(len(d.instructions))
		instr := d.instructions[instPos]
		for j := range posns {
			switch instr {
			case 'L':
				posns[j] = d.nodes[posns[j]].left
			case 'R':
				posns[j] = d.nodes[posns[j]].right
			default:
				return fmt.Errorf("invalid instruction")
			}
			if strings.HasSuffix(posns[j], "Z") {
				lastZ[j] = i
			}
			st := state{
				pos:     posns[j],
				instPos: instPos,
			}
			pst, ok := pastPosns[j][st]
			if ok {
				if loopLens[j] != -1 && loopLens[j] != i-pst {
					fmt.Printf("loop length changed on thread %d: was %d now is %d\n", j, loopLens[j], i-pst)
				}
				loopLens[j] = i - pst
			}
			pastPosns[j][st] = i
		}
		allZ := true
		for j := range posns {
			if !strings.HasSuffix(posns[j], "Z") {
				allZ = false
			}
		}
		if allZ {
			fmt.Printf("Part 2: %d\n", i+1)
			return nil
		}
		allLooped := true
		for j := range posns {
			if loopLens[j] == -1 {
				allLooped = false
			}
		}
		if allLooped {
			break
		}
		i += 1
	}
	fmt.Printf("Part 2: %d\n", utils.LCM(loopLens...)) // This was correct for my input, but seems sketchy
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
