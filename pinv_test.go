// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"testing"

	"github.com/gonum/matrix/mat64"
)

const epsilon = 1e-5

type pinvCase struct {
	r, c           int
	test, expected []float64
}

var pinvCases []pinvCase = []pinvCase{
	{
		// NOTE this is a singular matrix.
		3, 3,
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
		3, 3,
		[]float64{10, 2, 3, 4, 5, 6, 7, 8, 9},
		// Expected matrix, generated in octave
		// with pinv(X)
		[]float64{
			0.111111111111111, -0.222222222222223, 0.111111111111111,
			-0.222222222222222, -2.555555555555554, 1.777777777777777,
			0.111111111111111, 2.444444444444445, -1.555555555555555,
		},
	},
	{
		// Rectangular matrix.
		4, 3,
		[]float64{8, 1, 6, 3, 5, 7, 4, 9, 2, 9, 8, 7},
		// Expected matrix, generated in octave
		// with pinv(X)
		[]float64{
			0.10602310231023103, -0.15621562156215615, 0.02268976897689765, 0.05885588558855884,
			-0.07343234323432345, 0.01870187018701872, 0.09323432343234321, 0.01760176017601759,
			-0.00288778877887791, 0.19361936193619353, -0.08622112211221114, -0.02365236523652363,
		},
	},
}

func TestPinv(t *testing.T) {
	for _, c := range pinvCases {
		E := mat64.NewDense(c.c, c.r, c.expected)
		T := mat64.NewDense(c.r, c.c, c.test)
		R, err := pinv(T)
		if err != nil {
			t.Error(err)
		}
		if !mat64.EqualApprox(R, E, epsilon) {
			t.Error("not a (pseudo)inverse matrix")
			fT := mat64.Formatted(T)
			fR := mat64.Formatted(R)
			fE := mat64.Formatted(E)
			t.Logf("case:\n%v\n", fT)
			t.Logf("have:\n%v\n", fR)
			t.Logf("want:\n%v\n", fE)
			t.Log("=========================")
		}
	}
}

// Test that pinv(A) satisfies
// Moore–Penrose pseudoinverse properties:
// A x pinv(A) x A = A
func TestPinv_properties(t *testing.T) {
	for _, c := range pinvCases {
		A := mat64.NewDense(c.r, c.c, c.test)
		P := new(mat64.Dense)
		Ap, err := pinv(A)
		if err != nil {
			t.Error(err)
		}
		P.Product(A, Ap, A)
		if !mat64.EqualApprox(A, P, epsilon) {
			t.Error("first property of Moore–Penrose violated")
			fA := mat64.Formatted(A)
			t.Logf("case:\n%v\n", fA)
		}
	}
}
