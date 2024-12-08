package main

import (
	"fmt"
	"math"
	"os"

	utils "github.com/ghjm/advent_utils"
	"github.com/oleiade/lane/v2"
)

type data struct {
	lines []string
}

type point struct {
	x, y int
}

type state struct {
	pos         point
	direction   point
	numStraight int
}

func manhattanDistance(p1, p2 point) int {
	return utils.Abs(p1.x-p2.x) + utils.Abs(p1.y-p2.y)
}

var leftTurn = map[point]point{
	{0, -1}: {-1, 0},
	{-1, 0}: {0, 1},
	{0, 1}:  {1, 0},
	{1, 0}:  {0, -1},
}

var rightTurn = map[point]point{
	{0, -1}: {1, 0},
	{1, 0}:  {0, 1},
	{0, 1}:  {-1, 0},
	{-1, 0}: {0, -1},
}

func (d *data) heatLoss(p point) int {
	return utils.MustAtoi(string(d.lines[p.y][p.x]))
}

func (d *data) neighbors(s state, part int) []state {
	var result []state
	if s.direction.x == 0 && s.direction.y == 0 {
		for _, dp := range []point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			np := point{s.pos.x + dp.x, s.pos.y + dp.y}
			if np.x >= 0 && np.y >= 0 && np.x < len(d.lines[0]) && np.y < len(d.lines) {
				result = append(result, state{
					pos:         np,
					direction:   dp,
					numStraight: 0,
				})
			}
		}
		return result
	}
	var dirs []point
	if part == 1 {
		dirs = append(dirs, leftTurn[s.direction], rightTurn[s.direction])
		if s.numStraight < 2 {
			dirs = append(dirs, s.direction)
		}
	} else {
		if s.numStraight < 9 {
			dirs = append(dirs, s.direction)
		}
		if s.numStraight >= 3 {
			dirs = append(dirs, leftTurn[s.direction], rightTurn[s.direction])
		}
	}
	for _, dp := range dirs {
		np := point{s.pos.x + dp.x, s.pos.y + dp.y}
		if np.x >= 0 && np.y >= 0 && np.x < len(d.lines[0]) && np.y < len(d.lines) {
			newNs := 0
			if dp == s.direction {
				newNs = s.numStraight + 1
			}
			result = append(result, state{
				pos:         np,
				direction:   dp,
				numStraight: newNs,
			})
		}
	}
	return result
}

func (d *data) findMinPath(part int) ([]state, int) {
	goal := point{len(d.lines[0]) - 1, len(d.lines) - 1}
	start := state{point{0, 0}, point{0, 0}, 0}
	pq := lane.NewMinPriorityQueue[state, int]()
	pq.Push(start, manhattanDistance(start.pos, goal))
	gs := utils.NewDefaultMap[state, int](math.MaxInt32)
	cameFrom := make(map[state]state)
	visited := make(map[state]struct{})
	gs.Set(start, 0)
	for {
		s, _, ok := pq.Pop()
		if !ok {
			break
		}
		visited[s] = struct{}{}
		if s.pos == goal {
			current := s
			path := []state{current}
			heatLoss := d.heatLoss(current.pos)
			for {
				prev, ok := cameFrom[current]
				if !ok {
					break
				}
				current = prev
				path = append([]state{current}, path...)
				heatLoss += d.heatLoss(current.pos)
			}
			heatLoss -= d.heatLoss(current.pos) // starting square heat loss doesn't count
			return path, heatLoss
		}
		for _, neigh := range d.neighbors(s, part) {
			_, ok := visited[neigh]
			if ok {
				continue
			}
			tentativeScore := gs.Get(s) + d.heatLoss(neigh.pos)
			if tentativeScore < gs.Get(neigh) {
				cameFrom[neigh] = s
				gs.Set(neigh, tentativeScore)
				pq.Push(neigh, tentativeScore+manhattanDistance(neigh.pos, goal))
			}
		}
	}
	panic("min path not found")
}

func run() error {
	d := data{}
	err := utils.OpenAndReadLines("input17.txt", func(s string) error {
		d.lines = append(d.lines, s)
		return nil
	})
	if err != nil {
		return err
	}
	_, part1 := d.findMinPath(1)
	fmt.Printf("Part 1: %d\n", part1)
	_, part2 := d.findMinPath(2)
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
