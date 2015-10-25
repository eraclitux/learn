// Copyright (c) 2015 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

// Package sml exposes some machine learning alghoritms.
//
// Unsupervised learning:
//
//	- k means clustering
//
// K means clustering
//
// Example of data
//
//	Hours	Choices		Stars	Price
//	12,	"A,C",		5,	15.10
//	1,	"D"		1,	1
//
// Categorical features are supported. They must be translated to an array of 0 and 1:
//
//	Hours	Choices		Stars	Price
//	12,	"[1,0,1,0]",	5,	15.10
//	1,	"[0,0,0,1]"	1,	1
package sml

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/eraclitux/stracer"
)

func hummingD(a, b uint) uint {
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

// elementsDistance returns distance for two elelemnt of same type
// (quantitative, nominal, cardinal, binary).
// Returning value is ∈ [0,1].
func elementsDistance(a1, a2 interface{}) (d float64, er error) {
	switch a1.(type) {
	case float64:
		return manhattan(a1.(float64), a2.(float64)), nil
	case *Category:
		return a1.(*Category).distance(a2.(*Category)), nil
	default:
		return -1, UnknownType
	}
}

// distance calculates distance between
// two different rows using average of single elements
// distance to account heterogeneous data.
// FIXME check that ∈ of weights are <=1
func distance(s, v []interface{}, weights []float64) (float64, error) {
	var total float64
	for i, e := range s {
		t, err := elementsDistance(e, v[i])
		if err != nil {
			return -1, err
		}
		if weights == nil {
			total += t
		} else {
			total += t * weights[i]
		}
	}
	return total / float64(len(s)), nil
}

func createRandCategory(l uint) *Category {
	sS := []string{}
	for i := 0; i < int(l); i++ {
		sN := strconv.Itoa(rand.Intn(2))
		sS = append(sS, sN)
	}
	return newCategory(strings.Join(sS, ","))
}
func createRandomCentroids(k int, s []interface{}) ([][]interface{}, error) {
	l := len(s)
	r := make([][]interface{}, 0, k)
	for i := 0; i < k; i++ {
		c := make([]interface{}, l)
		for i, e := range s {
			switch e.(type) {
			case float64:
				c[i] = rand.Float64()
			case *Category:
				c[i] = createRandCategory(e.(*Category).choicesN)
			default:
				return nil, UnknownType
			}
		}
		r = append(r, c)
	}
	return r, nil
}

func zeroCentroid(c []interface{}) {
	for i, e := range c {
		switch e.(type) {
		case float64:
			c[i] = float64(0)
		case *Category:
			c[i].(*Category).zero()
		default:
			panic("unknown type zeroing centroid")
		}
	}
}

// incrementCentroid adds quantities to centorid elements
// to calculate the mean after.
func incrementCentroid(c []interface{}, d []interface{}) {
	for i, e := range c {
		switch e.(type) {
		case float64:
			c[i] = e.(float64) + d[i].(float64)
		case *Category:
			e.(*Category).add(d[i].(*Category))
		default:
			panic("unknown type incremententing centroid")
		}
	}
}

func centerCentroid(c []interface{}, l int) {
	for i, e := range c {
		switch e.(type) {
		case float64:
			c[i] = e.(float64) / float64(l)
		case *Category:
			e.(*Category).mean(l)
		default:
			panic("unknown type zeroing centroid")
		}
	}
}

// FIXME TESTME
func moveCentroids(centroids [][]interface{}, dataMap []Point, data Table) {
	// Maps number of elements that belongs to a centroid.
	eleMap := map[int]int{}
	for _, p := range dataMap {
		eleMap[p.K]++
	}
	for k := 0; k < len(centroids); k++ {
		if eleMap[k] == 0 {
			stracer.Traceln(k, "is a zero element centroid, not zeroing")
			continue
		}
		zeroCentroid(centroids[k])
	}
	for i, p := range dataMap {
		incrementCentroid(centroids[p.K], data.Row(i))
	}
	for k := 0; k < len(centroids); k++ {
		if eleMap[k] == 0 {
			continue
		}
		centerCentroid(centroids[k], eleMap[k])
	}
}

// Kmc computes k means clustering.
//
// Data MUST be normalized before to be passed, Normalize function could be used.
func Kmc(data Table, k int, weights []float64) (result *KmcResult, er error) {
	// FIXME randomly centroids with zero elemtns are created which take to higher SSE.
	// FIXME check for unnormalized data!
	// This assigns all elements to centroid 0 as default.
	result = &KmcResult{}
	dataMap := make([]Point, data.Len())
	// Set max distance for all elements.
	for i := 0; i < data.Len(); i++ {
		dataMap[i].Distance = 1
	}
	centroids := make([][]interface{}, 0)
	centroids, er = createRandomCentroids(k, data.Row(0))
	if er != nil {
		return
	}
	changed := true
	for {
		changed = false
		for i := 0; i < data.Len(); i++ {
			e := data.Row(i)
			for j := 0; j < k; j++ {
				var d float64
				d, er = distance(e, centroids[j], nil)
				if er != nil {
					return
				}
				if d < dataMap[i].Distance {
					changed = true
					dataMap[i].Distance = d
					dataMap[i].K = j
				}
			}
		}
		if !changed {
			break
		}
		moveCentroids(centroids, dataMap, data)
		stracer.Traceln("centroids moved", centroids)
	}
	for _, p := range dataMap {
		result.TotalSSE += math.Pow(p.Distance, 2)
	}
	result.Map = dataMap
	result.Centroids = centroids
	return
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
