package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	utils "github.com/ghjm/advent_utils"
)

type hailstone struct {
	pos0 utils.Point3D[float64]
	vel  utils.Point3D[float64]
}

type data struct {
	stones []hailstone
}

const EPSILON float64 = 0.000000000001

// converts a hailstone line to the form ax+by=c
func getLineTerms(h hailstone) (a, b, c float64) {
	a = h.vel.Y
	b = -h.vel.X
	c = h.vel.Y*h.pos0.X - h.vel.X*h.pos0.Y
	return
}

func findIntersect(stone1, stone2 hailstone) (utils.Point3D[float64], bool) {
	a1, b1, c1 := getLineTerms(stone1)
	a2, b2, c2 := getLineTerms(stone2)
	if a1*b2-a2*b1 < EPSILON {
		// parallel lines
		return utils.Point3D[float64]{}, false
	}
	result := utils.Point3D[float64]{}
	result.X = (c1*b2 - b1*c2) / (a1*b2 - b1*a2)
	result.Y = (a1*c2 - c1*a2) / (a1*b2 - b1*a2)
	return result, true
}

func findTimeAtPoint(stone hailstone, p utils.Point3D[float64]) (float64, bool) {
	if math.Abs(stone.vel.X) < EPSILON && math.Abs(stone.vel.Y) < EPSILON {
		if stone.pos0.X == p.X && stone.pos0.Y == p.Y {
			return 0, true
		} else {
			return 0, false
		}
	} else if math.Abs(stone.vel.X) < EPSILON {
		if stone.pos0.X != p.X {
			return 0, false
		}
		return (p.Y - stone.pos0.Y) / stone.vel.Y, true
	} else if math.Abs(stone.vel.Y) < EPSILON {
		if stone.pos0.Y != p.Y {
			return 0, false
		}
		return (p.X - stone.pos0.X) / stone.vel.X, true
	} else {
		return (p.X - stone.pos0.X) / stone.vel.X, true
	}
}

func (d *data) part1() int {
	count := 0
	for i, sA := range d.stones {
		for j, sB := range d.stones {
			if i == j {
				continue
			}
			p, ok := findIntersect(sA, sB)
			if !ok {
				continue
			}
			if p.X < 200000000000000 || p.X > 400000000000000 || p.Y < 200000000000000 || p.Y > 400000000000000 {
				continue
			}
			t, ok := findTimeAtPoint(sA, p)
			if !ok {
				panic("error")
			}
			if t < 0 {
				continue
			}
			t, ok = findTimeAtPoint(sB, p)
			if !ok {
				panic("error")
			}
			if t < 0 {
				continue
			}
			count++
		}
	}
	return count
}

func run() error {
	d := data{}
	err := utils.OpenAndReadLines("input24.txt", func(s string) error {
		var stone hailstone
		s = strings.Replace(s, " ", "", -1)
		s1 := strings.Split(s, "@")
		if len(s1) != 2 {
			return fmt.Errorf("invalid input")
		}
		p1 := strings.Split(s1[0], ",")
		stone.pos0 = utils.Point3D[float64]{
			X: utils.MustAtof(p1[0]),
			Y: utils.MustAtof(p1[1]),
			Z: utils.MustAtof(p1[2]),
		}
		v1 := strings.Split(s1[1], ",")
		stone.vel = utils.Point3D[float64]{
			X: utils.MustAtof(v1[0]),
			Y: utils.MustAtof(v1[1]),
			Z: utils.MustAtof(v1[2]),
		}
		d.stones = append(d.stones, stone)
		return nil
	})
	if err != nil {
		return err
	}
	fmt.Printf("Part 1: %d\n", d.part1())
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
