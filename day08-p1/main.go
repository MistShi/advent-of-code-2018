package main

import (
	"fmt"
	"log"

	"github.com/campoy/advent-of-code-2018/aoc"
)

func main() {
	nums, err := aoc.ReadIntSeqFromFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	n := readNode(nums)
	fmt.Println(n.sumMetadata())
}

type node struct {
	children []node
	metadata []int
}

func readNode(nums *aoc.IntSeq) node {
	n := node{
		children: make([]node, nums.Next()),
		metadata: make([]int, nums.Next()),
	}
	for i := range n.children {
		n.children[i] = readNode(nums)
	}
	for i := range n.metadata {
		n.metadata[i] = nums.Next()
	}
	return n
}

func (n node) sumMetadata() int {
	sum := aoc.Sum(n.metadata)
	for _, c := range n.children {
		sum += c.sumMetadata()
	}
	return sum
}
