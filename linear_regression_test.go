// Copyright (c) 2015 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import "testing"

func loadDataset(t *testing.T) Table {
	// Load all data in memory.
	trainData, err := ReadAllCSV("datasets/anscombe.csv")
	if err != nil {
		t.Fatal(err)
	}
	return trainData

}

func TestNewLinearRegression(t *testing.T) {
	trainData := loadDataset(t)
	var tab memoryTable = make([][]interface{}, 1)
	tab[0] = []interface{}{5.0, 5.0, 5.0, 8.0, 5.68, 4.74, 5.73, 6.89}
	lr, err := NewLinearRegression(trainData)
	if err != nil {
		t.Fatal(err)
	}
	y, err := lr.Predict(tab)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(y)
	t.Fatal("test me")
}

func TestLinearRegression_Fit(t *testing.T) {
	t.Fatal("testing me")
}
func TestLinearRegression_Predict(t *testing.T) {
	var tab memoryTable = [][]interface{}{
		[]interface{}{"not", "allowed"},
	}
	t.Log(tab)
	t.Fail()
}
