package main

import (
	"bytes"
	"testing"
)

func TestWriteMetrics(t *testing.T) {
	suppliers := []Supplier{
		Supplier{0, "Test Restaurant", 1, 0.0},
	}
	var b bytes.Buffer

	writeMetrics(suppliers, &b)

	expected := "reachedOrderValue{supplier=\"Test Restaurant\",office=\"1\"} 0.00\n"
	if b.String() != expected {
		t.Errorf("Expected %s to equal %s", b.String(), expected)
	}
}
