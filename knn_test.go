// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

// Run benchmarks:
// go test -run NONE -bench . -benchmem

import (
	"fmt"
	"reflect"
	"testing"
)

func TestKSamples_CheckUpdate(t *testing.T) {
	samples := newKSamples(3)
	expectedSamples := newKSamples(3)
	// FIXME converte string category
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
	row := []interface{}{0.2, 0.4, 0.1, newCategory("one", nil)}
	samples.checkUpdate(0.5, row)
	row = []interface{}{0.1, 0.3, 0.0, newCategory("one", nil)}
	samples.checkUpdate(0.5, row)
	row = []interface{}{0.1, 0.3, 0.0, newCategory("two", nil)}
	samples.checkUpdate(0.4, row)
	row = []interface{}{0.14, 0.33, 0.23, newCategory("two", nil)}
	samples.checkUpdate(0.6, row)
	row = []interface{}{0.13, 0.33, 0.23, newCategory("two", nil)}
	samples.checkUpdate(0.6, row)

	nearest := samples.getNearest()
	if nearest != "two" {
		t.Fatal("nearest:", nearest)
	}
}

func loadTrainSet(t *testing.T, set string) (Table, []float64, []float64, []string) {
	path := fmt.Sprintf("datasets/%s.csv", set)
	trainSet, err := ReadAllCSV(path)
	if err != nil {
		t.Fatal(err)
	}
	mu, sigma, catSet, err := Normalize(trainSet, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	return trainSet, mu, sigma, catSet
}

// Test bruteForcekNN using dataset with
// a numerical only features.
func TestBruteForcekNN_numerical(t *testing.T) {
	trainSet, mu, sigma, catSet := loadTrainSet(t, "iris")
	clf, err := bruteForcekNN(trainSet, 3)
	if err != nil {
		t.Fatal(err)
	}
	// Categorize single sample.
	var testSet MemoryTable = make([][]interface{}, 1)
	testSet[0] = []interface{}{5.2, 3.4, 1.3, 0.1}
	_, _, _, err = Normalize(testSet, mu, sigma, catSet)
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

// Test bruteForcekNN using dataset with
// numerical and categorical features.
func TestBruteForcekNN_mixed(t *testing.T) {
	trainSet, mu, sigma, catSet := loadTrainSet(t, "adult_train")
	clf, err := bruteForcekNN(trainSet, 3)
	if err != nil {
		t.Fatal(err)
	}
	// Categorize single sample.
	var testSet MemoryTable = make([][]interface{}, 1)
	testSet[0] = []interface{}{
		25.0, "Private", 226802.0, "11th", 7.0, "Never-married", "Machine-op-inspct", "Own-child", "Black", "Male", 0.0, 0.0, 40.0, "United-States",
	}
	_, _, _, err = Normalize(testSet, mu, sigma, catSet)
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
	expected := "<=50K"
	if r[0] != expected {
		t.Errorf("want: %s, got: %s", expected, r[0])
	}
}

func TestKdTreekNN(t *testing.T) {
	trainSet, mu, sigma, catSet := loadTrainSet(t, "iris")
	clf, err := kdTreekNN(trainSet, 3)
	if err != nil {
		t.Fatal(err)
	}
	// Categorize single sample.
	var testSet MemoryTable = make([][]interface{}, 1)
	testSet[0] = []interface{}{5.2, 3.4, 1.3, 0.1}
	_, _, _, err = Normalize(testSet, mu, sigma, catSet)
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
	mu, sigma, catSet, err := Normalize(trainSet, nil, nil, nil)
	if err != nil {
		b.Fatal(err)
	}
	// Categorize single sample.
	var testSet MemoryTable = make([][]interface{}, 1)
	testSet[0] = []interface{}{5.2, 3.4, 1.3, 0.1}
	_, _, _, err = Normalize(testSet, mu, sigma, catSet)
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
