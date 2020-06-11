package hw02_unpack_string //nolint:golint,stylecheck

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type test struct {
	input    string
	expected string
	err      error
}

func TestUnpack(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    "a4bc2d5e",
			expected: "aaaabccddddde",
		},
		{
			input:    "abccd",
			expected: "abccd",
		},
		{
			input:    "3abc",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "45",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "aaa10b",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "",
			expected: "",
		},
		{
			input:    "d\n5abc",
			expected: "d\n\n\n\n\nabc",
		},
		{
			input:    "d;3abc",
			expected: "d;;;abc",
		},
	} {
		t.Run(fmt.Sprintf("subtest-for-%q", tst.input), func(t *testing.T) {
			result, err := Unpack(tst.input)
			require.Equal(t, tst.err, err)
			require.Equal(t, tst.expected, result)
		})
	}
}

func TestUnpackError(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    "dг2bc",
			expected: "dггdc",
		},
	} {
		t.Run(fmt.Sprintf("subtest-for-%q", tst.input), func(t *testing.T) {
			result, _ := Unpack(tst.input)
			require.NotEqual(t, tst.expected, result)
		})
	}
}

func TestUnpackWithEscape(t *testing.T) {
	t.Skip() // Remove if task with asterisk completed

	for _, tst := range [...]test{
		{
			input:    `qwe\4\5`,
			expected: `qwe45`,
		},
		{
			input:    `qwe\45`,
			expected: `qwe44444`,
		},
		{
			input:    `qwe\\5`,
			expected: `qwe\\\\\`,
		},
		{
			input:    `qwe\\\3`,
			expected: `qwe\3`,
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}
