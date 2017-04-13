// Copyright (c) 2015 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"encoding/csv"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

// height,weigth,choices
var dataCSV = []byte(
	`100.34,23,"[1,0,0,0]"
	10.4,3,"[0,1,0,0]"
	400.4, -67,"[0,0,0,1]"`)

func TestLoadCSV(t *testing.T) {
	tempPath := os.TempDir() + "/learn_test.csv"
	err := ioutil.WriteFile(tempPath, dataCSV, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Remove(tempPath)
	}()
	r, _ := LoadCSV(tempPath)
	if r == nil {
		t.Fatal("empty csv reader")
	}

}

func TestNormalize(t *testing.T) {
	// [0.2306153846153846 1 1000] [0 0.7777777777777778 0100] [1 0 0001]
	var expected MemoryTable = [][]interface{}{
		[]interface{}{
			float64(0.2306153846153846),
			float64(1),
			newCategory("[1,0,0,0]"),
		},
		[]interface{}{
			float64(0),
			float64(0.7777777777777778),
			newCategory("[0,1,0,0]"),
		},
		[]interface{}{
			float64(1),
			float64(0),
			newCategory("[0,0,0,1]"),
		},
	}
	f, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	rC := csvClosable{f, csv.NewReader(strings.NewReader(string(dataCSV)))}
	defer rC.Close()
	data, err := Normalize(rC)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(data, expected) {
		t.Fatalf("expected: %v got: %v", expected, data)
		t.Log("normalized data:", data)
	}
}
