package main

import "testing"

func TestStack(t *testing.T) {
	s := new_bytestack()
	s.in('a')
	s.in('b')

	if s.out() != 'b' {
		t.Errorf("out() failed")
	}

	if s.peek() != 'a' {
		t.Errorf("peek() failed")
	}

	if s.out() != 'a' {
		t.Errorf("out() failed")
	}

	if !s.is_empty() {
		t.Errorf("is_empty() failed")
	}
}
