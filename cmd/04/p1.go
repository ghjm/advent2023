package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	utils "github.com/ghjm/advent_utils"
)

type card struct {
	id      int
	winners map[int]struct{}
	picked  map[int]struct{}
}

type data struct {
	cards     []card
	cardIndex map[int]card
}

func (d *data) makeIndex() {
	d.cardIndex = make(map[int]card)
	for _, c := range d.cards {
		d.cardIndex[c.id] = c
	}
}

func (d *data) findCard(id int) (card, error) {
	if d.cardIndex == nil {
		d.makeIndex()
	}
	c, ok := d.cardIndex[id]
	if ok {
		return c, nil
	}
	return card{}, fmt.Errorf("not found")
}

func (c card) numWinners() int {
	result := 0
	for p := range c.picked {
		_, ok := c.winners[p]
		if ok {
			result += 1
		}
	}
	return result
}

func run() error {
	d := data{
		cards: make([]card, 0),
	}
	cardRe := regexp.MustCompile(`^Card +(\d+): (.+) \| (.+)$`)
	err := utils.OpenAndReadLines("input4.txt", func(s string) error {
		m := cardRe.FindStringSubmatch(s)
		id, err := strconv.Atoi(m[1])
		if err != nil {
			return err
		}
		c := card{
			id:      id,
			winners: make(map[int]struct{}),
			picked:  make(map[int]struct{}),
		}
		for k, v := range map[string]map[int]struct{}{
			m[2]: c.winners,
			m[3]: c.picked,
		} {
			for _, ns := range strings.Split(k, " ") {
				if len(ns) > 0 {
					n, err := strconv.Atoi(ns)
					if err != nil {
						return err
					}
					v[n] = struct{}{}
				}
			}
		}
		d.cards = append(d.cards, c)
		return nil
	})
	if err != nil {
		return err
	}
	part1 := 0
	for _, c := range d.cards {
		nw := c.numWinners()
		points := 0
		for i := 0; i < nw; i++ {
			if points == 0 {
				points = 1
			} else {
				points *= 2
			}
		}
		part1 += points
	}
	fmt.Printf("Part 1: %d\n", part1)
	open := make([]card, 0)
	closed := make([]card, 0)
	for _, c := range d.cards {
		open = append(open, c)
	}
	for len(open) > 0 {
		c := open[0]
		open = open[1:]
		closed = append(closed, c)
		nw := c.numWinners()
		for i := c.id + 1; i < c.id+1+nw; i++ {
			nc, err := d.findCard(i)
			if err != nil {
				return err
			}
			open = append(open, nc)
		}
	}
	fmt.Printf("Part 2: %d\n", len(closed))
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
