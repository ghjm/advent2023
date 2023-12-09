package main

import (
	"fmt"
	"github.com/ghjm/advent2023/pkg/utils"
	"os"
	"strings"
)

type data struct {
	sequences [][]int
}

func getDifferences(seq []int) []int {
	result := make([]int, 0)
	if len(seq) < 2 {
		return result
	}
	for i := 1; i < len(seq); i++ {
		result = append(result, seq[i]-seq[i-1])
	}
	return result
}

func run() error {
	d := data{
		sequences: make([][]int, 0),
	}
	err := utils.OpenAndReadLines("input9.txt", func(s string) error {
		if s == "" {
			return nil
		}
		newSeq := make([]int, 0)
		for _, f := range strings.Fields(s) {
			newSeq = append(newSeq, utils.MustAtoi(f))
		}
		d.sequences = append(d.sequences, newSeq)
		return nil
	})
	if err != nil {
		return err
	}
	part1 := 0
	part2 := 0
	for _, seq := range d.sequences {
		diffList := make([][]int, 1)
		diffList[0] = seq
		for {
			diffList = append(diffList, getDifferences(diffList[len(diffList)-1]))
			allZeroes := true
			for _, v := range diffList[len(diffList)-1] {
				if v != 0 {
					allZeroes = false
					break
				}
			}
			if allZeroes {
				break
			}
		}
		diffList[len(diffList)-1] = append(diffList[len(diffList)-1], 0)
		for i := len(diffList) - 2; i >= 0; i-- {
			diffList[i] = append(diffList[i], diffList[i][len(diffList[i])-1]+diffList[i+1][len(diffList[i+1])-1])
		}
		part1 += diffList[0][len(diffList[0])-1]
		diffList[len(diffList)-1] = append([]int{0}, diffList[len(diffList)-1]...)
		for i := len(diffList) - 2; i >= 0; i-- {
			diffList[i] = append([]int{diffList[i][0] - diffList[i+1][0]}, diffList[i]...)
		}
		part2 += diffList[0][0]
	}
	fmt.Printf("Part 1: %d\n", part1)
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
