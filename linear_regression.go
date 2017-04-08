// Copyright (c) 2015 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"errors"

	"github.com/eraclitux/trace"
	"github.com/gonum/matrix/mat64"
)

type Regression interface {
	Fit(Table) error
	Predict(Table) ([]float64, error)
}

type LinearRegression struct {
	// X is a representation of design matrix,
	// an (m x n+1) matrix with cases as rows
	// wheree X(i,1) is always 1.
	//	1 x1 ... xn+1
	//	...
	//	1 x1 ... xn+1
	// FIXME we really need to store this?
	X *mat64.Dense
	// y is a vector of observed values
	// of dependent variable.
	y               *mat64.Dense
	gradientDescent bool
	theta           *mat64.Dense
}

// NewLinearRegression returns an implementation of
// Regression interface.
//
// Data is a table with training samples as rows.
// Last element in the row MUST be
// the observed values of dependent variable y.
//
// Design matrix X (sample data as row with X(i,1)==1)
// will be created automatically.
//
// NOTE: this function will load whole Table in memory.
func NewLinearRegression(Data Table) (Regression, error) {
	// m: number of samples
	// n: number of features

	// TODO automatic select between gradient descent and normal equation.

	m := Data.Len()
	if m <= 0 {
		return nil, NoData
	}
	s1, err := Data.Row(0)
	if err != nil {
		return nil, NoData
	}
	// o = n+1 because table stores
	// values of y in the last column.
	// +1 for theta0 == 1 terms.
	o := len(s1)
	// TODO check if m < n

	// Build design matrix and y vector.
	X := mat64.NewDense(m, o, nil)
	yRows := make([]float64, m)
	for i := 0; i < m; i++ {
		// Row for design matrix.
		row := make([]float64, o)
		r, err := Data.Row(i)
		if err != nil {
			return nil, NoData
		}
		// x0
		row[0] = 1
		for j, e := range r {
			if f, ok := e.(float64); !ok {
				return nil, unknownType(e)
			} else {
				if j >= o {
					return nil, errors.New("index out of range")
				} else if j == o-1 {
					// y element
					yRows[i] = f
					trace.Println("y:", f)
				} else {
					row[j+1] = f
				}
			}
		}
		trace.Println("row:", row)
		X.SetRow(i, row)
	}
	y := mat64.NewDense(m, 1, yRows)
	trace.Println(X)
	trace.Println(y)
	// theta = pinv(X'*X)*X'*y
	//var M mat64.Dense
	return &LinearRegression{
		// FIXME make design matrix.
		X: X,
		y: y,
	}, nil
}

// Fit can be used to update training data for the specific
// estimeted y.
func (lr *LinearRegression) Fit(t Table) error {
	return nil
}

// Predict given a Table with feature set to predict in its rows:
// 	x1 x2 ... xn
//	...
// 	x1 x2 ... xn
// returns a slice of float of estimeted y values.
//
// Features must be stored as float64 in Table.
func (lr *LinearRegression) Predict(t Table) ([]float64, error) {
	var y mat64.Dense
	ys := make([]float64, t.Len())
	for i := 0; i < t.Len(); i++ {
		row, err := t.Row(i)
		if err != nil {
			return nil, err
		}
		// +1 for x0
		l := len(row) + 1
		a := make([]float64, l)
		for _, e := range row {
			if _, ok := e.(float64); !ok {
				return nil, unknownType(e)
			}
		}
		// for x0
		a[0] = 1
		X := mat64.NewDense(1, l, a)
		// y is a (1 x 1) matrix
		y.Mul(X, lr.theta)
		ys[i] = y.At(1, 1)
	}
	return ys, nil
}
