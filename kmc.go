// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// BUG(eraclitux): randomly returns same
// category in tests.
func createRandCategory(l uint) *category {
	sS := []string{}
	for i := 0; i < int(l); i++ {
		sN := strconv.Itoa(rand.Intn(2))
		sS = append(sS, sN)
	}
	return newCategory(strings.Join(sS, ","), nil) // FIXME
}

// FIXME Andrew Ng suggests to initialize centroids
// to points of training samples.
func createRandomCentroids(k int, s []interface{}) ([][]interface{}, error) {
	l := len(s)
	r := make([][]interface{}, 0, k)
	for i := 0; i < k; i++ {
		c := make([]interface{}, l)
		for i, e := range s {
			switch e.(type) {
			case float64:
				c[i] = rand.Float64()
			case *category:
				c[i] = createRandCategory(e.(*category).catNumber)
			case string:
				c[i] = ""
			default:
				return nil, unknownTypeErr(e)
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
		case *category:
			c[i].(*category).zero()
		case string:
			// do nothing for string features.
		default:
			panic("unknown type zeroing centroid")
		}
	}
}

// incrementCentroid adds quantities to centroids elements
// to calculate the mean after.
func incrementCentroid(c []interface{}, d []interface{}) {
	for i, e := range c {
		switch e.(type) {
		case float64:
			c[i] = e.(float64) + d[i].(float64)
		case *category:
			e.(*category).add(d[i].(*category))
		case string:
			// do nothing for string features.
		default:
			panic("unknown type increasing centroid")
		}
	}
}

func centerCentroid(c []interface{}, l int) {
	for i, e := range c {
		switch e.(type) {
		case float64:
			c[i] = e.(float64) / float64(l)
		case *category:
			e.(*category).mean(l)
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
// Data MUST be normalized before to be passed,
// Normalize function can be used for that.
func Kmc(data Table, k int, weights []float64) (result *KmcResult, er error) {
	// FIXME randomly centroids with zero elements are created which take to higher SSE.
	nRows, _ := data.Caps()
	// This assigns all elements to centroid 0 as default.
	result = &KmcResult{}
	dataMap := make([]Point, nRows)
	// Set max distance for all elements.
	for i := 0; i < nRows; i++ {
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
		for i := 0; i < nRows; i++ {
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
			er = err
			return
		}
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
