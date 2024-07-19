package lib

type Vector struct {
	X, Y int
}

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
