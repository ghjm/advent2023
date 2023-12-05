package main

import (
	"fmt"
	"github.com/ghjm/advent2023/pkg/utils"
	"math"
	"os"
	"regexp"
	"strings"
)

type almanacEntry struct {
	source      int64
	destination int64
	length      int64
}

type almanac struct {
	from    string
	to      string
	entries []*almanacEntry
}

type idRange struct {
	start int64
	end   int64
}

type idSet struct {
	ranges []idRange
}

type data struct {
	seeds    []int64
	almanacs []*almanac
}

func (a *almanac) lookup(id int64) int64 {
	for _, ae := range a.entries {
		if id >= ae.source && id < ae.source+ae.length {
			return ae.destination + id - ae.source
		}
	}
	return id
}

func (d *data) findLocation(idKind string, idValue int64) int64 {
	if idKind == "location" {
		return idValue
	}
	for _, a := range d.almanacs {
		if a.from == idKind {
			return d.findLocation(a.to, a.lookup(idValue))
		}
	}
	panic(fmt.Sprintf("unresolvable item kind %s", idKind))
}

func (d *data) getLocationRanges(idKind string, ids idSet) idSet {
	if idKind == "location" {
		return ids
	}
	for _, a := range d.almanacs {
		if a.from == idKind {
			result := idSet{
				ranges: make([]idRange, 0),
			}
			for _, r := range ids.ranges {
				ts := a.translateRange(r)
				for _, tr := range ts.ranges {
					result.ranges = append(result.ranges, tr)
				}
			}
			return d.getLocationRanges(a.to, result)
		}
	}
	panic(fmt.Sprintf("unresolvable item kind %s", idKind))
}

func (a *almanac) entryContaining(id int64) *almanacEntry {
	for _, ae := range a.entries {
		if ae.source <= id && id < ae.source+ae.length {
			return ae
		}
	}
	return nil
}

func (a *almanac) entryAfter(id int64) *almanacEntry {
	var min int64
	var result *almanacEntry
	min = math.MaxInt64
	for _, ae := range a.entries {
		if ae.source >= id && ae.source < min {
			min = ae.source
			result = ae
		}
	}
	return result
}

func (a *almanac) translateRange(r idRange) idSet {
	result := idSet{
		ranges: make([]idRange, 0),
	}
	var curPos int64
	curPos = r.start

	// Check if we are starting already in a range
	ae := a.entryContaining(curPos)
	if ae != nil {
		rangeEnd := ae.source + ae.length - 1
		if rangeEnd > r.end {
			rangeEnd = r.end
		}
		result.ranges = append(result.ranges, idRange{
			start: curPos + ae.destination - ae.source,
			end:   rangeEnd + ae.destination - ae.source,
		})
		curPos = rangeEnd + 1
	}
	for curPos < r.end {
		ae = a.entryAfter(curPos)
		if ae == nil || ae.source >= r.end {
			// Remaining range after all translations are finished
			result.ranges = append(result.ranges, idRange{
				start: curPos,
				end:   r.end,
			})
			break
		}

		// range from curPos to the beginning of the translation
		if ae.source > curPos {
			result.ranges = append(result.ranges, idRange{
				start: curPos,
				end:   ae.source - 1,
			})
			curPos = ae.source
		}

		// range of the translation, possibly shortened if our source range is shorter
		rangeEnd := ae.source + ae.length - 1
		if rangeEnd > r.end {
			rangeEnd = r.end
		}
		result.ranges = append(result.ranges, idRange{
			start: ae.destination,
			end:   rangeEnd + ae.destination - ae.source,
		})
		curPos = rangeEnd + 1

	}
	return result
}

func run() error {
	seedsRe := regexp.MustCompile(`^seeds: (.*)$`)
	mapStartRe := regexp.MustCompile(`^(.+)-to-(.*) map:$`)
	mapEntryRe := regexp.MustCompile(`^(\d+) (\d+) (\d+)$`)
	d := data{
		seeds:    make([]int64, 0),
		almanacs: make([]*almanac, 0),
	}
	var curAlmanac *almanac
	err := utils.OpenAndReadLines("input5.txt", func(s string) error {
		if s == "" {
			return nil
		}
		m := seedsRe.FindStringSubmatch(s)
		if m != nil {
			for _, seed := range strings.Split(m[1], " ") {
				d.seeds = append(d.seeds, utils.MustAtoi64(seed))
			}
			return nil
		}
		m = mapStartRe.FindStringSubmatch(s)
		if m != nil {
			curAlmanac = &almanac{
				from:    m[1],
				to:      m[2],
				entries: make([]*almanacEntry, 0),
			}
			d.almanacs = append(d.almanacs, curAlmanac)
			return nil
		}
		m = mapEntryRe.FindStringSubmatch(s)
		if m != nil {
			ae := &almanacEntry{
				source:      utils.MustAtoi64(m[2]),
				destination: utils.MustAtoi64(m[1]),
				length:      utils.MustAtoi64(m[3]),
			}
			curAlmanac.entries = append(curAlmanac.entries, ae)
			return nil
		}
		return fmt.Errorf("unmatched line: %s", s)
	})
	if err != nil {
		return err
	}
	var part1 int64
	part1 = math.MaxInt64
	for _, s := range d.seeds {
		loc := d.findLocation("seed", s)
		if loc < part1 {
			part1 = loc
		}
	}
	fmt.Printf("Part 1: %d\n", part1)

	seedSet := idSet{
		ranges: make([]idRange, 0),
	}
	for i := 0; i < len(d.seeds); i += 2 {
		seedSet.ranges = append(seedSet.ranges, idRange{
			start: d.seeds[i],
			end:   d.seeds[i] + d.seeds[i+1] - 1,
		})
	}
	locSet := d.getLocationRanges("seed", seedSet)
	var min int64
	min = math.MaxInt64
	for _, r := range locSet.ranges {
		if r.start < min {
			min = r.start
		}
	}
	fmt.Printf("Part 2: %d\n", min)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
