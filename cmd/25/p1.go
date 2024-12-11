package main

import (
	"fmt"
	"os"
	"strings"

	utils "github.com/ghjm/advent_utils"
	"github.com/twmb/algoimpl/go/graph"
)

type data struct {
	g     *graph.Graph
	nodes map[string]graph.Node
}

func run() error {
	d := data{
		g:     graph.New(graph.Undirected),
		nodes: make(map[string]graph.Node),
	}
	err := utils.OpenAndReadLines("input25.txt", func(s string) error {
		v1 := strings.Split(s, ": ")
		if len(v1) != 2 {
			return fmt.Errorf("invalid input")
		}
		_, ok := d.nodes[v1[0]]
		if !ok {
			d.nodes[v1[0]] = d.g.MakeNode()
		}
		for _, v2 := range strings.Split(v1[1], " ") {
			_, ok = d.nodes[v2]
			if !ok {
				d.nodes[v2] = d.g.MakeNode()
			}
			err := d.g.MakeEdge(d.nodes[v1[0]], d.nodes[v2])
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	var cut []graph.Edge
	for len(cut) != 3 {
		cut = d.g.RandMinimumCut(12, 4)
	}
	for _, edge := range cut {
		d.g.RemoveEdge(edge.Start, edge.End)
	}
	scc := d.g.StronglyConnectedComponents()
	if len(scc) != 2 {
		return fmt.Errorf("expected 2 components, got %d", len(scc))
	}
	fmt.Printf("Part 1: %d\n", len(scc[0])*len(scc[1]))
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
