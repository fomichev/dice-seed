package main

import (
	"math/big"
	"testing"
)

func TestConv(t *testing.T) {
	var tests = []string{
		"0",
		"5",
		"f",
		"ff",
		"ffff",
		"ffffffff",
		"ffffffffffffffff",
		"ffffffffffffffffffffffffffffffff",
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
	}

	for _, test := range tests {
		expected16 := test

		ret, _ := big.NewInt(0).SetString(expected16, 16)
		base6 := ret.Text(6)

		maxLen := ENT / 8 * 2
		if len(expected16) > maxLen {
			expected16 = expected16[len(expected16)-maxLen:]
		}

		n, _ := rollsToSeed(base6, 0)
		got16 := n.Text(16)

		if got16 != expected16 {
			t.Errorf("rollsToSeed(%s, 0):\ngot\t\t0x%s,\nexpected\t0x%s", base6, got16, expected16)
		}
	}
}

func TestError(t *testing.T) {
	var tests = []string{
		"17",
		"1a",
	}

	for _, test := range tests {
		_, err := rollsToSeed(test, 0)
		if err == nil {
			t.Errorf("rollsToSeed(%s, 0): expected error", test)
		}
	}
}
