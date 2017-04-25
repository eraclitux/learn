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
	"strconv"

	"github.com/eraclitux/trace"
)

type featureType uint8

const (
	_ featureType = iota
	cat
	float
	str
)

var checkerRgxp = regexp.MustCompile(`\[(.+)\]`)

// Normalize uses Table's Update()
// to modify rows with normalized values.
//
// Only numerical features are normalized
// with the formula:
//
//	x - mu
//	------
//	sigma
//
// If mu and sigma are nil then
// they are calculated and returned
// otherwise their calculation is skipped
// and passed values are used, NaN is returned for
// non numerical features.
func Normalize(data Table, mu, sigma []float64) ([]float64, []float64, error) {
	nRows, nColumns := data.Caps()
	sum := make([]float64, nColumns)
	if mu == nil || sigma == nil {
		mu = make([]float64, nColumns)
		sigma = make([]float64, nColumns)
		for i := 0; i < nRows; i++ {
			row, err := data.Row(i)
			if err != nil {
				return nil, nil, err
			}
			for i, e := range row {
				switch e.(type) {
				case float64:
					sum[i] += e.(float64)
				default:
					// NaN is used to indicate
					// an unprocessed feature.
					sum[i] = math.NaN()
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
				return nil, nil, err
			}
			for i, e := range row {
				switch e.(type) {
				case float64:
					sigma[i] += math.Pow(e.(float64)-mu[i], 2)
				default:
					// NaN is used to indicate
					// an unprocessed feature.
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
	}
	// Normalize.
	for i := 0; i < nRows; i++ {
		row, err := data.Row(i)
		if err != nil {
			return nil, nil, err
		}
		for i, e := range row {
			switch e.(type) {
			case float64:
				row[i] = (e.(float64) - mu[i]) / sigma[i]
			}
		}
		data.Update(i, row)
	}
	return mu, sigma, nil
}

// ReadAllCSV read whole file and load it
// in memory.
func ReadAllCSV(path string) (Table, error) {
	// FIXME test it!
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
		cleanStrings(row)
		iRow = make([]interface{}, len(row))
		for i, e := range row {
			switch kind(e) {
			case float:
				f, err := strconv.ParseFloat(e, 64)
				if err != nil {
					return nil, err
				}
				iRow[i] = f
			case cat:
				iRow[i] = newCategory(e)
			case str:
				iRow[i] = e
			default:
				// FIXME return error!
				panic("unknown type normalizing data")

			}
		}
		dataSlice = append(dataSlice, iRow)
	}
	return dataSlice, nil
}

// kind tries to identify data type
// from string for storing it into a Go type.
func kind(s string) featureType {
	// FIXME to not recall conversion methods
	// in the caller. Return also converted values.
	if checkerRgxp.MatchString(s) {
		trace.Println("clenerRgxp matched:", s)
		return cat
	}
	_, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return float
	}
	return str
}

func cleanStrings(row []string) {
	cleanerRgxp := regexp.MustCompile(`[[:space:]]`)
	for i, s := range row {
		row[i] = cleanerRgxp.ReplaceAllString(s, "")
	}
}
