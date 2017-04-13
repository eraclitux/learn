// Copyright (c) 2015 Andrea Masi. All rights reserved.
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

const tolerance = 1e-7

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
	if e := math.Abs(y[0] - 293081.464335); e > tolerance {
		t.Fatal("prediction over tolerance:", e)
	}
}
