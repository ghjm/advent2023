package main

import (
	"fmt"
	"os"
	"strings"

	utils "github.com/ghjm/advent_utils"
)

type record struct {
	condition string
	groups    []byte
}

type memoKey struct {
	condition string
	groups    string
}

type data struct {
	records []record
	memo    map[memoKey]int
}

// credit to shemetz for the idea of this solution
func (d *data) numArrangements(condition string, groups []byte) int {
	key := memoKey{condition, string(groups)}
	v, ok := d.memo[key]
	if ok {
		return v
	}
	var result int
	if len(condition) == 0 {
		if len(groups) == 0 {
			result = 1
		} else {
			result = 0
		}
	} else {
		switch condition[0] {
		case '.':
			result = d.numArrangements(condition[1:], groups)
		case '?':
			result = d.numArrangements("."+condition[1:], groups) + d.numArrangements("#"+condition[1:], groups)
		case '#':
			if len(groups) == 0 {
				result = 0
			} else if len(condition) < int(groups[0]) {
				result = 0
			} else if strings.Contains(condition[:groups[0]], ".") {
				result = 0
			} else {
				if len(groups) > 1 {
					if len(condition) < int(groups[0]+1) || condition[groups[0]] == '#' {
						result = 0
					} else {
						result = d.numArrangements(condition[groups[0]+1:], groups[1:])
					}
				} else {
					result = d.numArrangements(condition[groups[0]:], groups[1:])
				}
			}
		}
	}
	d.memo[key] = result
	return result
}

func run() error {
	d := data{}
	err := utils.OpenAndReadLines("input12.txt", func(s string) error {
		sf := strings.Fields(s)
		r := record{
			condition: sf[0],
		}
		for _, g := range strings.Split(sf[1], ",") {
			gi := utils.MustAtoi(g)
			if gi > 255 {
				panic("not byte representable")
			}
			r.groups = append(r.groups, byte(gi))
		}
		d.records = append(d.records, r)
		return nil
	})
	if err != nil {
		return err
	}

	for i := 1; i <= 2; i++ {
		sum := 0
		for _, r := range d.records {
			d.memo = make(map[memoKey]int)
			var rc string
			var rg []byte
			if i == 1 {
				rc = r.condition
				rg = r.groups
			} else {
				for j := 0; j < 5; j++ {
					rc = rc + r.condition
					if j < 4 {
						rc = rc + "?"
					}
					rg = append(rg, r.groups...)
				}
			}
			n := d.numArrangements(rc, rg)
			sum += n
		}
		fmt.Printf("Part %d: %d\n", i, sum)
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
