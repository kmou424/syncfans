package genericmap

type Pair[L, R any] struct {
	Left  L
	Right R
}

func NewPair[L, R any](left L, right R) Pair[L, R] {
	return Pair[L, R]{left, right}
}

func (p Pair[L, R]) Swap() Pair[R, L] {
	return Pair[R, L]{p.Right, p.Left}
}
