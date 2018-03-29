// Copyright 2018 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package io_test

import (
	"bytes"
	"fmt"
	"os"
	"reflect"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/io"
)

func ExampleArcDir() {
	//   0
	//  / \\
	// 1   2--\
	//      \-/
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(0, 2)
	g.AddEdge(2, 2)

	fmt.Println("Default WriteArcs All:  writes all arcs:")
	t := io.Text{}
	n, err := t.WriteAdjacencyList(g.AdjacencyList, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n\n", n, err)

	fmt.Println("Upper writes only arcs representing undirected edges:")
	t.WriteArcs = io.Upper
	n, err = t.WriteAdjacencyList(g.AdjacencyList, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// Default WriteArcs All:  writes all arcs:
	// 0: 1 2 2
	// 1: 0
	// 2: 0 0 2
	// bytes: 23, err: <nil>
	//
	// Upper writes only arcs representing undirected edges:
	// 0: 1 2 2
	// 2: 2
	//bytes: 14, err: <nil>
}

func ExampleFormat() {
	//   0
	//  / \\
	// 2-->3
	g := graph.AdjacencyList{
		0: {2, 3, 3},
		2: {3},
		3: {},
	}
	fmt.Println("Default Format Sparse:")
	var t io.Text
	n, err := t.WriteAdjacencyList(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n\n", n, err)

	fmt.Println("Format Dense:")
	t.Format = io.Dense
	n, err = t.WriteAdjacencyList(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n\n", n, err)

	fmt.Println("Format Arcs:")
	t.Format = io.Arcs
	n, err = t.WriteAdjacencyList(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n\n", n, err)
	// Output:
	// Default Format Sparse:
	// 0: 2 3 3
	// 2: 3
	// 3:
	// bytes: 17, err: <nil>
	//
	// Format Dense:
	// 2 3 3
	//
	// 3
	//
	// bytes: 10, err: <nil>
	//
	// Format Arcs:
	// 0 2
	// 0 3
	// 0 3
	// 2 3
	// bytes: 16, err: <nil>
}

func ExampleText_readNodeNames() {
	//   a   d
	//  / \   \
	// b   c   e
	r := bytes.NewBufferString(`
a b c  # source target target
d e
`)
	// For reading, default blank delimiter fields enable
	// delimiting by whitespace.
	t := io.Text{ReadNodeNames: true, Comment: "#"}
	g, names, m, err := t.ReadAdjacencyList(r)

	fmt.Println("names:")
	for i, n := range names {
		fmt.Println(i, n)
	}
	fmt.Println("graph:")
	for n, to := range g {
		fmt.Println(n, to)
	}
	fmt.Println(graph.OrderMap(m))
	fmt.Println("err:", err)
	// Output:
	// names:
	// 0 a
	// 1 b
	// 2 c
	// 3 d
	// 4 e
	// graph:
	// 0 [1 2]
	// 1 []
	// 2 []
	// 3 [4]
	// 4 []
	// map[a:0 b:1 c:2 d:3 e:4 ]
	// err: <nil>
}

func ExampleText_zeroValue() {
	//   0
	//  / \\
	// 2-->3
	g := graph.AdjacencyList{
		0: {2, 3, 3},
		2: {3},
		3: {},
	}
	var t io.Text // zero value
	var b bytes.Buffer
	t.WriteAdjacencyList(g, &b)
	fmt.Println("Written:")
	fmt.Print(b.String())
	// demonstrate round trip
	rt, _, _, _ := t.ReadAdjacencyList(&b)
	fmt.Println("Round trip:", reflect.DeepEqual(g, rt))
	// Output:
	// Written:
	// 0: 2 3 3
	// 2: 3
	// 3:
	// Round trip: true
}

func ExampleText_ReadLabeledAdjacencyList() {
	r := bytes.NewBufferString(`2 101 1 102 1 102

1 103`)
	g, err := io.Text{}.ReadLabeledAdjacencyList(r)
	for n, to := range g {
		fmt.Println(n, to)
	}
	fmt.Println("err: ", err)
	// Output:
	// 0 [{2 101} {1 102} {1 102}]
	// 1 []
	// 2 [{1 103}]
	// err:  <nil>
}

func ExampleText_WriteLabeledAdjacencyList() {
	//        0
	// (101) / \\ (102)
	//      2-->1
	//      (103)
	g := graph.LabeledAdjacencyList{
		0: {{2, 101}, {1, 102}, {1, 102}},
		2: {{1, 103}},
	}
	n, err := io.Text{}.WriteLabeledAdjacencyList(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// 2 101 1 102 1 102
	//
	// 1 103
	// bytes: 25, err: <nil>
}