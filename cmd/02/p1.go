package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	utils "github.com/ghjm/advent_utils"
)

type game struct {
	id     int
	rounds []map[string]int
}

func run() error {
	games := make([]game, 0)
	gameRe := regexp.MustCompile(`^Game (\d+): (.+)$`)
	drawRe := regexp.MustCompile(`(\d+) ([a-z]+)`)
	err := utils.OpenAndReadLines("input2.txt", func(s string) error {
		m := gameRe.FindStringSubmatch(s)
		id, err := strconv.Atoi(m[1])
		if err != nil {
			return err
		}
		g := game{
			id:     id,
			rounds: make([]map[string]int, 0),
		}
		for _, r := range strings.Split(m[2], ";") {
			rm := make(map[string]int)
			for _, d := range strings.Split(r, ",") {
				m = drawRe.FindStringSubmatch(d)
				c, err := strconv.Atoi(m[1])
				if err != nil {
					return err
				}
				rm[m[2]] = c
			}
			g.rounds = append(g.rounds, rm)
		}
		games = append(games, g)
		return nil
	})
	if err != nil {
		return err
	}
	part1 := 0
	for _, g := range games {
		gameOk := true
		for _, r := range g.rounds {
			if r["red"] > 12 || r["green"] > 13 || r["blue"] > 14 {
				gameOk = false
				break
			}
		}
		if gameOk {
			part1 += g.id
		}
	}
	fmt.Printf("Part 1: %d\n", part1)
	part2 := 0
	for _, g := range games {
		min := map[string]int{
			"red":   0,
			"green": 0,
			"blue":  0,
		}
		for _, r := range g.rounds {
			for k := range min {
				if r[k] > min[k] {
					min[k] = r[k]
				}
			}
		}
		part2 += min["red"] * min["green"] * min["blue"]
	}
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
