// Copyright (c) 2015 Andrea Masi. All rights reserved.
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

// csvClosable implements ReadCloser
type csvClosable struct {
	io.Closer
	*csv.Reader
}

var checkerRgxp *regexp.Regexp = regexp.MustCompile(`\[(.+)\]`)

// Normalize mathematically normalizes data.
//
// It loads data in memory. If this is not feasible because is "Big Data"
// a cutom type which implements Table interface backed by
// disk (file,database, etc) could be used.
//
// Every value (quantitative, nominal, cardinal, binary)
// is transformed to appropriate scalar/Category
// with elements ∈ {0,1}.
//
// Example of normalized data:
//
//	Hours	Choices		Stars	Price
//	1,	"[1,0,1,0]", 	1,	1
//	0,	"[0,0,0,1]" ,	0.25,	0
//
// Normalize uses the formula:
//
//	   x - Vmin
//	-----------
//	Vmax - Vmin
//
// to normalize all dataset.
// Where (Vmax, Vmin) are that maximun and minimun values for that feature.
//
// A good reference for data normalization:
// http://people.revoledu.com/kardi/tutorial/Similarity/MutivariateDistance.html
//
// BUG(eraclitux): better to use mean normalization
//
//	x - mean(x)
//	-----------
//	Vmax - Vmin
//
func Normalize(dataReadCloser ReadCloser) (Table, error) {
	defer dataReadCloser.Close()
	var dataSlice MemoryTable = [][]interface{}{}
	iRow := []interface{}{}
	// store maxs and mins
	maxs := []interface{}{}
	mins := []interface{}{}
	var j uint64
	for {
		row, err := dataReadCloser.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if j == 0 {
			maxs = make([]interface{}, len(row))
			mins = make([]interface{}, len(row))
			for i := 0; i < len(row); i++ {
				var f float64
				maxs[i] = f
				f = math.MaxFloat64
				mins[i] = f
			}
		}
		j++
		trace.Println("row from csv:", row)
		cleanStrings(row)
		trace.Println("cleaned row:", row)
		iRow = make([]interface{}, len(row))
		for i, e := range row {
			switch kind(e) {
			case float:
				f, err := strconv.ParseFloat(e, 64)
				if err != nil {
					return nil, err
				}
				iRow[i] = f
				if f > maxs[i].(float64) {
					maxs[i] = f
				} else if f < mins[i].(float64) {
					mins[i] = f
				}
			case cat:
				iRow[i] = newCategory(e)
			case str:
				iRow[i] = e
			default:
				panic("unknown type normalizing data")

			}
		}
		dataSlice = append(dataSlice, iRow)
	}
	// Normalize
	for _, row := range dataSlice {
		for i, e := range row {
			switch e.(type) {
			case float64:
				row[i] = (e.(float64) - mins[i].(float64)) / (maxs[i].(float64) - mins[i].(float64))
			}
		}
	}
	return dataSlice, nil
}

// LoadCSV reads data from a file and returns
// an implementation of ReadCloser.
// Example of data
//
//	Hours	Choices		Stars	Price
//
//	12,	"A,C",		5,	15.10
//	1,	"D",		1,	1
//
func LoadCSV(path string) (ReadCloser, error) {
	// we cannot call
	//defer f.Close()
	// or a caller will get  "bad file descriptor" when reading
	f, err := os.Open(path)
	if err != nil {
		return csvClosable{
			nil,
			nil,
		}, err
	}
	r := csv.NewReader(f)
	return csvClosable{
		f,
		r,
	}, nil
}

// ReadAllCSV read whole file and load it
// in memory.
func ReadAllCSV(path string) (Table, error) {
	// FIXME test it!
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
	// FIXME to not recall converiosn methods
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
