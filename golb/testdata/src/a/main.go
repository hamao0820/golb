package main

import (
	"fmt"

	"github.com/hamao0820/ac-library-go/lib/sample"
	"github.com/hamao0820/ac-library-go/lib/sample/nested"
	segmenttree "github.com/hamao0820/ac-library-go/lib/segment_tree"
	vc "github.com/hamao0820/ac-library-go/lib/vector"
)

func main() {
	v1 := vc.NewVector(1, 2)
	v2 := vc.NewVector(3, 4)
	v3 := v1.Add(v2) // v3 = v1 + v2
	// v3.Scale(2)

	st := segmenttree.NewSegmentTree(10, 0, func(a, b int) int { return a + b })
	st.Update(0, 1)

	fmt.Println(v3.X, v3.Y)
	fmt.Println(sample.Sample(1, 2))
	fmt.Println(nested.Nested())
	hello()
}

func hello() {
	fmt.Println("hello")
}

func unused() {
	fmt.Println("unused")
}
