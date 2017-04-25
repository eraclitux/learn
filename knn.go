// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"math"
)

// Classifier models a classification
// problem (binary or multi-labels).
type Classifier interface {
	// Predict returns a Table
	// which stores predicted labels
	// as single fields rows.
	Predict(Table) (Table, error)
}

type kNNClassifier struct {
	trainData Table
	k         int
}

// Predict calculates category for each element in testData.
func (k *kNNClassifier) Predict(testData Table) (Table, error) {
	nRows, _ := testData.Caps()
	var prediction MemoryTable = make([][]interface{}, nRows)
	for j := 0; j < nRows; j++ {
		testRow, err := testData.Row(j)
		if err != nil {
			return nil, err
		}
		samples := newKSamples(k.k)
		trainDataRows, _ := k.trainData.Caps()
		for i := 0; i < trainDataRows; i++ {
			trainRow, err := k.trainData.Row(i)
			d, err := distance(testRow, trainRow, nil)
			if err != nil {
				return nil, err
			}
			samples.checkUpdate(d, trainRow)
		}
		prediction[j] = []interface{}{samples.getNearest()}
	}
	return prediction, nil
}

type kSample struct {
	row      []interface{}
	distance float64
}

type kSamples []kSample

// newKSamples initialize kSamples
// with maximum distance.
func newKSamples(n int) kSamples {
	var samples kSamples = make([]kSample, n)
	for i := range samples {
		// FIXME if nil OK avoid this allocation
		r := []interface{}{}
		samples[i] = kSample{
			distance: math.MaxFloat64,
			row:      r,
		}
	}
	return samples
}

// checkUpdate checks if row is nearer that the others stored,
// updating samples in case.
func (t kSamples) checkUpdate(d float64, row []interface{}) {
	indexToChange := -1
	var maxDistance float64
	for i, e := range t {
		if e.distance > maxDistance {
			maxDistance = e.distance
			indexToChange = i
		}
	}
	// FIXME is -1 check really needed?
	if d < maxDistance && indexToChange != -1 {
		t[indexToChange].row = row
		t[indexToChange].distance = d
	}
}

// getNearest returns the classified label for
// given slice of k samples.
func (t kSamples) getNearest() string {
	m := make(map[string]int)
	for _, e := range t {
		// get label as last column in row.
		// FIXME check this assertion
		label := e.row[len(e.row)-1].(string)
		if _, ok := m[label]; ok {
			m[label]++
		} else {
			m[label] = 1
		}
	}
	max := 0
	label := ""
	for k, v := range m {
		// empty string is used to initialize k samples.
		if v > max && k != "" {
			label = k
			max = v
		}
	}
	return label
}

// NewkNN returns a new kNN Classifier.
// Categories must be stored as last field in Table's rows.
//
// Given m number of training samples and n their number of features,
// current brute force implemetation is at least O(nm²).
func NewkNN(trainData Table, k int) Classifier {
	return &kNNClassifier{
		trainData: trainData,
		k:         k,
	}
}
