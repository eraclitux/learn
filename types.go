// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"errors"
	"fmt"
)

// ErrNoData is returned
// in case of problems retrieving data
// from Table's underlying storage.
var ErrNoData = errors.New("learn: no data")

// Point stores data about
// kmc's points.
type Point struct {
	K        int     // The index of centroid to which point belongs.
	Distance float64 // Distance from centroid.
}

// KmcResult stores result
// of k mean clustering.
// FIXME divide TotalSSE for number of samples
// to have a smaller number.
type KmcResult struct {
	Map       []Point
	Centroids [][]interface{}
	TotalSSE  float64 // Sum of squared errors
}

func (r *KmcResult) String() string {
	return fmt.Sprintf(
		"%d clusters, total SSE: %f",
		len(r.Centroids), r.TotalSSE,
	)
}

// Table models tabular data.
type Table interface {
	Caps() (int, int)                    // Returns rows and columns numbers.
	Row(i int) ([]interface{}, error)    // Returns i-th row.
	Update(i int, r []interface{}) error // Substitutes i-th row with r.
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
