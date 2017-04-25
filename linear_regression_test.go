// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"math"
	"testing"
)

func loadDataset(t *testing.T) Table {
	// Load all data in memory.
	trainData, err := ReadAllCSV("datasets/linear_test.csv")
	if err != nil {
		t.Fatal(err)
	}
	return trainData
}

const floatTolerance = 1e-7

func floatsAreEqual(a, b float64) bool {
	equality := false
	if math.IsNaN(a) && math.IsNaN(b) {
		equality = true
	}
	if e := math.Abs(a - b); e < floatTolerance {
		equality = true
	}
	return equality
}

func TestLinearRegression_Predict(t *testing.T) {
	trainData := loadDataset(t)
	var tab MemoryTable = make([][]interface{}, 1)
	// Hand crafted case.
	tab[0] = []interface{}{1650.0, 3.0}
	lr, err := NewLinearRegression(trainData)
	if err != nil {
		t.Fatal(err)
	}
	y, err := lr.Predict(tab)
	if err != nil {
		t.Fatal(err)
	}
	wanted := 293081.464335
	if !floatsAreEqual(y[0], wanted) {
		t.Fatalf("wrong prediction, want: %f got: %f\n", y[0], wanted)
	}
}
