package main

import (
	"fmt"
	"github.com/ghjm/advent2023/pkg/utils"
	"os"
	"regexp"
	"strings"
)

type rule struct {
	cat  byte
	cond byte
	qty  int
	dest string
}

type part map[byte]int

type data struct {
	workflows map[string][]rule
	parts     []part
}

type minmax struct {
	min int
	max int
}

type area map[byte]minmax

func (d *data) runWorkflow(p part, wfid string) bool {
	switch wfid {
	case "A":
		return true
	case "R":
		return false
	}
	for _, r := range d.workflows[wfid] {
		if r.cond == 0 {
			return d.runWorkflow(p, r.dest)
		}
		v := p[r.cat]
		if (r.cond == '>' && v > r.qty) || (r.cond == '<' && v < r.qty) {
			return d.runWorkflow(p, r.dest)
		}
	}
	panic("workflow run failed")
}

func (a area) copy() area {
	newA := make(map[byte]minmax)
	for k, v := range a {
		newA[k] = v
	}
	return newA
}

func (d *data) getWholeArea() area {
	a := make(map[byte]minmax)
	for _, p := range d.parts {
		for k := range p {
			a[k] = minmax{1, 4000}
		}
	}
	return a
}

func (d *data) getAcceptRanges(region area, workflow string) []area {
	switch workflow {
	case "A":
		return []area{region}
	case "R":
		return nil
	}
	var result []area
	remaining := region.copy()
	for _, r := range d.workflows[workflow] {
		if r.cond == 0 {
			result = append(result, d.getAcceptRanges(remaining, r.dest)...)
			continue
		}
		switch r.cond {
		case '>':
			curArea := remaining.copy()
			curArea[r.cat] = minmax{r.qty + 1, remaining[r.cat].max}
			remaining[r.cat] = minmax{remaining[r.cat].min, r.qty}
			result = append(result, d.getAcceptRanges(curArea, r.dest)...)
		case '<':
			curArea := remaining.copy()
			curArea[r.cat] = minmax{remaining[r.cat].min, r.qty - 1}
			remaining[r.cat] = minmax{r.qty, remaining[r.cat].max}
			result = append(result, d.getAcceptRanges(curArea, r.dest)...)
		}
	}
	return result
}

func calcSizeOfAreas(areas []area) int64 {
	var sum int64
	for _, a := range areas {
		var size int64
		size = 1
		for _, v := range a {
			size *= int64(v.max - v.min + 1)
		}
		sum += size
	}
	return sum
}

func run() error {
	d := data{
		workflows: make(map[string][]rule),
		parts:     make([]part, 0),
	}

	inRules := true

	outerRuleRE := regexp.MustCompile(`^([a-z]+)\{(.+)\}$`)
	innerRuleRE := regexp.MustCompile(`^([a-z])([<>])(\d+):([a-zA-Z]+)$`)

	err := utils.OpenAndReadLines("input19.txt", func(s string) error {
		switch inRules {
		case true:
			if s == "" {
				inRules = false
				return nil
			}
			mo := outerRuleRE.FindStringSubmatch(s)
			if mo == nil {
				panic("bad data")
			}
			var r []rule
			for _, rs := range strings.Split(mo[2], ",") {
				mi := innerRuleRE.FindStringSubmatch(rs)
				if mi == nil {
					r = append(r, rule{
						dest: rs,
					})
				} else {
					r = append(r, rule{
						cat:  mi[1][0],
						cond: mi[2][0],
						qty:  utils.MustAtoi(mi[3]),
						dest: mi[4],
					})
				}
			}
			d.workflows[mo[1]] = r
		case false:
			s = strings.Trim(s, "{}")
			p := make(map[byte]int)
			for _, v := range strings.Split(s, ",") {
				t := strings.Split(v, "=")
				if len(t) != 2 || len(t[0]) != 1 {
					panic("bad data")
				}
				p[t[0][0]] = utils.MustAtoi(t[1])
			}
			d.parts = append(d.parts, p)
		}
		return nil
	})
	if err != nil {
		return err
	}
	part1 := 0
	for _, p := range d.parts {
		s := d.runWorkflow(p, "in")
		if s {
			for _, v := range p {
				part1 += v
			}
		}
	}
	fmt.Printf("Part 1: %d\n", part1)

	r := d.getAcceptRanges(d.getWholeArea(), "in")
	fmt.Printf("Part 2: %d\n", calcSizeOfAreas(r))

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
