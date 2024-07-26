package main

import (
	"fmt"
	"math"
)

func main() {
	v1 := NewVector(1, 2)
	v2 := NewVector(3, 4)
	v3 := v1.Add(v2)
	fmt.Println(v3.X, v3.Y)
	fmt.Println(Sample(1, 2))
	fmt.Println(Nested())
	hello()
}
func hello() {
	fmt.Println("hello")
}

//--------------------------------------------------以下は生成コード--------------------------------------------------

// util/util.go
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// vector/vector.go
type Vector struct{ X, Y int }

func NewVector(x, y int) Vector {
	return Vector{x, y}
}
func (v Vector) Add(v2 Vector) Vector {
	return Vector{v.X + v2.X, v.Y + v2.Y}
}
func (v *Vector) Scale(s int) {
	v.X *= s
	v.Y *= s
}

// nested/nested.go
func Nested() string {
	return "{{nested}}"
}

// sample/sample.go
func Sample(x, y int) int {
	return Min(x, y)
}

// util/consts.go
const (
	MOD998    = 998244353
	MOD107    = 1000000007
	MAX       = math.MaxInt64
	MIN       = math.MinInt64
	Yes       = "Yes"
	No        = "No"
	YES       = "YES"
	NO        = "NO"
	Takahashi = "Takahashi"
	Aoki      = "Aoki"
)
