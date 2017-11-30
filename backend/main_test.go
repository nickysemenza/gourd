package main

import "testing"

func TestTruth(t *testing.T) {
	if false {
		t.Error("oh no")
	}
}
