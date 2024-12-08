package main

import (
	"fmt"
	"os"
	"strings"

	utils "github.com/ghjm/advent_utils"
)

type data struct {
	lines    []string
	startPos point
}

type point struct {
	x, y int
}

type state struct {
	pos   point
	steps int
}

func (d *data) legalMovesFrom(p point) []point {
	var results []point
	for _, dir := range []point{{-1, 0}, {1, 0}, {0, 1}, {0, -1}} {
		np := point{p.x + dir.x, p.y + dir.y}
		if np.x >= 0 && np.y >= 0 && np.x < len(d.lines[0]) && np.y < len(d.lines) &&
			(d.lines[np.y][np.x] == '.' || d.lines[np.y][np.x] == 'S') {
			results = append(results, np)
		}
	}
	return results
}

func run() error {
	d := data{}
	err := utils.OpenAndReadLines("input21.txt", func(s string) error {
		d.lines = append(d.lines, s)
		idx := strings.Index(s, "S")
		if idx >= 0 {
			d.startPos = point{idx, len(d.lines) - 1}
		}
		return nil
	})
	if err != nil {
		return err
	}
	open := []state{state{d.startPos, 0}}
	visited := make(map[state]struct{})
	solutions := make(map[point]struct{})
	lmf := make(map[point][]point)
	max := 0
	for len(open) > 0 {
		s := open[0]
		open = open[1:]
		if s.steps > max {
			max = s.steps
			fmt.Printf("New max: %d, len(open)=%d, len(visited)=%d\n", max, len(open), len(visited))
		}
		if s.steps == 64 {
			solutions[s.pos] = struct{}{}
			continue
		}
		nps, ok := lmf[s.pos]
		if !ok {
			nps = d.legalMovesFrom(s.pos)
			lmf[s.pos] = nps
		}
		for _, np := range nps {
			ns := state{np, s.steps + 1}
			_, ok := visited[ns]
			if !ok {
				open = append(open, ns)
			}
			visited[ns] = struct{}{}
		}
	}
	fmt.Printf("%d\n", len(solutions))
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
