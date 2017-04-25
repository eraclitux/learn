// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"errors"
	"fmt"
)

// NoRow error should be returned in Table's Row
// in case of problems retrieving data from underlying
// storage.
var NoRow error = errors.New("learn: no row with this index")

var NoData error = errors.New("learn: no data")

type Point struct {
	// The index of centroid to which point belongs.
	K int
	// Distance from centroid.
	Distance float64
}

// BUG(eraclitux): divide TotalSSE for number of samples
// to have a smaller number.
type KmcResult struct {
	Map       []Point
	Centroids [][]interface{}
	// Sum of squared errors
	TotalSSE float64
}

func (r *KmcResult) String() string {
	return fmt.Sprintf(
		"%d clusters, total SSE: %f",
		len(r.Centroids), r.TotalSSE,
	)
}

// Table models tabular data.
type Table interface {
	// Returns total elements.
	Caps() (int, int)
	// Returns i-th row.
	Row(int) ([]interface{}, error)
	Update(int, []interface{}) error
	//Headers() []string
	// Maybe useful in future:
	// for d.Next {d.Row()}
	//Row() []interface{}
	//Next() bool
	// Useful?
	//NFeatures() // returns number of features?
}

// MemoryTable is a Table that stores data in memory.
type MemoryTable [][]interface{}

func (t MemoryTable) Caps() (int, int) {
	var rows, colums int
	if t != nil {
		rows = len(t)
	}
	if t[0] != nil {
		colums = len(t[0])
	}
	return rows, colums
}
func (t MemoryTable) Len() int { return len(t) }
func (t MemoryTable) Row(i int) ([]interface{}, error) {
	if i >= len(t) {
		return nil, NoRow
	}
	e := t[i]
	return e, nil
}
func (t MemoryTable) Update(i int, r []interface{}) error {
	if i >= len(t) {
		return NoRow
	}
	t[i] = r
	return nil
}
