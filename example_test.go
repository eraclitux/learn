// Copyright (c) 2015 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package sml_test

import (
	"fmt"
	"log"

	"github.com/eraclitux/sml"
)

func ExampleKmc() {
	rC, er := sml.LoadCSV("./datasets/iris_nolabels.csv")
	if er != nil {
		return
	}
	// Load all data in memory.
	data, er := sml.Normalize(rC)
	if er != nil {
		return
	}
	result, err := sml.Kmc(data, 3, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
