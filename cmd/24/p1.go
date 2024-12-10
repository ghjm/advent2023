package main

import (
	"fmt"
	"golang.org/x/exp/constraints"
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

func cross3D[T constraints.Integer | constraints.Float](a, b utils.Point3D[T]) utils.Point3D[T] {
	return utils.Point3D[T]{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func dot3D[T constraints.Integer | constraints.Float](a, b utils.Point3D[T]) T {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func sub3D[T constraints.Integer | constraints.Float](a, b utils.Point3D[T]) utils.Point3D[T] {
	return utils.Point3D[T]{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}

func div3D[T constraints.Integer | constraints.Float](a utils.Point3D[T], b T) utils.Point3D[T] {
	return utils.Point3D[T]{
		X: a.X / b,
		Y: a.Y / b,
		Z: a.Z / b,
	}
}

func lin[T constraints.Integer | constraints.Float](r T, a utils.Point3D[T], s T, b utils.Point3D[T], t T, c utils.Point3D[T]) utils.Point3D[T] {
	return utils.Point3D[T]{
		X: r*a.X + s*b.X + t*c.X,
		Y: r*a.Y + s*b.Y + t*c.Y,
		Z: r*a.Z + s*b.Z + t*c.Z,
	}
}

func findPlane(a, b hailstone) (utils.Point3D[float64], float64) {
	p12 := sub3D(a.pos0, b.pos0)
	v12 := sub3D(a.vel, b.vel)
	vv := cross3D(a.vel, b.vel)
	return cross3D(p12, v12), dot3D(p12, vv)
}

func isIndependent(a, b utils.Point3D[float64]) bool {
	c := cross3D(a, b)
	return c.X != 0 || c.Y != 0 || c.Z != 0
}

func findRock(h1, h2, h3 hailstone) float64 {
	a, A := findPlane(h1, h2)
	b, B := findPlane(h1, h3)
	c, C := findPlane(h2, h3)
	w := lin(A, cross3D(b, c), B, cross3D(c, a), C, cross3D(a, b))
	t := dot3D(a, cross3D(b, c))
	w = div3D(w, t)
	w1 := sub3D(h1.vel, w)
	w2 := sub3D(h2.vel, w)
	ww := cross3D(w1, w2)
	E := dot3D(ww, cross3D(h2.pos0, w2))
	F := dot3D(ww, cross3D(h1.pos0, w1))
	G := dot3D(h1.pos0, ww)
	S := dot3D(ww, ww)
	rock := lin(E, w1, -F, w2, G, ww)
	return (rock.X + rock.Y + rock.Z) / S
}

func (d *data) part2() int64 {
	stoneA := d.stones[0]
	var stoneB, stoneC hailstone
	stoneBFound := false
	for i := 1; i < len(d.stones); i++ {
		if !stoneBFound {
			if isIndependent(stoneA.vel, d.stones[i].vel) {
				stoneB = d.stones[i]
				stoneBFound = true
			}
		} else {
			if isIndependent(stoneA.vel, d.stones[i].vel) && isIndependent(stoneB.vel, d.stones[i].vel) {
				stoneC = d.stones[i]
				break
			}
		}
	}
	return int64(math.Round(findRock(stoneA, stoneB, stoneC)))
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
	fmt.Printf("Part 2: %d\n", d.part2())
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
