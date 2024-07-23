package main

import (
	"fmt"
	"golb/golb/testdata/lib/util"
	"golb/golb/testdata/lib/vector"
)

func main() {
	v1 := vector.NewVector(1, 2)
	v2 := vector.NewVector(3, 4)
	v3 := v1.Add(v2)
	v3.Scale(2)

	fmt.Println(v3.X, v3.Y)
	fmt.Println(util.Min(1, 2))
	hello()
}

func hello() {
	fmt.Println("hello")
}
