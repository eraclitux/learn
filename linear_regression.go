// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"errors"

	"github.com/gonum/matrix/mat64"
)

// Regression models a regression
// problem.
type Regression interface {
	Predict(Table) ([]float64, error)
}

type linearRegression struct {
	theta mat64.Matrix
}

// Predict given a Table with samples in its rows:
// 	x1 x2 ... xn
//	...
// 	x1 x2 ... xn
// returns a slice of float of estimated y values.
//
// Features must be stored as float64 in Table.
func (lr *linearRegression) Predict(t Table) ([]float64, error) {
	// FIXME Table confuses here? use []float64
	// and return a single float64?
	m, _ := t.Caps()
	var y mat64.Dense
	ys := make([]float64, m)
	for i := 0; i < m; i++ {
		row, err := t.Row(i)
		if err != nil {
			return nil, err
		}
		// +1 for x0
		n := len(row) + 1
		a := make([]float64, n)
		// for x0
		a[0] = 1
		for i, e := range row {
			f, ok := e.(float64)
			if !ok {
				return nil, unknownTypeErr(f)
			}
			a[i+1] = f
		}
		X := mat64.NewDense(1, n, a)
		// y is a (1 x 1) matrix
		y.Mul(X, lr.theta)
		//fT := mat64.Formatted(X)
		//trace.Printf("X:\n%v\n", fT)
		ys[i] = y.At(0, 0)
	}
	return ys, nil
}

// NewLinearRegression returns Regression type
// for linear regression.
//
// Data is a Table with training samples as rows.
// Last element in the row MUST be
// the observed value of dependent variable y.
//
// Current implementation uses normal equation,
// data normalization is not necessary.
//
// Table will be loaded in memory.
func NewLinearRegression(Data Table) (Regression, error) {
	// m: number of samples
	// n: number of features
	// X is a representation of design matrix,
	// an (m x n+1) matrix with cases as rows
	// where X(i,1) is always 1.
	//	1 x1 ... xn+1
	//	...
	//	1 x1 ... xn+1
	// The 'ones column' is added.
	// Y is the vector of observed values
	// of dependent variable.
	var X, Y *mat64.Dense
	m, _ := Data.Caps()
	if m <= 0 {
		return nil, ErrNoData
	}
	s1, err := Data.Row(0)
	if err != nil {
		return nil, ErrNoData
	}
	// o = n+1 because table stores
	// values of y in the last column.
	// +1 for theta0 == 1 terms.
	o := len(s1)
	// Build design matrix and y vector.
	X = mat64.NewDense(m, o, nil)
	yRows := make([]float64, m)
	for i := 0; i < m; i++ {
		// Row for design matrix.
		row := make([]float64, o)
		r, err := Data.Row(i)
		if err != nil {
			return nil, ErrNoData
		}
		// x0
		row[0] = 1
		for j, e := range r {
			if f, ok := e.(float64); !ok {
				return nil, unknownTypeErr(e)
			} else {
				if j >= o {
					return nil, errors.New("index out of range")
				} else if j == o-1 {
					// y element
					yRows[i] = f
				} else {
					row[j+1] = f
				}
			}
		}
		X.SetRow(i, row)
	}
	Y = mat64.NewDense(m, 1, yRows)
	// theta = pinv(X'*X)*X'*y
	P := new(mat64.Dense)
	Theta := new(mat64.Dense)
	// (n x n)
	P.Mul(X.T(), X)
	Pi, err := pinv(P)
	if err != nil {
		return nil, err
	}
	Theta.Product(Pi, X.T(), Y)
	//fT := mat64.Formatted(Theta)
	//trace.Printf("Theta:\n%v\n", fT)
	return &linearRegression{
		theta: Theta,
	}, nil
}
