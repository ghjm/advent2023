package main

import (
	"fmt"
	"os"
	"strconv"
	"unicode"

	utils "github.com/ghjm/advent_utils"
)

type grid [][]rune
type point struct {
	x int
	y int
}
type gears map[point][]int

type data struct {
	g         grid
	partNums  []int
	gearParts gears
}

func (d *data) checkDig(digstr []rune, x int, y int) error {
	pn, err := strconv.Atoi(string(digstr))
	if err != nil {
		return err
	}
	points := make(map[point]struct{})
	for i := range digstr {
		for _, dy := range []int{-1, 0, 1} {
			for _, dx := range []int{-1, 0, 1} {
				cy := y + dy
				cx := x + dx - len(digstr) + i
				if cy >= 0 && cx >= 0 && cy < len(d.g) && cx < len(d.g[y]) {
					points[point{cx, cy}] = struct{}{}
				}
			}
		}
	}
	for p := range points {
		chk := d.g[p.y][p.x]
		if !unicode.IsDigit(chk) && chk != '.' {
			d.partNums = append(d.partNums, pn)
		}
		if chk == '*' {
			_, ok := d.gearParts[p]
			if !ok {
				d.gearParts[p] = make([]int, 0)
			}
			d.gearParts[p] = append(d.gearParts[p], pn)
		}
	}
	return nil
}

func run() error {
	d := data{
		g:         make(grid, 0),
		partNums:  make([]int, 0),
		gearParts: make(gears),
	}
	err := utils.OpenAndReadLines("input3.txt", func(s string) error {
		d.g = append(d.g, []rune(s))
		return nil
	})
	if err != nil {
		return err
	}
	for y := range d.g {
		digstr := make([]rune, 0)
		for x := range d.g[y] {
			c := d.g[y][x]
			if unicode.IsDigit(c) {
				digstr = append(digstr, c)
			} else if len(digstr) > 0 {
				err = d.checkDig(digstr, x, y)
				if err != nil {
					return err
				}
				digstr = make([]rune, 0)
			}
		}
		if len(digstr) > 0 {
			err = d.checkDig(digstr, len(d.g[y]), y)
			if err != nil {
				return err
			}
		}
	}
	part1 := 0
	for _, pn := range d.partNums {
		part1 += pn
	}
	fmt.Printf("Part 1: %d\n", part1)
	part2 := 0
	for _, gp := range d.gearParts {
		if len(gp) == 2 {
			part2 += gp[0] * gp[1]
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
