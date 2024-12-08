package main

import (
	"fmt"
	"os"

	utils "github.com/ghjm/advent_utils"
)

type point struct {
	x, y int
}

type data struct {
	cells      [][]byte
	galaxies   []point
	expandRows []int
	expandCols []int
}

func (d *data) findGalaxies() {
	d.galaxies = make([]point, 0)
	for y := range d.cells {
		for x := range d.cells[y] {
			if d.cells[y][x] == '#' {
				d.galaxies = append(d.galaxies, point{x, y})
			}
		}
	}
}

func (d *data) findExpansions() {
	for y := range d.cells {
		empty := true
		for x := range d.cells[y] {
			if d.cells[y][x] != '.' {
				empty = false
				break
			}
		}
		if empty {
			d.expandRows = append(d.expandRows, y)
		}
	}
	for x := range d.cells[0] {
		empty := true
		for y := range d.cells {
			if d.cells[y][x] != '.' {
				empty = false
				break
			}
		}
		if empty {
			d.expandCols = append(d.expandCols, x)
		}
	}
}

func (d *data) expandGalaxies(expansion int) []point {
	result := make([]point, 0)
	for _, p := range d.galaxies {
		np := p
		for _, er := range d.expandRows {
			if p.y > er {
				np.y += expansion - 1
			}
		}
		for _, ec := range d.expandCols {
			if p.x > ec {
				np.x += expansion - 1
			}
		}
		result = append(result, np)
	}
	return result
}

func run() error {
	d := data{
		cells:      make([][]byte, 0),
		expandRows: make([]int, 0),
		expandCols: make([]int, 0),
	}
	err := utils.OpenAndReadLines("input11.txt", func(s string) error {
		d.cells = append(d.cells, []byte(s))
		return nil
	})
	if err != nil {
		return err
	}
	d.findGalaxies()
	d.findExpansions()
	for i, exp := range []int{2, 1000000} {
		sum := 0
		expandedGalaxy := d.expandGalaxies(exp)
		for j, g1 := range expandedGalaxy {
			for _, g2 := range expandedGalaxy[j:] {
				sum += utils.Abs(g1.x-g2.x) + utils.Abs(g1.y-g2.y)
			}
		}
		fmt.Printf("Part %d: %d\n", i, sum)
	}
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
