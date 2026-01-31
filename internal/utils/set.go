package utils

type set[T comparable] struct {
	m map[T]bool
}

func NewSet[T comparable]() *set[T] {
	return &set[T]{m: map[T]bool{}}
}
func (s *set[T]) Contains(x T) bool {
	_, present := s.m[x]
	return present
}

func (s *set[T]) Add(x T) bool {
	if s.Contains(x) {
		return false
	}
	s.m[x] = true
	return true
}

func (s *set[T]) Len() int {
	return len(s.m)
}

func (s *set[T]) Erase(x T) bool {
	if !s.Contains(x) {
		return false
	}
	delete(s.m, x)
	return true
}

func (s *set[T]) ToList() []T {
	result := make([]T, 0, len(s.m))
	for key := range s.m {
		result = append(result, key)
	}
	return result
}
