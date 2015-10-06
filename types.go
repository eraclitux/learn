// Copyright (c) 2015 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package sml

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// UnknownType is returned in case a feature is an
// unrecognized type.
var UnknownType error = errors.New("unknown type")

type Point struct {
	// The index of centroid to which point belongs.
	K int
	// Distance from centroid.
	Distance float64
}

// Category models a categorical (aka nominal es choices A,B etc) feature.
type Category struct {
	data     uint
	choicesN uint
}

// choices is in the form "[1,0,1]"
func newCategory(choices string) *Category {
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
	return &Category{
		data:     uint(data),
		choicesN: l,
	}
}

func (c *Category) add(b *Category) {
	c.data += b.data
}

func (c *Category) zero() {
	c.data = 0
}

// mean calculates mean for an element of
// a centroid previously incremented l times.
// TODO test for overflow, if 0b0000 & 0b111110000 != 0
func (c *Category) mean(l int) {
	c.data = c.data / uint(l)
}

// Distance returns simple matching distance from the passed Category.
// Returning value is ∈ [0,1].
func (c *Category) distance(b *Category) float64 {
	return float64(hummingD(c.data, b.data)) / float64(c.choicesN)
}
func (c *Category) String() string {
	format := fmt.Sprintf("%%0%db", c.choicesN)
	return fmt.Sprintf(format, c.data)
}

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

// ReadCloser makes possible to close underling
// file descriptor in a caller.
//
// Read() method must returns a io.EOF when there is no more
// data to parse.
type ReadCloser interface {
	Read() ([]string, error)
	Close() error
}

// Table models tabular data.
type Table interface {
	// Returns total elements.
	Len() int
	// Returns i-th row.
	Row(i int) []interface{}
	// Maybe usefull it future:
	// for d.Next {d.Row()}
	//Row() []interface{}
	//Next() bool
}

type memoryTable [][]interface{}

func (t memoryTable) Len() int                { return len(t) }
func (t memoryTable) Row(i int) []interface{} { return t[i] }
