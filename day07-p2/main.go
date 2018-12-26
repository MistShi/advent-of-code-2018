package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {
	numWorkers := flag.Int("w", 5, "number of workers")
	baseDuration := flag.Int("d", 60, "base duration")
	flag.Parse()

	g, err := readGraph("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	h := makeHeap(g)

	workers := make([]int, *numWorkers)
	for len(h) > 0 {
		fmt.Println("-----------")
		fmt.Printf("heap: %v\n", h)
		fmt.Printf("workers: %v\n", workers)

		n := h.pop()
		follows := make([]string, 0, len(n.follows))
		for _, f := range n.follows {
			follows = append(follows, f.name)
		}
		fmt.Printf("%s %v\n", n.name, follows)
		w, t := minArg(workers)
		if t < n.availableAt {
			t = n.availableAt
		}
		workers[w] = n.startAt(t, *baseDuration)
		fmt.Printf("worker %d starting task %s at %d, will be done at %d\n", t, n.name, w, workers[w])
	}

	fmt.Println(max(workers))
}

func minArg(vs []int) (i, v int) {
	minI, minV := 0, vs[0]
	for i, v := range vs {
		if minV > v {
			minV = v
			minI = i
		}
	}
	return minI, minV
}

func max(vs []int) int {
	m := 0
	for _, v := range vs {
		if m < v {
			m = v
		}
	}
	return m
}

type node struct {
	name        string
	deps        []*node
	waitingFor  int
	follows     []*node
	doneAt      int
	availableAt int
}

func (n *node) duration() int { return int(n.name[0]) - 'A' + 1 }

func (n *node) startAt(t, baseDuration int) int {
	n.doneAt = t + n.duration() + baseDuration

	for _, f := range n.follows {
		f.oneDepDone()
	}
	return n.doneAt
}

func (n *node) oneDepDone() {
	n.waitingFor--
	if n.waitingFor > 0 {
		return
	}

	for _, d := range n.deps {
		if n.availableAt < d.doneAt {
			n.availableAt = d.doneAt
		}
	}
}

type graph map[string]*node

func readGraph(path string) (graph, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	g := make(graph)

	s := bufio.NewScanner(f)
	for s.Scan() {
		var a, b string
		fmt.Sscanf(s.Text(), "Step %s must be finished before step %s can begin.", &a, &b)
		g.addDep(a, b)
	}

	return g, s.Err()
}

func (g graph) addDep(a, b string) {
	na := g.node(a)
	nb := g.node(b)
	na.follows = append(na.follows, nb)
	nb.deps = append(nb.deps, na)
	nb.availableAt = -1
	nb.waitingFor++
}

func (g graph) node(name string) *node {
	if _, ok := g[name]; !ok {
		g[name] = &node{name: name}
	}
	return g[name]
}

type heap []*node

func makeHeap(g graph) heap {
	h := make(heap, 0, len(g))
	for _, n := range g {
		h = append(h, n)
	}
	return h
}

func (h heap) String() string {
	h.sort()
	names := make([]string, 0, len(h))
	for _, n := range h {
		if n.waitingFor > 0 {
			continue
		}
		names = append(names, fmt.Sprintf("%s [%d @%d]", n.name, n.waitingFor, n.availableAt))
	}
	return strings.Join(names, " | ")
}

func (h *heap) pop() *node {
	if len(*h) == 0 {
		return nil
	}

	h.sort()
	r := (*h)[0]
	*h = (*h)[1:]
	return r
}

func (h heap) sort() {
	sort.Slice(h, func(i, j int) bool {
		di, dj := h[i].waitingFor, h[j].waitingFor
		if di < dj {
			return true
		} else if dj < di {
			return false
		}
		ai, aj := h[i].availableAt, h[j].availableAt
		if ai >= 0 && ai < aj {
			return true
		} else if aj >= 0 && aj < ai {
			return false
		}
		return h[i].name < h[j].name
	})
}