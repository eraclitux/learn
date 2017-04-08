// Copyright (c) 2015 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"testing"

	"github.com/gonum/matrix/mat64"
)

func TestPinv(t *testing.T) {
	cases := []struct {
		test, expected []float64
	}{
		{
			// NOTE this is a singular matrix.
			[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			// Expected matrix, generated in octave
			// with pinv(X)
			[]float64{
				-6.38888888888888e-01, -1.66666666666666e-01, 3.05555555555555e-01,
				// NOTE E(2,2) == -3.608224830031759e-16, but it's ok,
				// it's caused by IEEE754 limits, EqualApprox deals with that.
				-5.55555555555556e-02, 3.81639164714898e-17, 5.55555555555556e-02,
				5.27777777777777e-01, 1.66666666666666e-01, -1.94444444444444e-01,
			},
		},
		{
			// Non singular matrix.
			[]float64{10, 2, 3, 4, 5, 6, 7, 8, 9},
			// Expected matrix, generated in octave
			// with pinv(X)
			[]float64{
				0.111111111111111, -0.222222222222223, 0.111111111111111,
				-0.222222222222222, -2.555555555555554, 1.777777777777777,
				0.111111111111111, 2.444444444444445, -1.555555555555555,
			},
		},
	}

	for _, e := range cases {
		E := mat64.NewDense(3, 3, e.expected)
		T := mat64.NewDense(3, 3, e.test)
		pinv(T)
		if !mat64.EqualApprox(T, E, 0.000001) {
			t.Error("not a (pseudo)inverse matrix")
		}
		fT := mat64.Formatted(T)
		fE := mat64.Formatted(E)
		t.Logf("have:\n%v\n", fT)
		t.Logf("want:\n%v\n", fE)
		t.Log("=========================")
	}
}
