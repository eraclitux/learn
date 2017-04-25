// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"reflect"
	"testing"
)

func TestKSamples_CheckUpdate(t *testing.T) {
	samples := newKSamples(3)
	expectedSamples := newKSamples(3)
	row := []interface{}{0.2, 0.4, 0.1, "one"}
	samples.checkUpdate(0.9, row)
	row = []interface{}{0.2, 0.4, 0.1, "one"}
	samples.checkUpdate(0.8, row)
	row = []interface{}{0.1, 0.3, 0.0, "two"}
	samples.checkUpdate(0.4, row)
	row = []interface{}{0.1, 0.3, 0.0, "two"}
	samples.checkUpdate(0.3, row)
	// checkUpdate does not enforce any order,
	// we specify a knowed order for expectedSamples
	row = []interface{}{0.1, 0.3, 0.0, "two"}
	expectedSamples[0] = kSample{
		distance: 0.3,
		row:      row,
	}
	row = []interface{}{0.2, 0.4, 0.1, "one"}
	expectedSamples[1] = kSample{
		distance: 0.8,
		row:      row,
	}
	row = []interface{}{0.1, 0.3, 0.0, "two"}
	expectedSamples[2] = kSample{
		distance: 0.4,
		row:      row,
	}

	if !reflect.DeepEqual(samples, expectedSamples) {
		t.Fatalf("expected: %v got: %v", expectedSamples, samples)
	}
}

func TestKSamples_GetNearest(t *testing.T) {
	samples := newKSamples(5)
	row := []interface{}{0.2, 0.4, 0.1, "one"}
	samples.checkUpdate(0.5, row)
	row = []interface{}{0.1, 0.3, 0.0, "one"}
	samples.checkUpdate(0.5, row)
	row = []interface{}{0.1, 0.3, 0.0, "two"}
	samples.checkUpdate(0.4, row)
	row = []interface{}{0.14, 0.33, 0.23, "two"}
	samples.checkUpdate(0.6, row)
	row = []interface{}{0.13, 0.33, 0.23, "two"}
	samples.checkUpdate(0.6, row)

	nearest := samples.getNearest()

	if nearest != "two" {
		t.Fatal("nearest:", nearest)
	}
}
