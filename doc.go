// Copyright (c) 2015 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

// Package learn exposes some machine learning alghoritms with a focus on practical usage.
//
// It tries to be as idiomatic as possible. Interfaces are used in public APIs when possible
// to make methods adaptable to custom needs.
// Table interface should makes (hopefully) easy to use storage other than memory
// when dealing with "Big Data" (database, filesystem etc..)
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
// choosed at runtime:
//
// - manhattan for numerical features
//
// - humming distance for categorical features
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
package learn
