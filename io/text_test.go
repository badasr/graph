// Copyright 2018 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package io_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/io"
)

/* consider:
should read with arcDir != All reconstruct undirected graphs?
graph.Directed.Undirected will do it, but it would be more efficient to do it
here.
*/

type allErr struct{}

func (allErr) Read([]byte) (int, error) {
	return 0, errors.New("always error")
}

func TestReadALDense(t *testing.T) {
	// test normal operation
	want := graph.AdjacencyList{
		0: {2, 1, 1},
		2: {1},
	}
	r := bytes.NewBufferString(`2 1 1

1`)
	got, _, _, err := io.Text{Format: io.Dense}.ReadAdjacencyList(r)
	if err != nil {
		t.Fatal("readALDense: ", err)
	}
	if !got.Equal(want) {
		for fr, to := range got {
			t.Log(fr, " ", to)
		}
		t.Fail()
	}
	// test MapNames invalid
	tx := io.Text{Format: io.Dense, MapNames: true}
	if _, _, _, err := tx.ReadAdjacencyList(nil); err == nil {
		t.Fatal("readALDense allowed MapNames")
	}
	// test a read error
	_, _, _, err = io.Text{Format: io.Dense}.ReadAdjacencyList(allErr{})
	if err == nil {
		t.Fatal("readALDense allowed read error")
	}
}

func TestReadALArcs(t *testing.T) {
	want := graph.AdjacencyList{
		0: {2, 1, 1},
		2: {1},
	}
	r := bytes.NewBufferString(`
0 2
0 1
0 1
2 1`)
	got, _, _, err := io.Text{Format: io.Arcs}.ReadAdjacencyList(r)
	if err != nil {
		t.Fatal("readALArcs: ", err)
	}
	if !got.Equal(want) {
		for fr, to := range got {
			t.Log(fr, " ", to)
		}
		t.Fail()
	}
}

func TestBadFormat(t *testing.T) {
	_, _, _, err := io.Text{Format: -1}.ReadAdjacencyList(nil)
	if err == nil {
		t.Fatal("ReadAdacencyList no err from bad Format")
	}
}

func TestSep(t *testing.T) {
	// test a base < 10
	want := graph.AdjacencyList{
		0: {2, 1, 1},
		2: {1},
	}
	r := bytes.NewBufferString(`
0:829181
2:91`)
	got, _, _, err := io.Text{Base: 8}.ReadAdjacencyList(r)
	if err != nil {
		t.Fatal("sep: ", err)
	}
	if !got.Equal(want) {
		for fr, to := range got {
			t.Log(fr, " ", to)
		}
		t.Fail()
	}
	// test a base > 10
	const zz = 36*36 - 1
	want = graph.AdjacencyList{
		zz: {0},
	}
	r = bytes.NewBufferString(`zz: 0`)
	got, _, _, err = io.Text{Base: 36}.ReadAdjacencyList(r)
	if err != nil {
		t.Fatal("sep: ", err)
	}
	if !got.Equal(want) {
		for fr, to := range got {
			t.Log(fr, " ", to)
		}
		t.Fail()
	}
}

func TestReadSplitInts(t *testing.T) {
	// a couple of cases for test coverage
	want := graph.AdjacencyList{
		0: {2, 1, 1},
		2: {1},
	}
	r := bytes.NewBufferString(` 
:
0: 2 1 1
2: 1`)
	got, _, _, err := io.Text{FrDelim: ":"}.ReadAdjacencyList(r)
	if err != nil {
		t.Fatal("readSplitInts test: ", err)
	}
	if !got.Equal(want) {
		for fr, to := range got {
			t.Log(fr, " ", to)
		}
		t.Fail()
	}
}

/* removed examples:

func ExampleText_base36() {
	//      100000
	//       / \\
	// 100001-->100002
	g := graph.AdjacencyList{
		100000: {100001, 100002, 100002},
		100001: {100002},
		100002: {},
	}
	b := &bytes.Buffer{}
	t := io.Text{Base: 36}
	t.WriteAdjacencyList(g, b)
	fmt.Println(b)
	rt, _, _, _ := t.ReadAdjacencyList(b)
	for fr, to := range rt {
		if len(to) > 0 {
			fmt.Println(fr, to)
		}
	}
	// Output:
	// 255s: 255t 255u 255u
	// 255t: 255u
	// 255u:
	//
	// 100000 [100001 100002 100002]
	// 100001 [100002]
}

func ExampleText_ReadAdjacencyList_dense() {
	r := bytes.NewBufferString(`2 1 1

1`)
	g, _, _, err := io.Text{Format: io.Dense}.ReadAdjacencyList(r)
	for n, to := range g {
		fmt.Println(n, to)
	}
	fmt.Println("err: ", err)
	// Output:
	// 0 [2 1 1]
	// 1 []
	// 2 [1]
	// err:  <nil>
}

func ExampleText_ReadAdjacencyList() {
	r := bytes.NewBufferString(`
0: 2 1 1
2: 1`)
	g, _, _, err := io.Text{}.ReadAdjacencyList(r)
	for n, to := range g {
		fmt.Println(n, to)
	}
	fmt.Println("err: ", err)
	// Output:
	// 0 [2 1 1]
	// 1 []
	// 2 [1]
	// err:  <nil>
}

func ExampleText_ReadAdjacencyList_arcs() {
	r := bytes.NewBufferString(`
0 2
0 1
0 1 // parallel
2 1
`)
	t := io.Text{Format: io.Arcs, Comment: "//"}
	g, _, _, err := t.ReadAdjacencyList(r)
	for n, to := range g {
		fmt.Println(n, to)
	}
	fmt.Println("err: ", err)
	// Output:
	// 0 [2 1 1]
	// 1 []
	// 2 [1]
	// err:  <nil>
}

func ExampleText_ReadAdjacencyList_arcNames() {
	r := bytes.NewBufferString(`
a b
a b // parallel
a c
c b
`)
	t := io.Text{Format: io.Arcs, ReadNodeNames: true, Comment: "//"}
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
	fmt.Println("err: ", err)
	// Output:
	// names:
	// 0 a
	// 1 b
	// 2 c
	// graph:
	// 0 [1 1 2]
	// 1 []
	// 2 [1]
	// map[a:0 b:1 c:2 ]
	// err:  <nil>
}

func ExampleText_WriteAdjacencyList() {
	//   0
	//  / \\
	// 2-->1
	g := graph.AdjacencyList{
		0: {2, 1, 1},
		2: {1},
	}
	n, err := io.Text{}.WriteAdjacencyList(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// 0: 2 1 1
	// 2: 1
	// bytes: 14, err: <nil>
}

func ExampleText_WriteAdjacencyList_names() {
	//   a   d
	//  / \   \
	// b   c   e
	g := graph.AdjacencyList{
		0: {1, 2},
		3: {4},
		4: {},
	}
	names := []string{"a", "b", "c", "d", "e"}
	t := io.Text{WriteNodeName: func(n graph.NI) string { return names[n] }}
	n, err := t.WriteAdjacencyList(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// a: b c
	// d: e
	// bytes: 12, err: <nil>
}

func ExampleText_WriteAdjacencyList_arcs() {
	//   0
	//  / \\
	// 2-->1
	g := graph.AdjacencyList{
		0: {2, 1, 1},
		2: {1},
	}
	n, err := io.Text{Format: io.Arcs}.WriteAdjacencyList(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// 0 2
	// 0 1
	// 0 1
	// 2 1
	// bytes: 16, err: <nil>
}

func ExampleText_WriteAdjacencyList_arcNames() {
	//   a
	//  / \\
	// c-->b
	g := graph.AdjacencyList{
		0: {2, 1, 1},
		2: {1},
	}
	names := []string{"a", "b", "c"}
	t := io.Text{
		Format:        io.Arcs,
		WriteNodeName: func(n graph.NI) string { return names[n] },
	}
	n, err := t.WriteAdjacencyList(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// a c
	// a b
	// a b
	// c b
	// bytes: 16, err: <nil>
}

func ExampleText_WriteAdjacencyList_undirectedDense() {
	//   0
	//  / \\
	// 1---2--\
	//      \-/
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(0, 2)
	g.AddEdge(2, 2)
	t := io.Text{Format: io.Dense, WriteArcs: io.Upper}
	n, err := t.WriteAdjacencyList(g.AdjacencyList, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// 1 2 2
	//
	// 2
	// bytes: 9, err: <nil>
}

func ExampleText_WriteAdjacencyList_undirectedNames() {
	//   a
	//  / \\
	// b---c--\
	//      \-/
	var g graph.Undirected
	names := []string{"a", "b", "c"}
	ni := map[string]graph.NI{}
	for i, s := range names {
		ni[s] = graph.NI(i)
	}
	g.AddEdge(ni["a"], ni["b"])
	g.AddEdge(ni["a"], ni["c"])
	g.AddEdge(ni["a"], ni["c"])
	g.AddEdge(ni["c"], ni["c"])
	t := io.Text{
		WriteArcs:     io.Upper,
		WriteNodeName: func(n graph.NI) string { return names[n] },
	}
	n, err := t.WriteAdjacencyList(g.AdjacencyList, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// a: b c c
	// c: c
	// bytes: 14, err: <nil>
}

func ExampleText_mapNamesTab() {
	//   a   d
	//  / \   \
	// b   c   e
	r := bytes.NewBufferString(`
a	b c  # source target target
d e
`)
	// For reading, default blank delimiter fields enable
	// delimiting by whitespace.
	t := io.Text{MapNames: true, FrDelim: "\t", Comment: "#"}
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
	// map[a:0 b:1 c:2 d:3 e:4]
	// err: <nil>
}

/* commented out:  these examples show dense format, all that is currently
implemented, but default for zero value Text should be sparse.
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
*/
