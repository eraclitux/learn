// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"math"
	"testing"
)

func TestNormalize(t *testing.T) {
	set := []string{
		"foo",
		"bar",
		"crow",
	}
	var testCase MemoryTable = [][]interface{}{
		{
			400.31,
			1.2,
			"foo",
		},
		{
			300.21,
			10.2,
			"bar",
		},
		{
			-600.54,
			-2.5,
			"crow",
		},
	}
	expectedMean := []float64{33.3266666666667, 2.96666666666667, math.NaN()}
	expectedSigma := []float64{551.221566916003, 6.53171748725657, math.NaN()}
	var expected MemoryTable = [][]interface{}{
		{
			0.665763742493870,
			-0.270475058070629,
			newCategory("foo", set),
		},
		{
			0.484167074279229,
			1.10741674719484,
			newCategory("bar", set),
		},
		{
			-1.14993081677310,
			-0.836941689124212,
			newCategory("crow", set),
		},
	}
	mu, sigma, set, err := Normalize(testCase, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	for i, row := range testCase {
		for j, e := range row {
			switch v := e.(type) {
			case float64:
				g := expected[i][j].(float64)
				if !floatsAreEqual(v, g) {
					t.Errorf("%10s: %v\n %10s: %+v", "expected", g, "got", v)
				}
			}
		}
	}
	for i, u := range mu {
		if !floatsAreEqual(u, expectedMean[i]) {
			t.Errorf("expected mu: %v, got: %v", expectedMean[i], u)
		}
		if !floatsAreEqual(sigma[i], expectedSigma[i]) {
			t.Errorf("expected sigma: %v, got: %v", expectedSigma[i], sigma[i])
		}
	}
}
