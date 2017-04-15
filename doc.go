// Copyright (c) 2015 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

// Package learn exposes some machine learning algorithms.
//
// Regression:
//
//	- linear regression
//
// Clustering:
//
//	- k means clustering
//
// Classification:
//
//	- kNN
//
// K means clustering
//
// Categorical and numerical features are supported.
//
// Method for distance calculation is automatically
// chosen at runtime:
//
// - Manhattan for numerical features
//
// - Hamming distance for categorical features
//
// Example of data
//
//	Hours	Choices		Stars	Price
//	12,	"A,C",		5,	15.10
//	1,	"D"		1,	1
//
// Categorical features must be translated to an array of 0 and 1:
//
//	Hours	Choices		Stars	Price
//	12,	"[1,0,1,0]",	5,	15.10
//	1,	"[0,0,0,1]"	1,	1
//
// Beware: experimental package, APIs are unstable
// and can change quickly.
package learn
