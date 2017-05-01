// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn_test

import (
	"fmt"
	"log"

	"github.com/eraclitux/learn"
)

func ExampleKmc() {
	data, err := learn.ReadAllCSV("datasets/iris_nolabels.csv")
	if err != nil {
		log.Fatal(err)
	}
	_, _, _, err = learn.Normalize(data, nil, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	result, err := learn.Kmc(data, 3, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
