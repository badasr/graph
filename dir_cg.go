// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

import (
	"github.com/soniakeys/bits"
)

// dir_RO.go is code generated from dir_cg.go by directives in graph.go.
// Editing dir_cg.go is okay.  It is the code generation source.
// DO NOT EDIT dir_RO.go.
// The RO means read only and it is upper case RO to slow you down a bit
// in case you start to edit the file.

// Balanced returns true if for every node in g, in-degree equals out-degree.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledDirected) Balanced() bool {
	for n, in := range g.InDegree() {
		if in != len(g.LabeledAdjacencyList[n]) {
			return false
		}
	}
	return true
}

// Copy makes a deep copy of g.
// Copy also computes the arc size ma, the number of arcs.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledDirected) Copy() (c LabeledDirected, ma int) {
	l, s := g.LabeledAdjacencyList.Copy()
	return LabeledDirected{l}, s
}

// Cyclic determines if g contains a cycle, a non-empty path from a node
// back to itself.
//
// Cyclic returns true if g contains at least one cycle.  It also returns
// an example of an arc involved in a cycle.
// Cyclic returns false if g is acyclic.
//
// Also see Topological, which detects cycles.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledDirected) Cyclic() (cyclic bool, fr NI, to Half) {
	a := g.LabeledAdjacencyList
	fr, to.To = -1, -1
	temp := bits.New(len(a))
	perm := bits.New(len(a))
	var df func(int)
	df = func(n int) {
		switch {
		case temp.Bit(n) == 1:
			cyclic = true
			return
		case perm.Bit(n) == 1:
			return
		}
		temp.SetBit(n, 1)
		for _, nb := range a[n] {
			df(int(nb.To))
			if cyclic {
				if fr < 0 {
					fr, to = NI(n), nb
				}
				return
			}
		}
		temp.SetBit(n, 0)
		perm.SetBit(n, 1)
	}
	for n := range a {
		if perm.Bit(n) == 1 {
			continue
		}
		if df(n); cyclic { // short circuit as soon as a cycle is found
			break
		}
	}
	return
}

// Dominators computes the immediate dominator for each node reachable from
// start.
//
// The slice returned as Dominators.Immediate will have the length of
// g.AdjacencyList.  Nodes without a path to end will have a value of -1.
//
// See also the method Doms.  Internally Dominators must construct the
// transpose of g and also compute a postordering of a spanning tree of the
// subgraph reachable from start.  If you happen to have either of these
// computed anyway, it can be more efficient to call Doms directly.
func (g LabeledDirected) Dominators(start NI) Dominators {
	a := g.LabeledAdjacencyList
	l := len(a)
	// ExampleDoms shows traditional depth-first postorder, but it works to
	// generate a reverse preorder.  Also breadth-first works instead of
	// depth-first and may allow Doms to run a little faster by presenting
	// a shallower tree.
	post := make([]NI, l)
	a.breadthFirst(start, nil, nil, func(n NI) bool {
		l--
		post[l] = n
		return true
	})
	tr, _ := g.Transpose()
	return g.Doms(tr, post[l:])
}

// Doms computes either immediate dominators or postdominators.
//
// The slice returned as Dominators.Immediate will have the length of
// g.AdjacencyList.  Nodes without a path to end will have a value of -1.
//
// But see also the simpler methods Dominators and PostDominators.
//
// Doms requires argument tr to be the transpose graph of receiver g,
// and requres argument post to be a post ordering of receiver g.  More
// specifically a post ordering of a spanning tree of the subgraph reachable
// from some start node in g.  The start node will always be the last node in
// this postordering so it does not need to passed as a separate argument.
//
// Doms can be used to construct either dominators or postdominators.
// To construct dominators on a graph f, generate a postordering p on f
// and call f.Doms(f.Transpose(), p).  To construct postdominators, generate
// the transpose t first, then a postordering p on t (not f), and call
// t.Doms(f, p).
//
// Caution:  The argument tr is retained in the returned Dominators object
// and is used by the method Dominators.Frontier.  It is not deep-copied
// so it is invalid to call Doms, modify the tr graph, and then call Frontier.
func (g LabeledDirected) Doms(tr LabeledDirected, post []NI) Dominators {
	a := g.LabeledAdjacencyList
	dom := make([]NI, len(a))
	pi := make([]int, len(a))
	for i, n := range post {
		pi[n] = i
	}
	intersect := func(b1, b2 NI) NI {
		for b1 != b2 {
			for pi[b1] < pi[b2] {
				b1 = dom[b1]
			}
			for pi[b2] < pi[b1] {
				b2 = dom[b2]
			}
		}
		return b1
	}
	for n := range dom {
		dom[n] = -1
	}
	start := post[len(post)-1]
	dom[start] = start
	for changed := false; ; changed = false {
		for i := len(post) - 2; i >= 0; i-- {
			b := post[i]
			var im NI
			fr := tr.LabeledAdjacencyList[b]
			var j int
			var fp Half
			for j, fp = range fr {
				if dom[fp.To] >= 0 {
					im = fp.To
					break
				}
			}
			for _, p := range fr[j:] {
				if dom[p.To] >= 0 {
					im = intersect(im, p.To)
				}
			}
			if dom[b] != im {
				dom[b] = im
				changed = true
			}
		}
		if !changed {
			return Dominators{dom, tr}
		}
	}
}

// PostDominators computes the immediate postdominator for each node that can
// reach node end.
//
// The slice returned as Dominators.Immediate will have the length of
// g.AdjacencyList.  Nodes without a path to end will have a value of -1.
//
// See also the method Doms.  Internally Dominators must construct the
// transpose of g and also compute a postordering of a spanning tree of the
// subgraph of the transpose reachable from end.  If you happen to have either
// of these computed anyway, it can be more efficient to call Doms directly.
//
// See the method Doms anyway for the caution note.  PostDominators calls
// Doms internally, passing receiver g as Doms argument tr.  The caution means
// that it is invalid to call PostDominators, modify the graph g, then call
// Frontier.
func (g LabeledDirected) PostDominators(end NI) Dominators {
	tr, _ := g.Transpose()
	a := tr.LabeledAdjacencyList
	l := len(a)
	post := make([]NI, l)
	a.breadthFirst(end, nil, nil, func(n NI) bool {
		l--
		post[l] = n
		return true
	})
	return tr.Doms(g, post[l:])
}

// called from Dominators.Frontier via interface
func (from LabeledDirected) domFrontiers(d Dominators) DominanceFrontiers {
	im := d.Immediate
	f := make(DominanceFrontiers, len(im))
	for i := range f {
		if im[i] >= 0 {
			f[i] = map[NI]struct{}{}
		}
	}
	for b, fr := range from.LabeledAdjacencyList {
		if len(fr) < 2 {
			continue
		}
		imb := im[b]
		for _, p := range fr {
			for runner := p.To; runner != imb; runner = im[runner] {
				f[runner][NI(b)] = struct{}{}
			}
		}
	}
	return f
}

// FromList transposes a labeled graph into a FromList.
//
// Receiver g should be connected as a tree or forest.  Specifically no node
// can have multiple incoming arcs.  If any node n in g has multiple incoming
// arcs, the method returns (nil, n) where n is a node with multiple
// incoming arcs.
//
// Otherwise (normally) the method populates the From members in a
// FromList.Path and returns the FromList and -1.
//
// Other members of the FromList are left as zero values.
// Use FromList.RecalcLen and FromList.RecalcLeaves as needed.
//
// Unusual cases are parallel arcs and loops.  A parallel arc represents
// a case of multiple arcs going to some node and so will lead to a (nil, n)
// return, even though a graph might be considered a multigraph tree.
// A single loop on a node that would otherwise be a root node, though,
// is not a case of multiple incoming arcs and so does not force a (nil, n)
// result.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledDirected) FromList() (*FromList, NI) {
	paths := make([]PathEnd, g.Order())
	for i := range paths {
		paths[i].From = -1
	}
	for fr, to := range g.LabeledAdjacencyList {
		for _, to := range to {
			if paths[to.To].From >= 0 {
				return nil, to.To
			}
			paths[to.To].From = NI(fr)
		}
	}
	return &FromList{Paths: paths}, -1
}

// InDegree computes the in-degree of each node in g
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledDirected) InDegree() []int {
	ind := make([]int, g.Order())
	for _, nbs := range g.LabeledAdjacencyList {
		for _, nb := range nbs {
			ind[nb.To]++
		}
	}
	return ind
}

// IsTree identifies trees in directed graphs.
//
// Return value isTree is true if the subgraph reachable from root is a tree.
// Further, return value allTree is true if the entire graph g is reachable
// from root.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledDirected) IsTree(root NI) (isTree, allTree bool) {
	a := g.LabeledAdjacencyList
	v := bits.New(len(a))
	v.SetAll()
	var df func(NI) bool
	df = func(n NI) bool {
		if v.Bit(int(n)) == 0 {
			return false
		}
		v.SetBit(int(n), 0)
		for _, to := range a[n] {
			if !df(to.To) {
				return false
			}
		}
		return true
	}
	isTree = df(root)
	return isTree, isTree && v.Zero()
}

// StronglyConnectedComponents identifies strongly connected components in
// a directed graph.
//
// The method calls the emit argument for each component identified.  Each
// component is a list of nodes.  The emit function must return true for the
// method to continue identifying components.  If emit returns false, the
// method returns immediately.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// The algorithm here is by David Pearce.  See also alt.SCCPathBased and
// alt.SCCTarjan.
func (g LabeledDirected) StronglyConnectedComponents(emit func([]NI) bool) {
	// See Algorithm 3 PEA FIND SCC2(V,E) in "An Improved Algorithm for
	// Finding the Strongly Connected Components of a Directed Graph"
	// by David J. Pearce.
	a := g.LabeledAdjacencyList
	rindex := make([]int, len(a))
	var S []NI
	index := 1
	c := len(a) - 1
	var visit func(NI) bool
	visit = func(v NI) bool {
		root := true
		rindex[v] = index
		index++
		for _, w := range a[v] {
			if rindex[w.To] == 0 {
				if !visit(w.To) {
					return false
				}
			}
			if rindex[w.To] < rindex[v] {
				rindex[v] = rindex[w.To]
				root = false
			}
		}
		if !root {
			S = append(S, v)
			return true
		}
		var scc []NI
		index--
		for last := len(S) - 1; last >= 0; last-- {
			w := S[last]
			if rindex[v] > rindex[w] {
				break
			}
			S = S[:last]
			rindex[w] = c
			scc = append(scc, w)
			index--
		}
		rindex[v] = c
		c--
		return emit(append(scc, v))
	}
	for v := range a {
		if rindex[v] == 0 && !visit(NI(v)) {
			break
		}
	}
}

// Condensation returns strongly connected components and their
// condensation graph.
//
// Components are ordered in a forward topological ordering.
func (g LabeledDirected) Condensation() (scc [][]NI, cd AdjacencyList) {
	var r [][]NI
	// problems:  1. need to prove that Pearce returns a reverse topological
	// ordering like Tarjan.  2.  Why the reversing?  why not just collect
	// the components and use them as they are?
	g.StronglyConnectedComponents(func(c []NI) bool {
		r = append(r, c)
		return true
	})
	scc = make([][]NI, len(r))
	last := len(r) - 1
	for i, ci := range r {
		scc[last-i] = ci
	}
	cd = make(AdjacencyList, len(scc)) // return value
	cond := make([]NI, g.Order())      // mapping from g node to cd node
	for cn := len(cd) - 1; cn >= 0; cn-- {
		c := scc[cn]
		for _, n := range c {
			cond[n] = NI(cn) // map g node to cd node
		}
		var tos []NI           // list of 'to' nodes
		m := bits.New(len(cd)) // tos map
		m.SetBit(cn, 1)
		for _, n := range c {
			for _, to := range g.LabeledAdjacencyList[n] {
				if ct := cond[to.To]; m.Bit(int(ct)) == 0 {
					m.SetBit(int(ct), 1)
					tos = append(tos, ct)
				}
			}
		}
		cd[cn] = tos
	}
	return
}

// Topological computes a topological ordering of a directed acyclic graph.
//
// For an acyclic graph, return value ordering is a permutation of node numbers
// in topologically sorted order and cycle will be nil.  If the graph is found
// to be cyclic, ordering will be nil and cycle will be the path of a found
// cycle.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledDirected) Topological() (ordering, cycle []NI) {
	i := -1
	return g.dfTopo(func() NI {
		i++
		if i < g.Order() {
			return NI(i)
		}
		return -1
	})
}

func (g LabeledDirected) dfTopo(f func() NI) (ordering, cycle []NI) {
	a := g.LabeledAdjacencyList
	ordering = make([]NI, len(a))
	i := len(ordering)
	temp := bits.New(len(a))
	perm := bits.New(len(a))
	var cycleFound bool
	var cycleStart NI
	var df func(NI)
	df = func(n NI) {
		switch {
		case temp.Bit(int(n)) == 1:
			cycleFound = true
			cycleStart = n
			return
		case perm.Bit(int(n)) == 1:
			return
		}
		temp.SetBit(int(n), 1)
		for _, nb := range a[n] {
			df(nb.To)
			if cycleFound {
				if cycleStart >= 0 {
					// a little hack: orderng won't be needed so repurpose the
					// slice as cycle.  this is read out in reverse order
					// as the recursion unwinds.
					x := len(ordering) - 1 - len(cycle)
					ordering[x] = n
					cycle = ordering[x:]
					if n == cycleStart {
						cycleStart = -1
					}
				}
				return
			}
		}
		temp.SetBit(int(n), 0)
		perm.SetBit(int(n), 1)
		i--
		ordering[i] = n
	}
	for {
		n := f()
		if n < 0 {
			return ordering[i:], nil
		}
		if perm.Bit(int(n)) == 1 {
			continue
		}
		df(n)
		if cycleFound {
			return nil, cycle
		}
	}
}

// TopologicalKahn computes a topological ordering of a directed acyclic graph.
//
// For an acyclic graph, return value ordering is a permutation of node numbers
// in topologically sorted order and cycle will be nil.  If the graph is found
// to be cyclic, ordering will be nil and cycle will be the path of a found
// cycle.
//
// This function is based on the algorithm by Arthur Kahn and requires the
// transpose of g be passed as the argument.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledDirected) TopologicalKahn(tr Directed) (ordering, cycle []NI) {
	// code follows Wikipedia pseudocode.
	var L, S []NI
	// rem for "remaining edges," this function makes a local copy of the
	// in-degrees and consumes that instead of consuming an input.
	rem := make([]int, g.Order())
	for n, fr := range tr.AdjacencyList {
		if len(fr) == 0 {
			// accumulate "set of all nodes with no incoming edges"
			S = append(S, NI(n))
		} else {
			// initialize rem from in-degree
			rem[n] = len(fr)
		}
	}
	for len(S) > 0 {
		last := len(S) - 1 // "remove a node n from S"
		n := S[last]
		S = S[:last]
		L = append(L, n) // "add n to tail of L"
		for _, m := range g.LabeledAdjacencyList[n] {
			// WP pseudo code reads "for each node m..." but it means for each
			// node m *remaining in the graph.*  We consume rem rather than
			// the graph, so "remaining in the graph" for us means rem[m] > 0.
			if rem[m.To] > 0 {
				rem[m.To]--         // "remove edge from the graph"
				if rem[m.To] == 0 { // if "m has no other incoming edges"
					S = append(S, m.To) // "insert m into S"
				}
			}
		}
	}
	// "If graph has edges," for us means a value in rem is > 0.
	for c, in := range rem {
		if in > 0 {
			// recover cyclic nodes
			for _, nb := range g.LabeledAdjacencyList[c] {
				if rem[nb.To] > 0 {
					cycle = append(cycle, NI(c))
					break
				}
			}
		}
	}
	if len(cycle) > 0 {
		return nil, cycle
	}
	return L, nil
}

// TopologicalSubgraph computes a topological ordering of a subgraph of a
// directed acyclic graph.
//
// The subgraph considered is that reachable from the specified node list.
//
// For an acyclic subgraph, return value ordering is a permutation of reachable
// node numbers in topologically sorted order and cycle will be nil.  If the
// subgraph is found to be cyclic, ordering will be nil and cycle will be
// the path of a found cycle.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledDirected) TopologicalSubgraph(nodes []NI) (ordering, cycle []NI) {
	i := -1
	return g.dfTopo(func() NI {
		i++
		if i < len(nodes) {
			return nodes[i]
		}
		return -1
	})
}
