// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"errors"
	"math"
)

func hammingD(a, b uint) uint {
	var dist uint
	val := a ^ b
	for val != 0 {
		dist++
		val &= val - 1
	}
	return dist
}

func manhattan(a, b float64) float64 {
	d := a - b
	if d < 0 {
		d *= -1
	}
	return d
}

// TODO benchmark vs manhattan
func euclidean(a, b float64) float64 {
	return math.Pow(a-b, 2)
}

// elementsDistance returns distance for two elements of same type
// (quantitative or categorical).
func elementsDistance(a1, a2 interface{}) (d float64, er error) {
	switch v1 := a1.(type) {
	case float64:
		v2, ok := a2.(float64)
		if !ok {
			return -1, typeMismatchErr(a1, a2)
		}
		return manhattan(v1, v2), nil
		//return euclidean(v1, v2), nil
	case *category:
		v2, ok := a2.(*category)
		if !ok {
			return -1, typeMismatchErr(a1, a2)
		}
		return v1.distance(v2), nil
	default:
		return -1, unknownTypeErr(a1)
	}
}

// distance calculates distance between
// two different rows using average of single elements
// distance to account heterogeneous
// features (numerical & categorical).
// Last element in trainRow
// is considered label if its type is category.
func distance(testRow, trainRow []interface{}, weights []float64) (float64, error) {
	// FIXME check that ∈ of weights are <=1
	if len(trainRow) <= len(testRow) {
		return math.NaN(), errors.New("learn: insufficient number of features in train sample")
	}
	var total float64
	for i, e := range testRow {
		t, err := elementsDistance(e, trainRow[i])
		if err != nil {
			return -1, err
		}
		if weights == nil {
			total += t
		} else {
			total += t * weights[i]
		}
	}
	numFeatures := float64(len(testRow))
	return total / numFeatures, nil
}
