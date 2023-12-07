package main

import (
	"fmt"
	"github.com/ghjm/advent2023/pkg/utils"
	"os"
	"sort"
	"strings"
)

type hand struct {
	cards            string
	bid              int
	handType         int
	wildcardHandType int
}

type data struct {
	hands []hand
}

func cardsLess(c1 string, c2 string, cardValues map[byte]int) bool {
	if len(c1) != len(c2) {
		panic("impossible cards")
	}
	for i := range c1 {
		if cardValues[c1[i]] < cardValues[c2[i]] {
			return true
		}
		if cardValues[c1[i]] > cardValues[c2[i]] {
			return false
		}
	}
	return false
}

func cardsLessP1(c1 string, c2 string) bool {
	return cardsLess(c1, c2, map[byte]int{'2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, 'T': 10, 'J': 11, 'Q': 12, 'K': 13, 'A': 14})
}

func cardsLessP2(c1 string, c2 string) bool {
	return cardsLess(c1, c2, map[byte]int{'J': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, 'T': 10, 'Q': 12, 'K': 13, 'A': 14})
}

func (d *data) lessThanP1(i int, j int) bool {
	hi, hj := d.hands[i], d.hands[j]
	if hi.handType < hj.handType {
		return true
	}
	if hi.handType > hj.handType {
		return false
	}
	return cardsLessP1(hi.cards, hj.cards)
}

func (d *data) lessThanP2(i int, j int) bool {
	hi, hj := d.hands[i], d.hands[j]
	if hi.wildcardHandType < hj.wildcardHandType {
		return true
	}
	if hi.wildcardHandType > hj.wildcardHandType {
		return false
	}
	return cardsLessP2(hi.cards, hj.cards)
}

func makefreqs(cards string) []int {
	fc := make(map[rune]int)
	for _, c := range cards {
		_, ok := fc[c]
		if ok {
			fc[c] += 1
		} else {
			fc[c] = 1
		}
	}
	result := make([]int, 0)
	for _, v := range fc {
		result = append(result, v)
	}
	sort.Ints(result)
	return result
}

func handType(freqs []int) int {
	switch {
	case len(freqs) == 1 && freqs[0] == 5:
		// 5 of a kind
		return 7
	case len(freqs) == 2 && freqs[0] == 1 && freqs[1] == 4:
		// 4 of a kind
		return 6
	case len(freqs) == 2 && freqs[0] == 2 && freqs[1] == 3:
		// full house
		return 5
	case len(freqs) == 3 && freqs[0] == 1 && freqs[1] == 1 && freqs[2] == 3:
		// three of a kind
		return 4
	case len(freqs) == 3 && freqs[0] == 1 && freqs[1] == 2 && freqs[2] == 2:
		// two pair
		return 3
	case len(freqs) == 4 && freqs[0] == 1 && freqs[1] == 1 && freqs[2] == 1 && freqs[3] == 2:
		// pair
		return 2
	case len(freqs) == 5 && freqs[0] == 1 && freqs[1] == 1 && freqs[2] == 1 && freqs[3] == 1 && freqs[4] == 1:
		// high card
		return 1
	default:
		panic("impossible hand")
	}
}

func wildcardHandType(cards string) int {
	h2c := make(map[string]struct{}, 1)
	h2c[cards] = struct{}{}
	for i := range cards {
		for h := range h2c {
			if h[i] == 'J' {
				for _, nc := range "23456789TQKA" {
					nh := h[:i] + string(nc) + h[i+1:]
					h2c[nh] = struct{}{}
				}
			}
		}
	}
	bestType := 0
	for h := range h2c {
		t := handType(makefreqs(h))
		if t > bestType {
			bestType = t
		}
	}
	return bestType
}

func run() error {
	d := data{
		hands: make([]hand, 0),
	}
	err := utils.OpenAndReadLines("input7.txt", func(s string) error {
		fields := strings.Fields(s)
		d.hands = append(d.hands, hand{
			cards:            fields[0],
			bid:              utils.MustAtoi(fields[1]),
			handType:         handType(makefreqs(fields[0])),
			wildcardHandType: wildcardHandType(fields[0]),
		})
		return nil
	})
	if err != nil {
		return err
	}
	sort.Slice(d.hands, d.lessThanP1)
	part1 := 0
	for i, h := range d.hands {
		part1 += (i + 1) * h.bid
	}
	fmt.Printf("Part 1: %d\n", part1)
	sort.Slice(d.hands, d.lessThanP2)
	part2 := 0
	for i, h := range d.hands {
		part2 += (i + 1) * h.bid
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
