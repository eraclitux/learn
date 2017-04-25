// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import "math"

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
// (quantitative, nominal, cardinal, binary).
func elementsDistance(a1, a2 interface{}) (d float64, er error) {
	// TODO add Geo type of lat/long with distance (http://www.movable-type.co.uk/scripts/latlong.html)
	switch a1.(type) {
	case float64:
		return manhattan(a1.(float64), a2.(float64)), nil
		//return euclidean(a1.(float64), a2.(float64)), nil
	case *category:
		return a1.(*category).distance(a2.(*category)), nil
	default:
		return -1, unknownType(a1)
	}
}

// distance calculates distance between
// two different rows using average of single elements
// distance to account heterogeneous data.
func distance(s, v []interface{}, weights []float64) (float64, error) {
	// FIXME check that ∈ of weights are <=1
	var total float64
	// Some feature are ignored (es string)
	// we cannot use len(s) to calculate the average.
	var numFeatures float64
	for i, e := range s {
		// FIXME refactor with type switch
		// or move this into elementsDistance
		// Ignore string features.
		if _, ok := e.(string); ok {
			continue
		}
		t, err := elementsDistance(e, v[i])
		if err != nil {
			return -1, err
		}
		if weights == nil {
			total += t
		} else {
			total += t * weights[i]
		}
		numFeatures++
	}
	return total / numFeatures, nil
}
