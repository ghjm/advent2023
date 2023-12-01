package main

import (
	"fmt"
	"github.com/ghjm/advent2023/pkg/utils"
	"os"
	"unicode"
)

func setFirstLast(first *int, last *int, value int) {
	if *first == 0 {
		*first = value
	}
	*last = value
}

func run() error {
	sumP1 := 0
	sumP2 := 0
	digits := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	err := utils.OpenAndReadLines("input1.txt", func(s string) error {
		var firstP1, lastP1, firstP2, lastP2 int
		for i, c := range s {
			if unicode.IsDigit(c) {
				ci := int(c - '0')
				setFirstLast(&firstP1, &lastP1, ci)
				setFirstLast(&firstP2, &lastP2, ci)
			} else {
				for j, d := range digits {
					if (len(s) >= i+len(d)) && (s[i:i+len(d)] == d) {
						setFirstLast(&firstP2, &lastP2, j)
						break
					}
				}
			}
		}
		sumP1 += firstP1*10 + lastP1
		sumP2 += firstP2*10 + lastP2
		return nil
	})
	if err != nil {
		return err
	}
	fmt.Printf("Part 1: %d\nPart 2: %d\n", sumP1, sumP2)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
