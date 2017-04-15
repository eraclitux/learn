// Copyright (c) 2015 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn_test

import (
	"fmt"
	"log"

	"github.com/eraclitux/learn"
)

func ExampleNewLinearRegression() {
	trainData, err := learn.ReadAllCSV("datasets/linear_test.csv")
	if err != nil {
		log.Fatal(err)
	}
	var tab learn.MemoryTable = make([][]interface{}, 1)
	// no need to normalize as normal equation
	// is used.
	tab[0] = []interface{}{1650.0, 3.0}
	lr, err := learn.NewLinearRegression(trainData)
	if err != nil {
		log.Fatal(err)
	}
	y, err := lr.Predict(tab)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("predicted price for a (%.f sq-ft, %.f rooms) house: $%.f", tab[0][0], tab[0][1], y[0])
	// Output:
	// predicted price for a (1650 sq-ft, 3 rooms) house: $293081
}
