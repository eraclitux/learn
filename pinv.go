// Copyright (c) 2015 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"github.com/gonum/matrix"
	"github.com/gonum/matrix/mat64"
)

// pinv uses SVD to calculate pseudo inverse of
// a given matrix.
func pinv(X *mat64.Dense) {
	// FIXME use interface?
	// FIXME do not modify arg but return a pointer
	// Using SVD to calculate pseudo-inverse:
	// https://en.wikipedia.org/wiki/Moore%E2%80%93Penrose_pseudoinverse#Singular_value_decomposition_.28SVD.29

	svd := new(mat64.SVD)
	_ = svd.Factorize(X, matrix.SVDFull)
	singValues := svd.Values(nil)

	l := len(singValues)
	// Assemble sigma pseudo inverse matrix.
	// We get the pseudo inverse by taking the reciprocal
	// of each non-zero element on the diagonal,
	// leaving the zeros in place, and then
	// transposing the matrix.
	S := mat64.NewDense(l, l, nil)
	for i, e := range singValues {
		if singValues[i] != 0 {
			singValues[i] = 1 / e
			S.Set(i, i, singValues[i])
		}
	}
	var U, V, L, M mat64.Dense
	U.UFromSVD(svd)
	V.VFromSVD(svd)
	L.Mul(&V, S.T())
	M.Mul(&L, U.T())
	*X = M
}
