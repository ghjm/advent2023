package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	utils "github.com/ghjm/advent_utils"
)

type data struct {
	steps []string
}

type lens struct {
	id          string
	focalLength int
}

func hash(s string) int {
	h := 0
	for _, c := range s {
		h += int(c)
		h *= 17
		h %= 256
	}
	return h
}

func run() error {
	d := data{}
	err := utils.OpenAndReadLines("input15.txt", func(s string) error {
		d.steps = append(d.steps, strings.Split(s, ",")...)
		return nil
	})
	if err != nil {
		return err
	}
	part1 := 0
	for _, s := range d.steps {
		part1 += hash(s)
	}
	fmt.Printf("Part 1: %d\n", part1)

	instrRe := regexp.MustCompile(`^([a-z]+)([-=])(\d*)$`)
	boxes := make(map[int][]lens)
	for _, s := range d.steps {
		m := instrRe.FindStringSubmatch(s)
		if m == nil {
			return fmt.Errorf("non-matching instruction")
		}
		id, op := m[1], m[2]
		var length int
		if m[3] != "" {
			length = utils.MustAtoi(m[3])
		}
		h := hash(id)
		_, ok := boxes[h]
		if !ok {
			boxes[h] = make([]lens, 0)
		}
		switch op {
		case "-":
			for i, b := range boxes[h] {
				if b.id == id {
					boxes[h] = append(boxes[h][:i], boxes[h][i+1:]...)
					break
				}
			}
		case "=":
			found := false
			for i, b := range boxes[h] {
				if b.id == id {
					boxes[h][i].focalLength = length
					found = true
					break
				}
			}
			if !found {
				boxes[h] = append(boxes[h], lens{
					id:          id,
					focalLength: length,
				})
			}
		}
	}

	part2 := 0
	for i, b := range boxes {
		for j, l := range b {
			part2 += (i + 1) * (j + 1) * l.focalLength
		}
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
