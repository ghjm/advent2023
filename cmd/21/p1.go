package main

import (
	"fmt"
	utils "github.com/ghjm/advent_utils"
	"github.com/ghjm/advent_utils/board"
	"os"
)

type data struct {
	board    *board.StdBoard
	startPos utils.StdPoint
}

func (d *data) legalMovesFrom(p utils.StdPoint) []utils.StdPoint {
	var results []utils.StdPoint
	for _, dir := range []utils.StdPoint{{-1, 0}, {1, 0}, {0, 1}, {0, -1}} {
		np := p.Add(dir)
		if d.board.Contains(np) {
			v := d.board.Get(np)
			if v == '.' || v == 'S' {
				results = append(results, np)
			}
		}
	}
	return results
}

func run() error {
	d := data{
		board: board.NewStdBoard(board.WithStorage(&board.FlatBoard{})),
	}
	d.board.MustFromFile("input21.txt")
	d.board.Iterate(func(p utils.StdPoint, v rune) bool {
		if v == 'S' {
			d.startPos = p
			return false
		}
		return true
	})
	open := []utils.PointPlusData[int, int]{{d.startPos, 0}}
	visited := make(map[utils.StdPoint]int)
	for len(open) > 0 {
		s := open[0]
		open = open[1:]
		_, ok := visited[s.Point]
		if ok {
			continue
		}
		visited[s.Point] = s.Data
		for _, m := range d.legalMovesFrom(s.Point) {
			_, ok := visited[m]
			if !ok {
				open = append(open, utils.PointPlusData[int, int]{Point: m, Data: s.Data + 1})
			}
		}
	}
	var p1 int
	for _, v := range visited {
		if v <= 64 && v%2 == 0 {
			p1++
		}
	}
	fmt.Printf("Part 1: %d\n", p1)

	p2Steps := uint64(26501365)
	w := uint64(d.board.Bounds().Width())
	n := (p2Steps - (w / 2)) / w
	even := n * n
	odd := (n + 1) * (n + 1)
	var oddLocs, evenLocs, oddCorners, evenCorners uint64
	for _, v := range visited {
		if v%2 == 0 {
			evenLocs++
		} else {
			oddLocs++
		}
		if v > 65 {
			if v%2 == 0 {
				evenCorners++
			} else {
				oddCorners++
			}
		}
	}
	p2 := odd*oddLocs + even*evenLocs - (n+1)*oddCorners + n*evenCorners
	fmt.Printf("Part 2: %d\n", p2)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
