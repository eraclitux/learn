// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn_test

import (
	"fmt"
	"log"

	"github.com/eraclitux/learn"
)

func ExampleNewkNN() {
	trainSet, err := learn.ReadAllCSV("datasets/iris.csv")
	if err != nil {
		log.Fatal(err)
	}
	mu, sigma, catSet, err := learn.Normalize(trainSet, nil, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	clf, err := learn.NewkNN(trainSet, 3)
	if err != nil {
		log.Fatal(err)
	}
	// Categorize single sample.
	var testSet learn.MemoryTable = make([][]interface{}, 1)
	testSet[0] = []interface{}{5.2, 3.4, 1.3, 0.1}
	_, _, _, err = learn.Normalize(testSet, mu, sigma, catSet)
	if err != nil {
		log.Fatal(err)
	}
	prediction, err := clf.Predict(testSet)
	if err != nil {
		log.Fatal(err)
	}
	r, err := prediction.Row(0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("predicted category:", r[0])

	// OUTPUT:
	// predicted category: setosa
}
