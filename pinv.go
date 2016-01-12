// Copyright (c) 2015 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package nml

import (
	"math"

	"github.com/gonum/matrix/mat64"
)

// Used to correct floating point errors.
// Everything equal or less than small+epsilon is considered zero (right?)
const epsilon float64 = 0.0000001

// pinv uses SVD to calculate pseudo inverse of
// a given matrix.
func pinv(X *mat64.Dense) {
	// Using SVD to calculate pseudo-inverse:
	// https://en.wikipedia.org/wiki/Moore%E2%80%93Penrose_pseudoinverse#Singular_value_decomposition_.28SVD.29

	// small as smallest non zero float math.SmallestNonzeroFloat64
	svd := mat64.SVD(X, epsilon, math.SmallestNonzeroFloat64, true, true)

	l := len(svd.Sigma)
	// Assemble sigma pseudo inverse matrix.
	// We get the pseudo inverse by taking the reciprocal
	// of each non-zero element on the diagonal,
	// leaving the zeros in place, and then
	// transposing the matrix.
	S := mat64.NewDense(l, l, nil)
	for i, e := range svd.Sigma {
		if svd.Sigma[i] != 0 {
			svd.Sigma[i] = 1 / e
			S.Set(i, i, svd.Sigma[i])
		}
	}
	var L, M mat64.Dense
	L.Mul(svd.V, S.T())
	M.Mul(&L, svd.U.T())
	*X = M
}
