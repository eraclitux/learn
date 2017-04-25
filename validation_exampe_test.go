// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn_test

import (
	"fmt"
	"log"

	"github.com/eraclitux/learn"
)

func ExampleValidate() {
	// Cross validation
	trainSet, err := learn.ReadAllCSV("datasets/iris_train.csv")
	if err != nil {
		log.Fatal(err)
	}
	mu, sigma, err := learn.Normalize(trainSet, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	testSet, err := learn.ReadAllCSV("datasets/iris_test.csv")
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = learn.Normalize(testSet, mu, sigma)
	if err != nil {
		log.Fatal(err)
	}
	clf := learn.NewkNN(trainSet, 3)
	predictedLabels, err := clf.Predict(testSet)
	if err != nil {
		log.Fatal(err)
	}
	confMatrix, err := learn.ConfusionM(testSet, predictedLabels)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(confMatrix)
	report := learn.Validate(confMatrix)
	fmt.Println(report)

	// OUTPUT:
	//             setosa(1):           5           0           0
	//         versicolor(2):           0           7           0
	//          virginica(3):           0           0           3
	//
	//      feature | precision | recall |
	//       setosa |      1.00 |   1.00 |
	//   versicolor |      1.00 |   1.00 |
	//    virginica |      1.00 |   1.00 |
	// Overall accuracy: 1.00
}
