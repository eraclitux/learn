package nml

import (
	"math"
	"testing"
)

func TestComputeAccuracy(t *testing.T) {
	var cm = ConfMatrix{mm: make(map[string]map[string]int)}
	cm.mm["true"] = make(map[string]int)
	cm.mm["false"] = make(map[string]int)
	cm.mm["true"]["true"] = 5
	cm.mm["true"]["false"] = 0
	cm.mm["false"]["true"] = 0
	cm.mm["false"]["false"] = 10
	oA := computeAccuracy(cm)
	if oA != 1 {
		t.Fatal("expected an accuracy of 1, got:", oA)
	}
	cm.mm["true"]["true"] = 0
	cm.mm["true"]["false"] = 1
	cm.mm["false"]["true"] = 8
	cm.mm["false"]["false"] = 0

	oA = computeAccuracy(cm)
	if oA != 0 {
		t.Fatal("expected an accuracy of 0, got:", oA)
	}

}

func TestComputePrecision(t *testing.T) {
	m := make(map[string]int)
	m["one"] = 10
	m["two"] = 0
	m["three"] = 0
	p := computePrecision("one", m)
	if p != 1 {
		t.Fatal("expected a precision of 1, got:", p)
	}
	m["one"] = 0
	m["two"] = 10
	m["three"] = 1110
	p = computePrecision("one", m)
	if p != 0 {
		t.Fatal("expected a precision of 0, got:", p)
	}

	m["one"] = 0
	m["two"] = 0
	m["three"] = 0
	p = computePrecision("one", m)
	if !math.IsNaN(p) {
		t.Fatal("expected a precision of NaN, got:", p)
	}
}

func TestComputeRecall(t *testing.T) {
	var cm = ConfMatrix{mm: make(map[string]map[string]int)}
	cm.mm["one"] = map[string]int{"one": 5, "two": 0, "three": 1}
	cm.mm["two"] = map[string]int{"one": 0, "two": 6, "three": 1}
	cm.mm["three"] = map[string]int{"one": 0, "two": 1, "three": 1}
	r := computeRecall("one", cm)
	if r != 1 {
		t.Fatal("expected a recall of 1, got:", r)
	}
	cm.mm["one"] = map[string]int{"one": 0, "two": 0, "three": 1}
	cm.mm["two"] = map[string]int{"one": 1, "two": 6, "three": 1}
	cm.mm["three"] = map[string]int{"one": 0, "two": 1, "three": 1}
	r = computeRecall("one", cm)
	if r != 0 {
		t.Fatal("expected a recall of 0, got:", r)
	}
}
