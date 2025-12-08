package pool

type Reseter interface {
	Reset()
}

type Pool[T Reseter] struct {
	Store []T
}

func NewPool[T Reseter]() *Pool[T] {
	return &Pool[T]{
		Store: make([]T, 0, 10),
	}
}

func (p *Pool[T]) Get() T {
	if len(p.Store) > 0 {
		tmpP := p.Store[len(p.Store)-1]
		p.Store = p.Store[:len(p.Store)-1]
		return tmpP
	}
	var empty T
	return empty
}

func (p *Pool[T]) Put(item T) {
	p.Store = append(p.Store, item)
}
