// Copyright (c) 2015 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package sml

import "github.com/eraclitux/stracer"

type Classifier interface {
	// Predict returns a Table
	// which stores predicted labels
	// as single elements rows.
	Predict(Table) (Table, error)
	//Fit()
	//CrossValidation()
}

type ValidationReport struct {
	Recal float64
	Prior float64
}

type kNNClassifier struct {
	trainData Table
	k         int
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
	for i, _ := range samples {
		r := []interface{}{}
		samples[i] = kSample{
			distance: 1,
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
	if d < maxDistance && indexToChange != -1 {
		t[indexToChange].row = row
		t[indexToChange].distance = d
	}
}

// getNearest return the clasified label for
// given slice of k samples.
func (t kSamples) getNearest() string {
	stracer.Traceln("k samples:", t)
	m := make(map[string]int)
	for _, e := range t {
		// get label as last column in row.
		label := e.row[len(e.row)-1].(string)
		if _, ok := m[label]; ok {
			m[label]++
		} else {
			m[label] = 1
		}
	}
	max := 0
	label := ""
	stracer.Traceln("calulated nearest map:", m)
	for k, v := range m {
		// empty string is used to initialize k samples.
		if v > max && k != "" {
			label = k
			max = v
		}
	}
	return label
}

// Predict calculates category for each element in testData.
func (k *kNNClassifier) Predict(testData Table) (Table, error) {
	var prediction memoryTable = make([][]interface{}, testData.Len())
	for j := 0; j < testData.Len(); j++ {
		testRow, err := testData.Row(j)
		if err != nil {
			return nil, err
		}
		samples := newKSamples(k.k)
		for i := 0; i < k.trainData.Len(); i++ {
			trainRow, err := k.trainData.Row(i)
			d, err := distance(testRow, trainRow, nil)
			if err != nil {
				return nil, err
			}
			//stracer.Traceln("trainRow label:", trainRow[len(trainRow)-1], "testRow label:", testRow[len(testRow)-1], "distance:", d)
			samples.checkUpdate(d, trainRow)
		}
		prediction[j] = []interface{}{samples.getNearest()}
	}
	return prediction, nil
}

// NewkNNClassifier returns a new kNN Classifier.
//
// Given m number of training samples and n their number of features,
// current brute force implemetation is at least O(nm²).
func NewkNNClassifier(trainData Table, k int) Classifier {
	return &kNNClassifier{
		trainData: trainData,
		k:         k,
	}
}
