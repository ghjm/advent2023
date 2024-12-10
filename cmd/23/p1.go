package main

import (
	"fmt"
	utils "github.com/ghjm/advent_utils"
	"github.com/ghjm/advent_utils/board"
	"os"
	"time"
)

type data struct {
	board      *board.StdBoard
	tmpBoard   *board.StdBoard
	startPoint utils.StdPoint
	endPoint   utils.StdPoint
}

var directions = []struct {
	d     utils.StdPoint
	slope rune
}{
	{d: utils.StdPoint{X: -1, Y: 0}, slope: '<'},
	{d: utils.StdPoint{X: 1, Y: 0}, slope: '>'},
	{d: utils.StdPoint{X: 0, Y: -1}, slope: '^'},
	{d: utils.StdPoint{X: 0, Y: 1}, slope: 'v'},
}

func (d *data) legalMovesWithSlope(p utils.StdPoint) []utils.StdPoint {
	var results []utils.StdPoint
	for _, dir := range directions {
		np := p.Add(dir.d)
		if !d.board.Contains(np) {
			continue
		}
		nv := d.tmpBoard.Get(np)
		if nv == '.' || nv == dir.slope {
			results = append(results, np)
		}
	}
	return results
}

func (d *data) legalMovesIgnoringSlope(p utils.StdPoint) []utils.StdPoint {
	var results []utils.StdPoint
	for _, dir := range directions {
		np := p.Add(dir.d)
		if !d.tmpBoard.Contains(np) {
			continue
		}
		if d.tmpBoard.Get(np) != '#' {
			results = append(results, np)
		}
	}
	return results
}

func (d *data) search(start utils.StdPoint, legalMoves func(utils.StdPoint) []utils.StdPoint) int {
	if d.tmpBoard.Get(start) == 'O' {
		return 0
	}
	if start.Equal(d.endPoint) {
		sum := 0
		d.tmpBoard.Iterate(func(p utils.Point[int], v rune) bool {
			if v == 'O' {
				sum++
			}
			return true
		})
		return sum
	}
	d.tmpBoard.Set(start, 'O')
	defer d.tmpBoard.Set(start, d.board.Get(start))
	maxLength := 0
	for _, np := range legalMoves(start) {
		nextLength := d.search(np, legalMoves)
		if nextLength > maxLength {
			maxLength = nextLength
		}
	}
	return maxLength
}

func run() error {
	d := data{
		board: board.NewStdBoard(board.WithStorage(&board.FlatBoard{})),
	}
	err := d.board.FromFile("input23.txt")
	if err != nil {
		return err
	}
	for x := 0; x < d.board.Bounds().Width(); x++ {
		sp := utils.StdPoint{X: x, Y: 0}
		v := d.board.Get(sp)
		if v == '.' {
			d.startPoint = sp
			break
		}
	}
	for x := 0; x < d.board.Bounds().Width(); x++ {
		sp := utils.StdPoint{X: x, Y: d.board.Bounds().Height() - 1}
		if d.board.Get(sp) == '.' {
			d.endPoint = sp
			break
		}
	}
	start := time.Now()
	d.tmpBoard = d.board.Copy()
	fmt.Printf("Part 1: %d\n", d.search(d.startPoint, d.legalMovesWithSlope))
	fmt.Printf("time: %v\n", time.Since(start))
	d.tmpBoard = d.board.Copy()
	fmt.Printf("Part 2: %d\n", d.search(d.startPoint, d.legalMovesIgnoringSlope))
	fmt.Printf("time: %v\n", time.Since(start))
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
