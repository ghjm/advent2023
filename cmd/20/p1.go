package main

import (
	"fmt"
	"github.com/ghjm/advent2023/pkg/utils"
	"os"
	"regexp"
	"strings"
)

type data struct {
	components map[string]component
}

type pulse struct {
	from string
	dest string
	high bool
}

type component interface {
	RecvPulse(from string, high bool) []pulse
	Neighbors() []string
	Reset()
}

func sendToNeighbors(from string, high bool, neighbors []string) []pulse {
	var result []pulse
	for _, n := range neighbors {
		result = append(result, pulse{
			from: from,
			dest: n,
			high: high,
		})
	}
	return result
}

type broadcaster struct {
	id        string
	neighbors []string
}

func (c *broadcaster) RecvPulse(_ string, high bool) []pulse {
	return sendToNeighbors(c.id, high, c.neighbors)
}

func (c *broadcaster) Neighbors() []string {
	return c.neighbors
}

func (c *broadcaster) Reset() {
}

type flipflop struct {
	id        string
	neighbors []string
	on        bool
}

func (c *flipflop) RecvPulse(_ string, high bool) []pulse {
	if high {
		return nil
	}
	c.on = !c.on
	return sendToNeighbors(c.id, c.on, c.neighbors)
}

func (c *flipflop) Neighbors() []string {
	return c.neighbors
}

func (c *flipflop) Reset() {
	c.on = false
}

type conjunction struct {
	id         string
	neighbors  []string
	lastPulses map[string]bool
}

func (c *conjunction) RecvPulse(from string, high bool) []pulse {
	c.lastPulses[from] = high
	allHigh := true
	for _, p := range c.lastPulses {
		if !p {
			allHigh = false
			break
		}
	}
	return sendToNeighbors(c.id, !allHigh, c.neighbors)
}

func (c *conjunction) Neighbors() []string {
	return c.neighbors
}

func (c *conjunction) Reset() {
	for k := range c.lastPulses {
		c.lastPulses[k] = false
	}
}

func (d *data) pushButton() (int, int, bool) {
	pending := []pulse{{
		from: "button",
		dest: "broadcaster",
		high: false,
	}}
	lowPulses := 0
	highPulses := 0
	rxLowPulse := false
	for len(pending) > 0 {
		curPulse := pending[0]
		pending = pending[1:]
		if curPulse.high {
			highPulses++
		} else {
			lowPulses++
		}
		if !curPulse.high && curPulse.dest == "rx" {
			rxLowPulse = true
		}
		comp, ok := d.components[curPulse.dest]
		if ok {
			pending = append(pending, comp.RecvPulse(curPulse.from, curPulse.high)...)
		}
	}
	return lowPulses, highPulses, rxLowPulse
}

func run() error {
	d := data{
		components: make(map[string]component),
	}
	connRE := regexp.MustCompile(`^(.*) -> (.*)$`)
	err := utils.OpenAndReadLines("input20.txt", func(s string) error {
		m := connRE.FindStringSubmatch(s)
		if m == nil {
			panic("bad data")
		}
		neighbors := strings.Split(m[2], ", ")
		switch {
		case m[1] == "broadcaster":
			c := &broadcaster{
				id:        m[1],
				neighbors: neighbors,
			}
			d.components[m[1]] = c
		case m[1][0] == '&':
			c := &conjunction{
				id:         m[1][1:],
				neighbors:  neighbors,
				lastPulses: make(map[string]bool),
			}
			d.components[m[1][1:]] = c
		case m[1][0] == '%':
			c := &flipflop{
				id:        m[1][1:],
				neighbors: neighbors,
				on:        false,
			}
			d.components[m[1][1:]] = c
		}
		return nil
	})
	if err != nil {
		return err
	}
	for cn, cv := range d.components {
		cc, ok := cv.(*conjunction)
		if ok {
			for nn, nv := range d.components {
				for _, s := range nv.Neighbors() {
					if s == cn {
						cc.lastPulses[nn] = false
						break
					}
				}
			}
		}
	}
	var lowPulses, highPulses int
	for i := 0; i < 1000; i++ {
		lp, hp, _ := d.pushButton()
		lowPulses += lp
		highPulses += hp
	}
	fmt.Printf("Part 1: %d\n", lowPulses*highPulses)
	for _, c := range d.components {
		c.Reset()
	}
	i := 0
	for {
		i++
		_, _, rxLowPulse := d.pushButton()
		if rxLowPulse {
			break
		}
	}
	fmt.Printf("Part 2: %d\n", i)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
