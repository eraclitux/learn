// Copyright (c) 2015 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"errors"

	"github.com/gonum/matrix"
	"github.com/gonum/matrix/mat64"
)

// pinv uses SVD to calculate pseudo inverse of
// a given matrix.
func pinv(X *mat64.Dense) (mat64.Matrix, error) {
	// Using SVD to calculate pseudo-inverse:
	// https://en.wikipedia.org/wiki/Moore%E2%80%93Penrose_pseudoinverse#Singular_value_decomposition_.28SVD.29
	r, c := X.Dims()
	svd := new(mat64.SVD)
	ok := svd.Factorize(X, matrix.SVDFull)
	if !ok {
		return nil, errors.New("learn: not factorizable")
	}
	singValues := svd.Values(nil)
	// Assemble sigma pseudo inverse matrix.
	// We get the pseudo inverse by taking the reciprocal
	// of each non-zero element on the diagonal,
	// leaving the zeros in place, and then
	// transposing the matrix.
	Σi := mat64.NewDense(r, c, nil)
	for i, e := range singValues {
		if singValues[i] != 0 {
			singValues[i] = 1 / e
			Σi.Set(i, i, singValues[i])
		}
	}
	var U, V, P mat64.Dense
	U.UFromSVD(svd)
	V.VFromSVD(svd)
	P.Product(&V, Σi.T(), U.T())
	return &P, nil
}
