package main

import (
	"fmt"
	"testing"
)

func TestRe2Postfix(t *testing.T) {
	re := "a(gb|c)*g"
	postfix := re2postfix(re)
	fmt.Println(postfix)
}
