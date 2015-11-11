// Copyright (c) 2015 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

// Package mml exposes some machine learning alghoritms, in a minimalistic way.
//
// It tries to be as idiomatic as possible. Interfaces are used in public APIs when possible
// to make methods adaptable to custom needs.
// Table interface should makes (hopefully) easy to use storage other than memory
// when dealing with "Big Data" (database, filesystem etc..)
//
// Clustering:
//
//	- k means clustering
//
// Classification:
//
//	- kNN
//
// K means clustering
//
// Categorical and numeriacl features are supported.
//
// Method for Distance calculation is automatically
// choosed at runtime:
//
// - manhattan for numerical features
//
// - humming distance for categorical features
//
// Example of data
//
//	Hours	Choices		Stars	Price
//	12,	"A,C",		5,	15.10
//	1,	"D"		1,	1
//
// Categorical features must be translated to an array of 0 and 1:
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
			case string:
				c[i] = ""
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
		case string:
			// do nothing for string features.
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
		case string:
			// do nothing for string features.
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
		case string:
			// do nothing for string features.
		default:
			panic("unknown type centering centroid")
		}
	}
}

// FIXME TESTME
func moveCentroids(centroids [][]interface{}, dataMap []Point, data Table) error {
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
		row, err := data.Row(i)
		if err != nil {
			return err
		}
		incrementCentroid(centroids[p.K], row)
	}
	for k := 0; k < len(centroids); k++ {
		if eleMap[k] == 0 {
			continue
		}
		centerCentroid(centroids[k], eleMap[k])
	}
	return nil
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
	row, er := data.Row(0)
	if er != nil {
		return

	}
	centroids, er = createRandomCentroids(k, row)
	if er != nil {
		return
	}
	changed := true
	for {
		changed = false
		for i := 0; i < data.Len(); i++ {
			e, err := data.Row(i)
			if err != nil {
				er = err
				return
			}
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
		err := moveCentroids(centroids, dataMap, data)
		if err != nil {
			err = err
			return
		}
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
