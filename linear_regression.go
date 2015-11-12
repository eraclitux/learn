package nml

type Regression interface {
	Estimate(Table) (Table, error)
}

type linearRegression struct {
	// X is a representation of design matrix,
	// an (m x n+1) matrix with cases as rows
	// wheree X(i,1) is always 1.
	//	1 x1 ... xn+1
	//	...
	//	1 x1 ... xn+1
	X               Table
	y               Table
	gradientDescent bool
}

// NewLinearRegression returns an implementation of
// Regression interface.
//
// D is a table with training samples as rows.
// y is a table of single elements rows which are
// the observed values of dependent variable.
func NewLinearRegression(D, y Table) Regression {
	// m: number of samples
	// n: number of features
	// D represent an (m x n) matrix of training data.
	// y is a column vector (m x 1) with observed results.
	return &linearRegression{
		// FIXME make design matrix.
		X: D,
		y: y,
	}
}

// Estimate return a Table with single element rows representig
// esimeted y.
func (lr *linearRegression) Estimate(t Table) (Table, error) {
	// TODO automatic select between gradient descent and normal equation.
	//theta := make([]float64, 0)
	return nil, nil
}
