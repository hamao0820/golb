package main

import (
	"fmt"

	"github.com/hamao0820/ac-library-go/lib/sample"
	"github.com/hamao0820/ac-library-go/lib/util"
	vc "github.com/hamao0820/ac-library-go/lib/vector"
)

func main() {
	v1 := vc.NewVector(1, 2)
	v2 := vc.NewVector(3, 4)
	v3 := v1.Add(v2)
	v3.Scale(2)

	fmt.Println(v3.X, v3.Y)
	fmt.Println(util.Min(1, 2))
	fmt.Println(sample.Sample(1, 2))
	hello()
}

func hello() {
	fmt.Println("hello")
}
