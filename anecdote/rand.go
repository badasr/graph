// Copyright 2017 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package main

import (
	"fmt"

	"github.com/soniakeys/graph"
)

var gnpUSmall graph.Undirected
var gnpUSmallTag string

func GnpUSmall() (string, int, int) {
	const n = 1000
	gnpUSmallTag = fmt.Sprint("Gnp ", n, " nds")
	gnpUSmall = graph.GnpUndirected(n, .2, nil)
	return "Gnp undirected", n, gnpUSmall.ArcSize() / 2
}

var gnpULarge graph.Undirected
var gnpULargeTag string

func GnpULarge() (string, int, int) {
	const n = 2e4
	gnpULargeTag = fmt.Sprint("Gnp ", n, " nds")
	gnpULarge = graph.GnpUndirected(n, .105, nil)
	return "Gnp undirected", n, gnpULarge.ArcSize() / 2
}

var gnmUSmall graph.Undirected
var gnmUSmallTag string

func GnmUSmall() (string, int, int) {
	const n = 1000
	const m = 100e3
	gnmUSmallTag = fmt.Sprint("Gnm ", n, " nds")
	gnmUSmall = graph.GnmUndirected(n, m, nil)
	return "Gnm undirected", n, m
}

var gnmULarge graph.Undirected
var gnmULargeTag string

func GnmULarge() (string, int, int) {
	const n = 20e3
	const m = 20e6
	gnmULargeTag = fmt.Sprint("Gnm ", n, " nds")
	gnmULarge = graph.GnmUndirected(n, m, nil)
	return "Gnm undirected", n, m
}

func Gnm3Small() (string, int, int) {
	const n = 1000
	const m = 100e3
	graph.Gnm3Undirected(n, m, nil)
	return "Gnm3 undirected", n, m
}

func Gnm3Large() (string, int, int) {
	const n = 20e3
	const m = 20e6
	graph.Gnm3Undirected(n, m, nil)
	return "Gnm3 undirected", n, m
}

var gnpDSmall graph.Directed
var gnpDSmallTag string

func GnpDSmall() (string, int, int) {
	const n = 1000
	gnpDSmallTag = fmt.Sprint("Gnp ", n, " nds")
	gnpDSmall = graph.GnpDirected(n, .101, nil)
	return "Gnp directed", n, gnpDSmall.ArcSize()
}

var gnpDLarge graph.Directed
var gnpDLargeTag string

func GnpDLarge() (string, int, int) {
	const n = 2e4
	gnpDLargeTag = fmt.Sprint("Gnp ", n, " nds")
	gnpDLarge = graph.GnpDirected(n, .05, nil)
	return "Gnp directed", n, gnpDLarge.ArcSize()
}

var gnmDSmall graph.Directed
var gnmDSmallTag string

func GnmDSmall() (string, int, int) {
	const n = 1000
	const ma = 100e3
	gnmDSmallTag = fmt.Sprint("Gnm ", n, " nds")
	gnmDSmall = graph.GnmDirected(n, ma, nil)
	return "Gnm directed", n, ma
}

var gnmDLarge graph.Directed
var gnmDLargeTag string

func GnmDLarge() (string, int, int) {
	const n = 20e3
	const ma = 20e6
	gnmDLargeTag = fmt.Sprint("Gnm ", n, " nds")
	gnmDLarge = graph.GnmDirected(n, ma, nil)
	return "Gnm directed", n, ma
}

var eucSmall graph.Directed
var eucSmallTag string

func EucSmall() (string, int, int) {
	const n = 1e3
	const ma = 5e3
	var err error
	eucSmall, _, err = graph.Euclidean(n, ma, 1, 1, nil)
	if err != nil {
		return "nope", n, ma
	}
	eucSmallTag = "Euclidean "+h(n)+" nds"
	return "Euclidean", n, ma
}

var eucLarge graph.Directed
var eucLargeTag string

func EucLarge() (string, int, int) {
	const n = 1e6
	const ma = 5e6
	var err error
	eucLarge, _, err = graph.Euclidean(n, ma, 1, 1, nil)
	if err != nil {
		return "nope", n, ma
	}
	eucLargeTag = "Euclidean "+h(n)+" nds"
	return "Euclidean", n, ma
}

var kronDSmall graph.Directed
var kronDSmallTag string

func KronDSmall() (string, int, int) {
	kronDSmall, ma := graph.KroneckerDirected(11, 7, nil)
	kronDSmallTag = fmt.Sprint("Kron ", kronDSmall.Order(), "nds")
	return "Kronecker directed", kronDSmall.Order(), ma
}

var kronDLarge graph.Directed
var kronDLargeTag string

func KronDLarge() (string, int, int) {
	kronDLarge, ma := graph.KroneckerDirected(17, 21, nil)
	kronDLargeTag = fmt.Sprint("Kron ", kronDLarge.Order(), "nds")
	return "Kronecker directed", kronDLarge.Order(), ma
}

var kronUSmall graph.Undirected
var kronUSmallTag string

func KronUSmall() (string, int, int) {
	kronUSmall, m := graph.KroneckerUndirected(11, 7, nil)
	kronUSmallTag = fmt.Sprint("Kron ", kronUSmall.Order(), "nds")
	return "Kronecker undirected", kronUSmall.Order(), m
}

var kronULarge graph.Undirected
var kronULargeTag string

func KronULarge() (string, int, int) {
	kronULarge, m := graph.KroneckerUndirected(17, 21, nil)
	kronULargeTag = fmt.Sprint("Kron ", kronULarge.Order(), "nds")
	return "Kronecker undirected", kronULarge.Order(), m
}

var geoSmall graph.Undirected
var geoSmallTag string

func GeoSmall() (string, int, int) {
	const n = 1000
	var m int
	geoSmall, _, m = graph.Geometric(n, .1, nil)
	geoSmallTag = fmt.Sprint("Geom ", n, "nds")
	return "Geometric", n, m
}

var geoLarge graph.Undirected
var geoLargeTag string

func GeoLarge() (string, int, int) {
	const n = 3e4
	var m int
	geoLarge, _, m = graph.Geometric(n, .01, nil)
	geoLargeTag = fmt.Sprint("Geom ", n, "nds")
	return "Geometric", n, m
}

var chungLuSmall graph.Undirected
var chungLuSmallTag string

func ChungLuSmall() (string, int, int) {
	const n = 1e4
	w := make([]float64, n)
	for i := range w {
		w[i] = 5 + 10*float64(n-i)/float64(n)
	}
	chungLuSmall = graph.ChungLu(w, nil)
	chungLuSmallTag = "ChungLu " + h(n) + " nds"
	return "ChungLu", n, chungLuSmall.ArcSize() / 2
}

var chungLuLarge graph.Undirected
var chungLuLargeTag string

func ChungLuLarge() (string, int, int) {
	const n = 2e5
	w := make([]float64, n)
	for i := range w {
		w[i] = 2 + 50*n/float64(i+1)
	}
	chungLuLarge = graph.ChungLu(w, nil)
	chungLuLargeTag = "ChungLu " + h(n) + " nds"
	return "ChungLu", n, chungLuLarge.ArcSize() / 2
}
