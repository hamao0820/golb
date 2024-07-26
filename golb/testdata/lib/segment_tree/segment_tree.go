package segmenttree

type SegmentTree[T any] struct {
	data []T
	n    int
	e    T
	op   func(T, T) T
}

func NewSegmentTree[T any](n int, e T, op func(T, T) T) *SegmentTree[T] {
	segtree := new(SegmentTree[T])
	segtree.e = e
	segtree.op = op
	segtree.n = 1
	for segtree.n < n {
		segtree.n *= 2
	}
	segtree.data = make([]T, segtree.n*2-1)
	for i := 0; i < segtree.n*2-1; i++ {
		segtree.data[i] = segtree.e
	}
	return segtree
}

func (segtree *SegmentTree[T]) Update(idx int, x T) {
	idx += segtree.n - 1
	segtree.data[idx] = x
	for 0 < idx {
		idx = (idx - 1) / 2
		segtree.data[idx] = segtree.op(segtree.data[idx*2+1], segtree.data[idx*2+2])
	}
}

func (segtree *SegmentTree[T]) query(begin, end, idx, a, b int) T {
	if b <= begin || end <= a {
		return segtree.e
	}
	if begin <= a && b <= end {
		return segtree.data[idx]
	}
	v1 := segtree.query(begin, end, idx*2+1, a, (a+b)/2)
	v2 := segtree.query(begin, end, idx*2+2, (a+b)/2, b)
	return segtree.op(v1, v2)
}

func (segtree *SegmentTree[T]) Query(begin, end int) T {
	return segtree.query(begin, end, 0, 0, segtree.n)
}
