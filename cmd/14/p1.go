package main

import (
	"fmt"
	"os"

	utils "github.com/ghjm/advent_utils"
)

type data struct {
	grid [][]byte
}

func (d *data) copy() *data {
	nd := data{}
	for _, g := range d.grid {
		nd.grid = append(nd.grid, g[:])
	}
	return &nd
}

func (d *data) moveDirection(dy, dx int) bool {
	anyChanged := false
	for y := range d.grid {
		for x := range d.grid[y] {
			ny := y + dy
			nx := x + dx
			if ny < 0 || nx < 0 || ny >= len(d.grid) || nx >= len(d.grid[y]) {
				continue
			}
			if d.grid[y][x] == 'O' && d.grid[ny][nx] == '.' {
				d.grid[y][x] = '.'
				d.grid[ny][nx] = 'O'
				anyChanged = true
			}
		}
	}
	return anyChanged
}

func (d *data) moveAllNorth() {
	for d.moveDirection(-1, 0) {
	}
}

func (d *data) spin() {
	for d.moveDirection(-1, 0) {
	}
	for d.moveDirection(0, -1) {
	}
	for d.moveDirection(1, 0) {
	}
	for d.moveDirection(0, 1) {
	}
}

func (d *data) totalLoad() int {
	load := 0
	for y := range d.grid {
		for x := range d.grid[y] {
			if d.grid[y][x] == 'O' {
				load += len(d.grid) - y
			}
		}
	}
	return load
}

func (d *data) printout() {
	for y := range d.grid {
		fmt.Printf("%s\n", string(d.grid[y]))
	}
}

func (d *data) key() string {
	s := ""
	for _, g := range d.grid {
		s = s + string(g) + ","
	}
	return s
}

func run() error {
	d := data{}
	err := utils.OpenAndReadLines("input14.txt", func(s string) error {
		d.grid = append(d.grid, []byte(s))
		return nil
	})
	if err != nil {
		return err
	}
	p1 := d.copy()
	p1.moveAllNorth()
	fmt.Printf("Part 1: %d\n", p1.totalLoad())
	p2 := d.copy()
	positions := make(map[string]int)
	loads := make(map[int]int)
	var loopFrom, loopLen int
	for i := 0; i < 1000000000; i++ {
		p2.spin()
		loads[i] = p2.totalLoad()
		k := p2.key()
		p, ok := positions[k]
		if ok {
			loopFrom = p
			loopLen = i - p
			break
		}
		positions[k] = i
	}
	finalPos := loopFrom + ((1000000000 - loopFrom) % loopLen) - 1
	fmt.Printf("Part 2: %d\n", loads[finalPos])
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
