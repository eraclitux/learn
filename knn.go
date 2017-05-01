// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"math"

	kdtree "github.com/hongshibao/go-kdtree"
)

// Classifier models a classification
// problem (binary or multi-labels).
type Classifier interface {
	Predict(Table) (Table, error) // Returns a Table with predicted labels as rows.
}

type kNNBruteForceCls struct {
	trainData Table
	k         int
}

// Predict calculates category for each element in testData.
func (k *kNNBruteForceCls) Predict(testData Table) (Table, error) {
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

type kNNkdTreeCls struct {
	tree *kdtree.KDTree
	k    int
}

// Predict calculates category for each element in testData.
func (k *kNNkdTreeCls) Predict(testData Table) (Table, error) {
	nRows, _ := testData.Caps()
	var prediction MemoryTable = make([][]interface{}, nRows)
	for j := 0; j < nRows; j++ {
		testRow, err := testData.Row(j)
		if err != nil {
			return nil, err
		}
		targetPoint := makeKDTreePoint(testRow)
		neighbours := k.tree.KNN(targetPoint, k.k)
		prediction[j] = []interface{}{nearestLabelInTree(neighbours)}
	}
	return prediction, nil
}

func nearestLabelInTree(neighbours []kdtree.Point) string {
	m := make(map[string]int)
	for _, n := range neighbours {
		label := n.(*kdTreePoint).label
		if _, ok := m[label]; ok {
			m[label]++
		} else {
			m[label] = 1
		}
	}
	max := 0
	label := ""
	for k, v := range m {
		if v > max {
			label = k
			max = v
		}
	}
	return label
}

// NewkNN returns a new kNN Classifier.
// Labels must be stored as last field in Table's rows.
//
// Given m number of training samples and n their number of features,
// if m < 100 brute force is used, otherwise a k-d tree
// is built.
// Brute force implementation is at least O(n*m) but if m is low
// should be a better choice as avoids tree building overhead.
// Search in k-d tree is (n*log(m)) but
// when n > ~20 k-d tree could become O(n*m).
//
// BUG(eraclitux): categorical features are used with brute force
// but not with k-d tree.
func NewkNN(trainData Table, k int) (Classifier, error) {
	// FIXME 100 is arbitrary,
	// algorithm to use
	// should be based on m and n,
	// not just m.
	// Use benchmarks to find (m, n)
	// brute force or k-d tree.
	nRows, _ := trainData.Caps()
	if nRows < 100 {
		return bruteForcekNN(trainData, k)
	}
	return kdTreekNN(trainData, k)
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
		samples[i] = kSample{
			distance: math.MaxFloat64,
			row:      nil,
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
	if d < maxDistance {
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
		tmp := e.row[len(e.row)-1].(*category)
		label := tmp.label
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

type kdTreePoint struct {
	kdtree.Point
	features []float64
	label    string
}

func makeKDTreePoint(row []interface{}) *kdTreePoint {
	// TODO optimization: preallocate this somehow?
	features := []float64{}
	var label string
	for i, e := range row {
		switch v := e.(type) {
		case float64:
			features = append(features, v)
		// BUG categorical features are skipped
		case *category:
			// Last element in row
			// must be the sample's label.
			if i == len(row)-1 {
				label = v.label
			}
		}
	}
	return &kdTreePoint{
		features: features,
		label:    label,
	}
}

func (p *kdTreePoint) Dim() int {
	return len(p.features)
}

func (p *kdTreePoint) GetValue(i int) float64 {
	return p.features[i]
}

func (p *kdTreePoint) Distance(other kdtree.Point) float64 {
	var res float64
	for i := 0; i < p.Dim(); i++ {
		tmp := p.GetValue(i) - other.GetValue(i)
		res += tmp * tmp
	}
	return res
}

func (p *kdTreePoint) PlaneDistance(val float64, i int) float64 {
	tmp := p.GetValue(i) - val
	return tmp * tmp
}

func bruteForcekNN(trainData Table, k int) (*kNNBruteForceCls, error) {
	return &kNNBruteForceCls{
		trainData: trainData,
		k:         k,
	}, nil
}

func kdTreekNN(trainData Table, k int) (*kNNkdTreeCls, error) {
	nRows, _ := trainData.Caps()
	points := make([]kdtree.Point, nRows)
	for i := 0; i < nRows; i++ {
		row, err := trainData.Row(i)
		if err != nil {
			return nil, err
		}
		points[i] = makeKDTreePoint(row)
	}
	return &kNNkdTreeCls{
		tree: kdtree.NewKDTree(points),
		k:    k,
	}, nil
}
