package set

type Set[T comparable] struct {
	values map[T]struct{}
}

func New[T comparable]() *Set[T] {
	return &Set[T]{
		values: make(map[T]struct{}),
	}
}

func (s *Set[T]) Add(value T) {
	s.values[value] = struct{}{}
}

func (s *Set[T]) AddArray(values []T) {
	for i := range values {
		s.values[values[i]] = struct{}{}
	}
}

func (s *Set[T]) Remove(value T) {
	delete(s.values, value)
}

func (s *Set[T]) ToArray() []T {
	var array []T
	for value := range s.values {
		array = append(array, value)
	}
	return array
}
