// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"encoding/csv"
	"io"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
)

type featureType uint8

const (
	_ featureType = iota
	floatFeature
	stringFeature
)

// Normalize uses Table's Update()
// to modify rows with normalized values.
//
// Numerical features are normalized
// with the formula:
//
//	x - mu
//	------
//	sigma
//
// If mu, sigma or catSet are nil
// they are calculated and returned
// otherwise their computation is skipped
// and passed values are used.
//
// Categorical features are mapped to
// a representation suitable from
// other functions in the package.
func Normalize(data Table, mu, sigma []float64, catSet []string) ([]float64, []float64, []string, error) {
	nRows, nColumns := data.Caps()
	sum := make([]float64, nColumns)
	if mu == nil || sigma == nil || catSet == nil {
		catSet = make([]string, 0)
		mapSet := make(map[string]struct{})
		mu = make([]float64, nColumns)
		sigma = make([]float64, nColumns)
		for i := 0; i < nRows; i++ {
			row, err := data.Row(i)
			if err != nil {
				return nil, nil, nil, err
			}
			for i, e := range row {
				switch v := e.(type) {
				case float64:
					sum[i] += v
				case string:
					// NaN is used to indicate
					// to not compute mean
					// for the feature.
					sum[i] = math.NaN()
					mapSet[v] = struct{}{}
				default:
					return nil, nil, nil, unknownTypeErr(e)
				}
			}
		}
		// Compute means.
		for i, s := range sum {
			mu[i] = s
			if !math.IsNaN(s) {
				mu[i] = s / float64(nRows)
			}
		}
		// Compute standard deviation.
		for i := 0; i < nRows; i++ {
			row, err := data.Row(i)
			if err != nil {
				return nil, nil, nil, err
			}
			for i, e := range row {
				switch v := e.(type) {
				case float64:
					sigma[i] += math.Pow(v-mu[i], 2)
				default:
					// NaN is used to indicate
					// to not compute sigma
					// for the feature.
					sigma[i] = math.NaN()
				}
			}
		}
		for i, s := range sigma {
			sigma[i] = s
			if !math.IsNaN(s) {
				sigma[i] = math.Sqrt(s / float64(nRows-1))
			}
		}
		// Order map set.
		catSet = orderMapSet(mapSet)
	}
	// Normalize.
	for i := 0; i < nRows; i++ {
		row, err := data.Row(i)
		if err != nil {
			return nil, nil, nil, err
		}
		for i, e := range row {
			switch v := e.(type) {
			case float64:
				row[i] = (v - mu[i]) / sigma[i]
			case string:
				row[i] = newCategory(v, catSet)
			default:
				return nil, nil, nil, unknownTypeErr(e)
			}
		}
		data.Update(i, row)
	}
	return mu, sigma, catSet, nil
}

// ReadAllCSV read whole file and load it
// in memory.
func ReadAllCSV(path string) (Table, error) {
	// FIXME category must be exposed
	// to let Table to be used externally.
	// FIXME is it possible to preallocate length
	// having namers of file's rows?
	var dataSlice MemoryTable = [][]interface{}{}
	iRow := []interface{}{}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		// FIXME really needed?
		cleanStrings(row)
		iRow = make([]interface{}, len(row))
		// Elements in rows
		// are either float or string.
		for i, e := range row {
			switch kind(e) {
			case floatFeature:
				f, err := strconv.ParseFloat(e, 64)
				if err != nil {
					return nil, err
				}
				iRow[i] = f
			case stringFeature:
				// This must be normalized
				// with Normalize.
				iRow[i] = e
			default:
				return nil, unknownTypeErr(e)
			}
		}
		dataSlice = append(dataSlice, iRow)
	}
	return dataSlice, nil
}

// kind identifies data type
// from string for storing into a Go type.
func kind(s string) featureType {
	// FIXME avoid outer call to ParseFloat
	// returning also the number
	_, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return floatFeature
	}
	return stringFeature
}

func cleanStrings(row []string) {
	cleanerRgxp := regexp.MustCompile(`[[:space:]]`)
	for i, s := range row {
		row[i] = cleanerRgxp.ReplaceAllString(s, "")
	}
}

// orderCatSet returns an ordered
// slice of categories set.
func orderMapSet(set map[string]struct{}) []string {
	l := len(set)
	ss := make([]string, 0, l)
	for k := range set {
		ss = append(ss, k)
	}
	sort.Sort(sort.StringSlice(ss))
	return ss
}
