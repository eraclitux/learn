// Copyright (c) 2015 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"errors"
	"fmt"
)

// unknownType assembles an appropiate error
// for unrecognized types.
func unknownType(args ...interface{}) error {
	return errors.New("unrecognized type: " + fmt.Sprint(args...))
}