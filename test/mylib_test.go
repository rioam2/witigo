package test

import (
	"testing"

	witigo "github.com/rioam2/witigo/pkg"
)

func TestExampleFunction(t *testing.T) {
	result := witigo.MyFunction()        // Replace with actual function call
	expected := "Hello from MyFunction!" // Replace with expected result

	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
