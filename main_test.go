package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrtol(t *testing.T) {
	cases := map[string]struct {
		input     string
		current   int
		expectNum int
		expectI   int
	}{
		"first num": {
			input:     "5+20-1",
			current:   0,
			expectNum: 5,
			expectI:   1,
		},
		"second num": {
			input:     "5+20-1",
			current:   2,
			expectNum: 20,
			expectI:   2,
		},
		"third num": {
			input:     "5+20-1",
			current:   5,
			expectNum: 1,
			expectI:   1,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			num, i := strtol(tc.input, tc.current)
			assert.Equal(t, tc.expectNum, num)
			assert.Equal(t, tc.expectI, i)
		})
	}
}
