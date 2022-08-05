package main

type bytestack struct {
	vals []byte
}

func new_bytestack() *bytestack {
	return &bytestack{
		vals: make([]byte, 0),
	}
}

func (s *bytestack) in(val byte) {
	s.vals = append(s.vals, val)
}

func (s *bytestack) out() byte {
	val := s.vals[len(s.vals)-1]
	s.vals = s.vals[:len(s.vals)-1]
	return val
}

func (s *bytestack) peek() byte {
	return s.vals[len(s.vals)-1]
}

func (s *bytestack) is_empty() bool {
	return len(s.vals) == 0
}
