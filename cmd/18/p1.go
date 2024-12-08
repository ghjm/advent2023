package main

import (
	"fmt"
	"os"
	"regexp"

	utils "github.com/ghjm/advent_utils"
)

type instruction struct {
	p1Dir  byte
	p1Dist int64
	p2Dir  byte
	p2Dist int64
}

type point struct {
	x, y int64
}

type data struct {
	inst []instruction
}

var dirs = map[byte]point{
	'U': {0, -1},
	'D': {0, 1},
	'L': {-1, 0},
	'R': {1, 0},
}

func (d *data) getVertices(part int) []point {
	var prelim []point
	curPos := point{0, 0}
	var minX int64
	var minY int64
	for _, instr := range d.inst {
		var dir point
		var dist int64
		if part == 1 {
			dir = dirs[instr.p1Dir]
			dist = instr.p1Dist
		} else {
			dir = dirs[instr.p2Dir]
			dist = instr.p2Dist
		}
		curPos = point{curPos.x + dir.x*dist, curPos.y + dir.y*dist}
		if curPos.x < minX {
			minX = curPos.x
		}
		if curPos.y < minY {
			minY = curPos.y
		}
		prelim = append(prelim, curPos)
	}
	var result []point
	for _, c := range prelim {
		result = append(result, point{c.x - minX, c.y - minY})
	}
	return result
}

func shoelaceArea(points []point) int64 {
	var sum int64
	for i, p := range points {
		np := points[(i+1)%len(points)]
		sum += (p.y + np.y) * (p.x - np.x)
		sum += utils.Abs64(np.y-p.y) + utils.Abs64(np.x-p.x)
	}
	if sum%2 == 1 {
		panic("sum is odd")
	}
	return (sum / 2) + 1 // Why +1?  I don't know
}

var hexToDir = map[byte]byte{
	0: 'R',
	1: 'D',
	2: 'L',
	3: 'U',
}

func run() error {
	instRE := regexp.MustCompile(`^([RLDU]) (\d+) \(#([a-z0-9]+)\)$`)
	d := data{}
	err := utils.OpenAndReadLines("input18.txt", func(s string) error {
		m := instRE.FindStringSubmatch(s)
		d.inst = append(d.inst, instruction{
			p1Dir:  m[1][0],
			p1Dist: utils.MustAtoi64(m[2]),
			p2Dir:  hexToDir[byte(utils.MustAtoiHex64(string(m[3][len(m[3])-1])))],
			p2Dist: utils.MustAtoiHex64(m[3][0 : len(m[3])-1]),
		})
		return nil
	})
	if err != nil {
		return err
	}
	v := d.getVertices(1)
	fmt.Printf("Part 1: %d\n", shoelaceArea(v))
	v = d.getVertices(2)
	fmt.Printf("Part 2: %d\n", shoelaceArea(v))
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
