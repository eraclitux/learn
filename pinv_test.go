// Copyright (c) 2015 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package nml

import (
	"testing"

	"github.com/gonum/matrix/mat64"
)

func TestPinv(t *testing.T) {
	// NOTE this is a singular matrix.
	d := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	D := mat64.NewDense(3, 3, d)
	// Expected matrix, generatete in octave
	// with pinv(D)
	e := []float64{
		-6.38888888888888e-01, -1.66666666666666e-01, 3.05555555555555e-01,
		// NOTE D(2,2) -3.608224830031759e-16, but its ok,
		// it's caused by IEEE754 limits, EqualApprox deal with that.
		-5.55555555555556e-02, 3.81639164714898e-17, 5.55555555555556e-02,
		5.27777777777777e-01, 1.66666666666666e-01, -1.94444444444444e-01,
	}

	E := mat64.NewDense(3, 3, e)
	pinv(D)
	if !mat64.EqualApprox(D, E, 0.00001) {
		t.Log(D)
		t.Log(E)
		t.Fatal("not a (pseudo)inverse matrix...")
	}

	// Non singular matrix.
	d = []float64{10, 2, 3, 4, 5, 6, 7, 8, 9}
	D = mat64.NewDense(3, 3, d)
	// Expected matrix, generatete in octave
	// with pinv(D)
	e = []float64{
		0.111111111111111, -0.222222222222223, 0.111111111111111,
		-0.222222222222222, -2.555555555555554, 1.777777777777777,
		0.111111111111111, 2.444444444444445, -1.555555555555555,
	}

	E = mat64.NewDense(3, 3, e)
	pinv(D)
	if !mat64.EqualApprox(D, E, 0.00001) {
		t.Log(D)
		t.Log(E)
		t.Fatal("not a (pseudo)inverse matrix...")
	}
}
