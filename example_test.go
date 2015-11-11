// Copyright (c) 2015 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package nml_test

import (
	"fmt"
	"log"

	"github.com/eraclitux/nml"
)

func ExampleKmc() {
	rC, er := nml.LoadCSV("./datasets/iris_nolabels.csv")
	if er != nil {
		return
	}
	// NOTE this loads all data in memory.
	data, er := nml.Normalize(rC)
	if er != nil {
		return
	}
	result, err := nml.Kmc(data, 3, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
