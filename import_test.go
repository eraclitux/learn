// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"math"
	"testing"

	"github.com/eraclitux/trace"
)

func TestNormalize(t *testing.T) {
	var testCase MemoryTable = [][]interface{}{
		[]interface{}{
			400.31,
			1.2,
			newCategory("[1,0,0,0]"),
		},
		[]interface{}{
			300.21,
			10.2,
			newCategory("[0,1,0,0]"),
		},
		[]interface{}{
			-600.54,
			-2.5,
			newCategory("[0,0,0,1]"),
		},
	}
	// mean: 33.3266666666667 sigma: 551.221566916003
	// mean: 2.96666666666667 sigma: 6.53171748725657
	expectedMean := []float64{33.3266666666667, 2.96666666666667, math.NaN()}
	expectedSigma := []float64{551.221566916003, 6.53171748725657, math.NaN()}
	var expected MemoryTable = [][]interface{}{
		[]interface{}{
			0.665763742493870,
			-0.270475058070629,
			newCategory("[1,0,0,0]"),
		},
		[]interface{}{
			0.484167074279229,
			1.10741674719484,
			newCategory("[0,1,0,0]"),
		},
		[]interface{}{
			-1.14993081677310,
			-0.836941689124212,
			newCategory("[0,0,0,1]"),
		},
	}
	mu, sigma, err := Normalize(testCase, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	for i, row := range testCase {
		for j, e := range row {
			switch e.(type) {
			case float64:
				f := e.(float64)
				g := expected[i][j].(float64)
				if !floatsAreEqual(f, g) {
					t.Errorf("%10s: %v\n %10s: %+v", "expected", g, "got", f)
				}
			}
		}
	}
	for i, u := range mu {
		if !floatsAreEqual(u, expectedMean[i]) {
			t.Errorf("expected mu: %v, got: %v", expectedMean[i], u)
		}
		trace.Println("IsNaN:", math.IsNaN(u))
		if !floatsAreEqual(sigma[i], expectedSigma[i]) {
			t.Errorf("expected sigma: %v, got: %v", expectedSigma[i], sigma[i])
		}
	}
}
