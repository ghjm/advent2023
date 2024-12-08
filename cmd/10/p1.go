package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	utils "github.com/ghjm/advent_utils"
)

type point struct {
	x, y int
}

type data struct {
	cells     [][]byte
	start     point
	startType byte
}

func (d *data) getCell(p point) byte {
	return d.cells[p.y][p.x]
}

func (d *data) neighbors(p point) []point {
	pCh := d.getCell(p)
	if pCh == 'S' {
		pCh = d.startType
	}
	n1 := make([]point, 0)
	switch pCh {
	case '|':
		n1 = append(n1, point{p.x, p.y - 1}, point{p.x, p.y + 1})
	case '-':
		n1 = append(n1, point{p.x - 1, p.y}, point{p.x + 1, p.y})
	case 'L':
		n1 = append(n1, point{p.x, p.y - 1}, point{p.x + 1, p.y})
	case 'J':
		n1 = append(n1, point{p.x - 1, p.y}, point{p.x, p.y - 1})
	case '7':
		n1 = append(n1, point{p.x - 1, p.y}, point{p.x, p.y + 1})
	case 'F':
		n1 = append(n1, point{p.x + 1, p.y}, point{p.x, p.y + 1})
	}
	n2 := make([]point, 0)
	for _, np := range n1 {
		if np.x >= 0 && np.y >= 0 && np.x < len(d.cells[0]) && np.y < len(d.cells) {
			n2 = append(n2, np)
		}
	}
	return n2
}

func (d *data) calcStartType() {
	for _, c := range "|-LJ7F" {
		d.startType = byte(c)
		neigh := d.neighbors(d.start)
		if len(neigh) != 2 {
			continue
		}
		good := true
		for _, n := range neigh {
			pGood := false
			for _, n2 := range d.neighbors(n) {
				if n2 == d.start {
					pGood = true
					break
				}
			}
			if !pGood {
				good = false
				break
			}
		}
		if good {
			return
		}
	}
}

func run() error {
	d := data{
		cells: make([][]byte, 0),
	}
	err := utils.OpenAndReadLines("input10.txt", func(s string) error {
		d.cells = append(d.cells, []byte(s))
		idx := strings.Index(s, "S")
		if idx != -1 {
			d.start = point{
				x: idx,
				y: len(d.cells) - 1,
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	d.calcStartType()
	distances := make(map[point]int)
	distances[d.start] = 0
	traces := d.neighbors(d.start)
	for _, t := range traces {
		distances[t] = 1
	}
	finished := false
	for !finished {
		finished = true
		for i := range traces {
			od := distances[traces[i]]
			for _, n := range d.neighbors(traces[i]) {
				_, ok := distances[n]
				if !ok {
					traces[i] = n
					distances[n] = od + 1
					finished = false
					break
				}
			}
		}
	}
	part1 := 0
	for _, dist := range distances {
		if dist > part1 {
			part1 = dist
		}
	}
	fmt.Printf("Part 1: %d\n", part1)

	// credit to ricbit for this solution
	re1 := regexp.MustCompile(`F-*7|L-*J`)
	re2 := regexp.MustCompile(`F-*J|L-*7`)
	part2 := 0
	d.cells[d.start.y][d.start.x] = d.startType
	for y := range d.cells {
		for x := range d.cells[y] {
			_, ok := distances[point{x, y}]
			if !ok {
				d.cells[y][x] = '.'
			}
		}
		s := string(d.cells[y])
		s = re1.ReplaceAllLiteralString(s, "")
		s = re2.ReplaceAllLiteralString(s, "|")
		interior := 0
		for _, c := range s {
			if c == '|' {
				interior += 1
			}
			if interior%2 == 1 && c == '.' {
				part2++
			}
		}
	}
	fmt.Printf("Part 2: %d\n", part2)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
