package test

import "testing"

func TestDemo(t *testing.T) {
	around(func(c client) {
		t.Log("fuck")
	})
}
