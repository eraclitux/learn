// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"errors"
	"fmt"
)

// ErrNoData should be returned
// in case of problems retrieving data
// from in Table's underlying storage.
var ErrNoData = errors.New("learn: no data")

// Point stores data about
// kmc's points.
type Point struct {
	// The index of centroid to which point belongs.
	K int
	// Distance from centroid.
	Distance float64
}

// KmcResult stores result
// of k mean clustering.
// FIXME divide TotalSSE for number of samples
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
	// Returns rows
	// and columns numbers.
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

// Caps implements Table's Caps.
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

// Row implements Table's Row.
func (t MemoryTable) Row(i int) ([]interface{}, error) {
	if i >= len(t) {
		return nil, ErrNoData
	}
	e := t[i]
	return e, nil
}

// Update implements Table's Update.
func (t MemoryTable) Update(i int, r []interface{}) error {
	if i >= len(t) {
		return ErrNoData
	}
	t[i] = r
	return nil
}
