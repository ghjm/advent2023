package main

import (
	"fmt"
	"github.com/ghjm/advent2023/pkg/utils"
	"os"
)

type data struct {
	puzzles [][]string
}

func rotate(m []string) []string {
	result := make([]string, len(m[0]))
	for ry := range result {
		s := ""
		for rx := range m {
			s += string([]byte{m[rx][ry]})
		}
		result[ry] = s
	}
	return result
}

func checkSymmetry(m []string, line int) bool {
	p1 := line
	p2 := line + 1
	for {
		if p1 < 0 {
			return true
		}
		if p2 >= len(m) {
			return true
		}
		if m[p1] != m[p2] {
			return false
		}
		p1--
		p2++
	}
}

func opposite(c byte) byte {
	switch c {
	case '.':
		return '#'
	case '#':
		return '.'
	default:
		panic("impossible map char")
	}
}

func checkPuzzle(puz []string, skip ...int) int {
loop:
	for i := 0; i < len(puz)-1; i++ {
		if puz[i] == puz[i+1] {
			if checkSymmetry(puz, i) {
				for _, s := range skip {
					if i == s {
						continue loop
					}
				}
				return i
			}
		}
	}
	return -1
}

func run() error {
	d := data{}
	var curPuz []string
	err := utils.OpenAndReadLines("input13.txt", func(s string) error {
		if s == "" {
			d.puzzles = append(d.puzzles, curPuz)
			curPuz = nil
		} else {
			curPuz = append(curPuz, s)
		}
		return nil
	})
	if curPuz != nil {
		d.puzzles = append(d.puzzles, curPuz)
	}
	if err != nil {
		return err
	}
	part1 := 0
	part2 := 0
	for pi, pp := range d.puzzles {
		puzWithRot := [][]string{pp, rotate(pp)}
		sol1 := -1
		sol1Prn := -1
	puz1:
		for prn, puz := range puzWithRot {
			sol := checkPuzzle(puz)
			if sol >= 0 {
				sol1 = sol
				sol1Prn = prn
				if prn == 0 {
					part1 += 100 * (sol1 + 1)
				} else {
					part1 += sol1 + 1
				}
				break puz1
			}
		}
		if sol1 == -1 {
			return fmt.Errorf("puzzle %d had no part 1 solution", pi)
		}
		sol2 := -1
	puz2:
		for prn, puz := range puzWithRot {
			for y := range puz {
				for x := range puz[0] {
					var newPuz []string
					for _, p := range puz {
						newPuz = append(newPuz, p)
					}
					newPuz[y] = newPuz[y][:x] + string([]byte{opposite(newPuz[y][x])}) + newPuz[y][x+1:]
					var sol int
					if prn == sol1Prn {
						sol = checkPuzzle(newPuz, sol1)
					} else {
						sol = checkPuzzle(newPuz)
					}
					if sol >= 0 {
						sol2 = sol
						if prn == 0 {
							part2 += 100 * (sol2 + 1)
						} else {
							part2 += sol2 + 1
						}
						break puz2
					}
				}
			}
		}
		if sol2 == -1 {
			return fmt.Errorf("puzzle %d had no part 2 solution", pi)
		}
	}
	fmt.Printf("Part 1: %d\nPart 2: %d\n", part1, part2)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
