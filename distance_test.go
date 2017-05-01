package learn

import "testing"

func TestHammingD(t *testing.T) {
	cases := []struct {
		a, b uint
		d    uint
	}{
		{1, 3, 1},
		// 100, 011
		{4, 3, 3},
		// 1000
		{8, 0, 1},
		// 10000 01111
		{16, 15, 5},
	}
	for _, c := range cases {
		d := hammingD(c.a, c.b)
		if d != c.d {
			t.Fatalf("wrong hamming distance: %b <> %b = %d, wants %d", c.a, c.b, d, c.d)
		}
	}
}
