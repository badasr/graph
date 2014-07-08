package graph

import (
	"math/big"
	"math/rand"
	"time"
)

// scale is log(2) node extent.
// arcFactor is rato of arcs generated to node extent.
func NewDir(scale uint, edgeFactor float64) (g AdjacencyList, m int) {
	return krongen(scale, edgeFactor, true)
}

func NewUnDir(scale uint, edgeFactor float64) (g AdjacencyList, m int) {
	return krongen(scale, edgeFactor, false)
}

func krongen(scale uint, edgeFactor float64, dir bool) (g AdjacencyList, m int) {
	rand.Seed(time.Now().Unix())
	N := 1 << scale                      // node extent
	M := int(edgeFactor*float64(N) + .5) // number of arcs to generate
	a, b, c := 0.57, 0.19, 0.19          // initiator probabilities
	ab := a + b
	cNorm := c / (1 - ab)
	aNorm := a / ab
	ij := make([][2]int, M)
	var bm big.Int
	var nNodes int
	for k := range ij {
		var i, j int
		for b := 1; b < N; b <<= 1 {
			if rand.Float64() > ab {
				i |= b
				if rand.Float64() > cNorm {
					j |= b
				}
			} else if rand.Float64() > aNorm {
				j |= b
			}
		}
		if bm.Bit(i) == 0 {
			bm.SetBit(&bm, i, 1)
			nNodes++
		}
		if bm.Bit(j) == 0 {
			bm.SetBit(&bm, j, 1)
			nNodes++
		}
		ij[k] = [2]int{i, j}
	}
	p := rand.Perm(nNodes)
	px := 0
	r := make([]int, M)
	for i := range r {
		if bm.Bit(i) == 1 {
			r[i] = p[px]
			px++
		}
	}
	g = make(AdjacencyList, nNodes)
ij:
	for _, e := range ij {
		if e[0] == e[1] {
			continue
		}
		ri, rj := r[e[0]], r[e[1]]
		for _, nb := range g[ri] {
			if nb == rj {
				continue ij
			}
		}
		g[ri] = append(g[ri], rj)
		if !dir {
			g[rj] = append(g[rj], ri)
		}
		m++
	}
	return
}
