// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"fmt"
)

// unknownTypeErr assembles an error
// for unrecognized type of feature.
func unknownTypeErr(a interface{}) error {
	return fmt.Errorf("learn: type of \"%v\" must be float or string not %T", a, a)
}
func typeMismatchErr(a, b interface{}) error {
	return fmt.Errorf("learn: type mismatch in features \"v\" \"v\"\n", a, b)
}
