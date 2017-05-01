// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"reflect"
	"testing"
)

func TestZeroCentroid(t *testing.T) {
	toZero := []interface{}{
		float64(0.9),
		float64(0.3),
		newCategory("foo", nil),
	}
	zero := []interface{}{
		0.0,
		0.0,
		&category{},
	}
	zeroCentroid(toZero)
	if !reflect.DeepEqual(toZero, zero) {
		t.Fatal("not zeroed:", toZero)
	}
}

func TestIncrementCentroid(t *testing.T) {
	incrementing := []interface{}{
		float64(1),
		float64(0.5),
		newCategory("0,1,1,1", nil),
	}
	increment := []interface{}{
		float64(1),
		float64(0.5),
		newCategory("0,0,0,1", nil),
	}
	expected := []interface{}{
		float64(2),
		float64(1),
		newCategory("1,0,0,0", nil),
	}
	if testing.Verbose() {
		t.Logf("incrementing: %v, increment: %v\n", incrementing, increment)
	}
	incrementCentroid(incrementing, increment)
	if !reflect.DeepEqual(incrementing, expected) {
		t.Fatalf("wrongly incremented, expected: %v got: %v", expected, incrementing)
	}
	if testing.Verbose() {
		t.Logf("incremented: %v\n", incrementing)
	}
}

func TestCenterCentroid(t *testing.T) {
	incrementing := []interface{}{
		float64(1),
		float64(0.5),
		newCategory("0,0,0,1", nil),
	}
	increment := []interface{}{
		float64(3),
		float64(1.5),
		newCategory("0,0,1,1", nil),
	}
	expected := []interface{}{
		float64(2),
		float64(1),
		newCategory("0,0,1,0", nil),
	}
	incrementCentroid(incrementing, increment)
	centerCentroid(incrementing, 2)
	if !reflect.DeepEqual(incrementing, expected) {
		t.Fatalf("wrongly incremented, expected: %v got: %v", expected, incrementing)
	}
}

func TestCategory_Distance(t *testing.T) {
	cases := []struct {
		a, b     *category
		distance float64
	}{
		{newCategory("1,1,1,1", nil), newCategory("0,0,0,0", nil), 1},
		{newCategory("0,1,0,1", nil), newCategory("1,0,1,0", nil), 1},
		{newCategory("0,0,1,1", nil), newCategory("0,0,0,0", nil), 0.5},
	}
	for i, c := range cases {
		d := c.a.distance(c.b)
		if d != c.distance {
			t.Fatalf("in case %d, expected: %f, got: %f", i, c.distance, d)
		}
	}
}

func TestCategory_Mean(t *testing.T) {
	cases := []struct {
		a    *category
		cats []*category
		mean *category
	}{
		{
			newCategory("0,0,0,0", nil),
			[]*category{newCategory("1,1,1,1", nil), newCategory("1,1,1,1", nil)},
			newCategory("1,1,1,1", nil),
		},
		{
			newCategory("0,0,0,0", nil),
			[]*category{newCategory("0,0,1,1", nil), newCategory("1,1,0,0", nil), newCategory("1,1,0,0", nil)},
			newCategory("1,0,0,1", nil),
		},
		{
			newCategory("0,0,0,0", nil),
			[]*category{newCategory("0,0,1,1", nil), newCategory("1,1,0,0", nil), newCategory("1,1,0,0", nil), newCategory("1,1,0,0", nil)},
			newCategory("1,0,0,1", nil),
		},
	}
	for i, c := range cases {
		for _, e := range c.cats {
			c.a.add(e)
		}
		c.a.mean(len(c.cats))
		if !reflect.DeepEqual(c.a, c.mean) {
			t.Fatalf("in case %d, expected: %v, got: %v", i+1, c.mean, c.a)
		}
	}
}

func TestCreateRandCategory(t *testing.T) {
	a := createRandCategory(4)
	b := createRandCategory(4)
	if reflect.DeepEqual(a, b) {
		t.Fatal("equal!")
	}
}

func TestCreateRandomCentroids(t *testing.T) {
	f := []interface{}{
		float64(1),
		newCategory("[1,0,0,0]", nil),
	}
	a, err := createRandomCentroids(4, f)
	if err != nil {
		t.Fatal(err)
	}
	b, err := createRandomCentroids(4, f)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(a, b) {
		t.Fatal("equal!")
	}
}

func TestKmc(t *testing.T) {
	data, err := ReadAllCSV("datasets/iris.csv")
	if err != nil {
		t.Fatal(err)
	}
	_, _, _, err = Normalize(data, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	r, err := Kmc(data, 3, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("kmc result:", r)
}
