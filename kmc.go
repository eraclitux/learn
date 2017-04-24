// Copyright (c) 2015 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/eraclitux/trace"
)

// category models a categorical (aka nominal es choices A,B etc) feature.
// Choices must be translated to the form:
//	"[1,0,1]"
type category struct {
	data     uint
	choicesN uint
}

// choices is in the form "[1,0,1]"
func newCategory(choices string) *category {
	choices = checkerRgxp.ReplaceAllString(choices, "$1")
	l := uint(len(strings.Split(choices, `,`)))
	s := strings.Replace(choices, ",", "", -1)
	s = strings.Replace(s, " ", "", -1)
	if s == "" {
		s = "0"
	}
	data, err := strconv.ParseUint(s, 2, 32)
	if err != nil {
		// Fail fast.
		panic(fmt.Sprintf("in newCategory: %s", err))
	}
	return &category{
		data:     uint(data),
		choicesN: l,
	}
}

func (c *category) add(b *category) {
	c.data += b.data
}

func (c *category) zero() {
	c.data = 0
}

// mean calculates mean for an element of
// a centroid previously incremented l times.
// TODO test for overflow, if 0b0000 & 0b111110000 != 0
func (c *category) mean(l int) {
	c.data = c.data / uint(l)
}

// distance returns simple matching distance from the passed Category.
// Returning value is ∈ [0,1].
func (c *category) distance(b *category) float64 {
	return float64(hammingD(c.data, b.data)) / float64(c.choicesN)
}
func (c *category) String() string {
	format := fmt.Sprintf("%%0%db", c.choicesN)
	return fmt.Sprintf(format, c.data)
}

// BUG(eraclitux): randomly returns same
// category in tests.
func createRandCategory(l uint) *category {
	sS := []string{}
	for i := 0; i < int(l); i++ {
		sN := strconv.Itoa(rand.Intn(2))
		sS = append(sS, sN)
	}
	return newCategory(strings.Join(sS, ","))
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
				c[i] = createRandCategory(e.(*category).choicesN)
			case string:
				c[i] = ""
			default:
				return nil, unknownType(e)
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
			trace.Println(k, "is a zero element centroid, not zeroing")
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
	// FIXME check for not normalized data!

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
		trace.Println("centroids moved", centroids)
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
