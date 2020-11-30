package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/gadzira/home_work_melnichenko/hw07_file_copying/copy"
	"github.com/stretchr/testify/require"
)

type test struct {
	in     string
	out    string
	golden string
	offset int64
	limit  int64
}

func TestCopy(t *testing.T) {

	t.Run("If offset more than file size", func(t *testing.T) {
		err := copy.Copy("./testdata/input.txt", "output.txt", 7000, 25)
		expectedErr := errors.New("offset exceeds file size")
		require.Equal(t, expectedErr, err)
	})

	t.Run("If file unknown length", func(t *testing.T) {
		err := copy.Copy("/dev/urandom", "output.txt", 0, 0)
		expectedErr := errors.New("unsupported file")
		require.Equal(t, expectedErr, err)
	})

	for _, tst := range [...]test{
		{
			in:     "./testdata/input.txt",
			out:    "out_offset0_limit0.txt",
			golden: "./testdata/out_offset0_limit0.txt",
			offset: 0,
			limit:  0,
		},
		{
			in:     "./testdata/input.txt",
			out:    "out_offset0_limit10.txt",
			golden: "./testdata/out_offset0_limit10.txt",
			offset: 0,
			limit:  10,
		},
		{
			in:     "./testdata/input.txt",
			out:    "out_offset0_limit1000.txt",
			golden: "./testdata/out_offset0_limit1000.txt",
			offset: 0,
			limit:  1000,
		},
		{
			in:     "./testdata/input.txt",
			out:    "out_offset0_limit10000.txt",
			golden: "./testdata/out_offset0_limit10000.txt",
			offset: 0,
			limit:  10000,
		},
		{
			in:     "./testdata/input.txt",
			out:    "out_offset100_limit1000.txt",
			golden: "./testdata/out_offset100_limit1000.txt",
			offset: 100,
			limit:  1000,
		},
		{
			in:     "./testdata/input.txt",
			out:    "out_offset6000_limit1000.txt",
			golden: "./testdata/out_offset6000_limit1000.txt",
			offset: 6000,
			limit:  1000,
		},
	} {
		t.Run(fmt.Sprintf("subtest-for-%q", tst.out), func(t *testing.T) {
			_ = copy.Copy(tst.in, tst.out, tst.offset, tst.limit)
			out, _ := ioutil.ReadFile(tst.out)
			golden, _ := ioutil.ReadFile(tst.golden)
			defer os.Remove(tst.out)

			if !bytes.Equal(out, golden) {
				t.Errorf("incoming file and outcomming file not matched")
			}
		})
	}
}
