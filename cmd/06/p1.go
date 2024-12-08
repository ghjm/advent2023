package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	utils "github.com/ghjm/advent_utils"
)

type data struct {
	times            []int
	distances        []int
	longRaceTime     int
	longRaceDistance int
}

func quadraticZeroes(a, b, c float64) (float64, float64) {
	det := b*b - 4*a*c
	if det <= 0 {
		panic("I only do positive determinants")
	}
	return (-b + math.Sqrt(det)) / (2 * a), (-b - math.Sqrt(det)) / (2 * a)
}

func numberOfTimes(time, distance int) int {
	r1, r2 := quadraticZeroes(-1, float64(time), float64(-distance-1))
	if r1 > r2 {
		r2, r1 = r1, r2
	}
	min := int(math.Ceil(r1))
	max := int(math.Floor(r2))
	return max - min + 1
}

func run() error {
	d := data{
		times:     make([]int, 0),
		distances: make([]int, 0),
	}
	err := utils.OpenAndReadLines("input6.txt", func(s string) error {
		var p *[]int
		var lrp *int
		if strings.HasPrefix(s, "Time:") {
			p = &d.times
			lrp = &d.longRaceTime
		} else if strings.HasPrefix(s, "Distance:") {
			p = &d.distances
			lrp = &d.longRaceDistance
		} else {
			return fmt.Errorf("bad data")
		}
		for _, st := range strings.Fields(s)[1:] {
			*p = append(*p, utils.MustAtoi(st))
		}
		*lrp = utils.MustAtoi(strings.Join(strings.Fields(s)[1:], ""))
		return nil
	})
	if err != nil {
		return err
	}
	part1 := 1
	for i := range d.times {
		part1 *= numberOfTimes(d.times[i], d.distances[i])
	}
	fmt.Printf("Part 1: %d\n", part1)
	part2 := numberOfTimes(d.longRaceTime, d.longRaceDistance)
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
