package main

import (
	"fmt"
	"os"
	"strings"

	utils "github.com/ghjm/advent_utils"
)

type data struct {
	cuboids     []utils.StdCuboid
	supporting  utils.MapList[int, int]
	supportedBy utils.MapList[int, int]
}

func (d *data) collapse() {
	delta := utils.Point3D[int]{X: 0, Y: 0, Z: -1}
	changed := true
	for changed {
		changed = false
		for i := range d.cuboids {
			nc := utils.StdCuboid{
				P1: d.cuboids[i].P1.Add(delta),
				P2: d.cuboids[i].P2.Add(delta),
			}
			if nc.P1.Z < 0 {
				continue
			}
			overlap := false
			for j := range d.cuboids {
				if i == j {
					continue
				}
				if d.cuboids[j].Overlaps(nc) {
					overlap = true
					break
				}
			}
			if !overlap {
				d.cuboids[i] = nc
				changed = true
			}
		}
	}
}

func (d *data) calcSupports() {
	d.supporting.Clear()
	for rn := range d.cuboids {
		oc := d.cuboids[rn]
		checkCube := utils.StdCuboid{
			P1: utils.Point3D[int]{
				X: oc.P1.X,
				Y: oc.P1.Y,
				Z: oc.P2.Z,
			},
			P2: utils.Point3D[int]{
				X: oc.P2.X,
				Y: oc.P2.Y,
				Z: oc.P2.Z + 1,
			},
		}
		for i, c := range d.cuboids {
			if i == rn {
				continue
			}
			if checkCube.Overlaps(c) {
				d.supporting.Add(rn, i)
			}
		}
	}
	d.supportedBy.Clear()
	for _, k := range d.supporting.Keys() {
		for _, v := range d.supporting.Get(k) {
			d.supportedBy.Add(v, k)
		}
	}
}

func (d *data) canBeDisintegrated(rn int) bool {
	for _, s := range d.supporting.Get(rn) {
		if len(d.supportedBy.Get(s)) == 1 {
			return false
		}
	}
	return true
}

func (d *data) numThatWouldMove(rn int) int {
	queue := []int{rn}
	affected := make(map[int]struct{})
	affected[rn] = struct{}{}
	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]
		for _, candidate := range d.supporting.Get(item) {
			_, ok := affected[candidate]
			if ok {
				continue
			}
			supportsRemoved := 0
			for _, v := range d.supportedBy.Get(candidate) {
				_, ok = affected[v]
				if ok {
					supportsRemoved++
				}
			}
			if supportsRemoved == len(d.supportedBy.Get(candidate)) {
				affected[candidate] = struct{}{}
			}
			queue = append(queue, candidate)
		}
	}
	return len(affected) - 1
}

func run() error {
	d := data{}
	err := utils.OpenAndReadLines("input22.txt", func(s string) error {
		v := strings.Split(s, "~")
		if len(v) != 2 {
			return fmt.Errorf("invalid input")
		}
		var pts [][]int
		for _, vp := range v {
			p1 := strings.Split(vp, ",")
			if len(p1) != 3 {
				return fmt.Errorf("invalid input")
			}
			var pt []int
			for _, pc := range p1 {
				pt = append(pt, utils.MustAtoi(pc))
			}
			pts = append(pts, pt)
		}
		c := utils.StdCuboid{
			P1: utils.Point3D[int]{
				X: pts[0][0],
				Y: pts[0][1],
				Z: pts[0][2] - 1,
			},
			P2: utils.Point3D[int]{
				X: pts[1][0] + 1,
				Y: pts[1][1] + 1,
				Z: pts[1][2],
			},
		}
		c.OrderCoords()
		d.cuboids = append(d.cuboids, c)
		return nil
	})
	if err != nil {
		return err
	}
	d.collapse()
	d.calcSupports()
	p1 := 0
	p2 := 0
	for i := range d.cuboids {
		if d.canBeDisintegrated(i) {
			p1++
		}
		p2 += d.numThatWouldMove(i)
	}
	fmt.Printf("Part 1: %d\n", p1)
	fmt.Printf("Part 1: %d\n", p2)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
