package main

import (
	"fmt"
	"github.com/ghjm/advent2023/pkg/utils"
	"os"
)

type point struct {
	x, y int
}

type state struct {
	pos point
	dir byte
}

type data struct {
	lines []string
}

var directions = map[byte]point{
	'U': {0, -1},
	'D': {0, 1},
	'L': {-1, 0},
	'R': {1, 0},
}

var movements = map[byte]map[byte][]byte{
	'U': {
		'.':  {'U'},
		'/':  {'R'},
		'\\': {'L'},
		'|':  {'U'},
		'-':  {'L', 'R'},
	},
	'D': {
		'.':  {'D'},
		'/':  {'L'},
		'\\': {'R'},
		'|':  {'D'},
		'-':  {'L', 'R'},
	},
	'L': {
		'.':  {'L'},
		'/':  {'D'},
		'\\': {'U'},
		'|':  {'U', 'D'},
		'-':  {'L'},
	},
	'R': {
		'.':  {'R'},
		'/':  {'U'},
		'\\': {'D'},
		'|':  {'U', 'D'},
		'-':  {'R'},
	},
}

func (d *data) numEnergized(start state) int {
	open := []state{start}
	visited := make(map[state]struct{})
	for len(open) > 0 {
		s := open[0]
		open = open[1:]
		visited[s] = struct{}{}
		for _, m := range movements[s.dir][d.lines[s.pos.y][s.pos.x]] {
			dp := directions[m]
			ns := state{
				pos: point{s.pos.x + dp.x, s.pos.y + dp.y},
				dir: m,
			}
			if ns.pos.x >= 0 && ns.pos.y >= 0 && ns.pos.x < len(d.lines[0]) && ns.pos.y < len(d.lines) {
				_, ok := visited[ns]
				if !ok {
					open = append(open, ns)
				}
			}
		}
	}
	mp := make(map[point]struct{})
	for s := range visited {
		mp[s.pos] = struct{}{}
	}
	return len(mp)
}

func run() error {
	d := data{}
	err := utils.OpenAndReadLines("input16.txt", func(s string) error {
		d.lines = append(d.lines, s)
		return nil
	})
	if err != nil {
		return err
	}
	fmt.Printf("Part 1: %d\n", d.numEnergized(state{point{0, 0}, 'R'}))
	maxE := 0
	for y := range d.lines {
		e := d.numEnergized(state{point{0, y}, 'R'})
		maxE = utils.Max(maxE, e)
		e = d.numEnergized(state{point{len(d.lines[0]) - 1, y}, 'L'})
		maxE = utils.Max(maxE, e)
	}
	for x := range d.lines[0] {
		e := d.numEnergized(state{point{x, 0}, 'D'})
		maxE = utils.Max(maxE, e)
		e = d.numEnergized(state{point{x, len(d.lines) - 1}, 'U'})
		maxE = utils.Max(maxE, e)
	}
	fmt.Printf("Part 2: %d\n", maxE)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
