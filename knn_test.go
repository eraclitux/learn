// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

// Run benchmarks:
// go test -run NONE -bench . -benchmem

import (
	"reflect"
	"testing"
)

func TestKSamples_CheckUpdate(t *testing.T) {
	samples := newKSamples(3)
	expectedSamples := newKSamples(3)
	row := []interface{}{0.2, 0.4, 0.1, "one"}
	samples.checkUpdate(0.9, row)
	row = []interface{}{0.2, 0.4, 0.1, "one"}
	samples.checkUpdate(0.8, row)
	row = []interface{}{0.1, 0.3, 0.0, "two"}
	samples.checkUpdate(0.4, row)
	row = []interface{}{0.1, 0.3, 0.0, "two"}
	samples.checkUpdate(0.3, row)
	// checkUpdate does not enforce any order,
	// we specify a known order for expectedSamples
	row = []interface{}{0.1, 0.3, 0.0, "two"}
	expectedSamples[0] = kSample{
		distance: 0.3,
		row:      row,
	}
	row = []interface{}{0.2, 0.4, 0.1, "one"}
	expectedSamples[1] = kSample{
		distance: 0.8,
		row:      row,
	}
	row = []interface{}{0.1, 0.3, 0.0, "two"}
	expectedSamples[2] = kSample{
		distance: 0.4,
		row:      row,
	}

	if !reflect.DeepEqual(samples, expectedSamples) {
		t.Fatalf("expected: %v got: %v", expectedSamples, samples)
	}
}

func TestKSamples_GetNearest(t *testing.T) {
	samples := newKSamples(5)
	row := []interface{}{0.2, 0.4, 0.1, "one"}
	samples.checkUpdate(0.5, row)
	row = []interface{}{0.1, 0.3, 0.0, "one"}
	samples.checkUpdate(0.5, row)
	row = []interface{}{0.1, 0.3, 0.0, "two"}
	samples.checkUpdate(0.4, row)
	row = []interface{}{0.14, 0.33, 0.23, "two"}
	samples.checkUpdate(0.6, row)
	row = []interface{}{0.13, 0.33, 0.23, "two"}
	samples.checkUpdate(0.6, row)

	nearest := samples.getNearest()

	if nearest != "two" {
		t.Fatal("nearest:", nearest)
	}
}

func loadTestSet(t *testing.T) (Table, []float64, []float64) {
	trainSet, err := ReadAllCSV("datasets/iris.csv")
	if err != nil {
		t.Fatal(err)
	}
	mu, sigma, err := Normalize(trainSet, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	return trainSet, mu, sigma
}

func TestBruteForcekNN(t *testing.T) {
	trainSet, mu, sigma := loadTestSet(t)
	clf, err := bruteForcekNN(trainSet, 3)
	if err != nil {
		t.Fatal(err)
	}
	// Categorize single sample.
	var testSet MemoryTable = make([][]interface{}, 1)
	testSet[0] = []interface{}{5.2, 3.4, 1.3, 0.1}
	_, _, err = Normalize(testSet, mu, sigma)
	if err != nil {
		t.Fatal(err)
	}
	prediction, err := clf.Predict(testSet)
	if err != nil {
		t.Fatal(err)
	}
	r, err := prediction.Row(0)
	if err != nil {
		t.Fatal(err)
	}
	expected := "setosa"
	if r[0] != expected {
		t.Errorf("want: %s, got: %s", expected, r[0])
	}
}

func TestKdTreekNN(t *testing.T) {
	trainSet, mu, sigma := loadTestSet(t)
	clf, err := kdTreekNN(trainSet, 3)
	if err != nil {
		t.Fatal(err)
	}
	// Categorize single sample.
	var testSet MemoryTable = make([][]interface{}, 1)
	testSet[0] = []interface{}{5.2, 3.4, 1.3, 0.1}
	_, _, err = Normalize(testSet, mu, sigma)
	if err != nil {
		t.Fatal(err)
	}
	prediction, err := clf.Predict(testSet)
	if err != nil {
		t.Fatal(err)
	}
	r, err := prediction.Row(0)
	if err != nil {
		t.Fatal(err)
	}
	expected := "setosa"
	if r[0] != expected {
		t.Errorf("want: %s, got: %s", expected, r[0])
	}
}

//
// Benchmarks
//
func BenchmarkNN(b *testing.B) {
	trainSet, err := ReadAllCSV("datasets/iris.csv")
	if err != nil {
		b.Fatal(err)
	}
	mu, sigma, err := Normalize(trainSet, nil, nil)
	if err != nil {
		b.Fatal(err)
	}
	// Categorize single sample.
	var testSet MemoryTable = make([][]interface{}, 1)
	testSet[0] = []interface{}{5.2, 3.4, 1.3, 0.1}
	_, _, err = Normalize(testSet, mu, sigma)
	// Benchmark classifier creation
	// and prediction.
	b.Run("kdTree", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			clf, _ := kdTreekNN(trainSet, 3)
			_, err := clf.Predict(testSet)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	clfTree, _ := kdTreekNN(trainSet, 3)
	// Benchmark only the prediction.
	b.Run("kdTree-pdct", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := clfTree.Predict(testSet)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	// Benchmark classifier creation
	// and prediction.
	b.Run("bruteF", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			clf, _ := bruteForcekNN(trainSet, 3)
			_, err := clf.Predict(testSet)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	clfBrute, _ := bruteForcekNN(trainSet, 3)
	// Benchmark only the prediction.
	b.Run("bruteF-pdct", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := clfBrute.Predict(testSet)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
