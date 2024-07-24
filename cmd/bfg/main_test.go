package main

import (
	"io"
	"testing"
)

func TestInputReader(t *testing.T) {
	table := []struct {
		name  string
		path  string
		dflag bool
		err   bool
	}{
		{
			"stdin test",
			"",
			true,
			false,
		},
		{
			"no stdin test",
			"",
			false,
			true,
		},
		{
			"file test",
			"./main.go",
			true,
			false,
		},
		{
			"bad file test",
			"./file_not_found",
			true,
			true,
		},
		{
			"explicit stdin test",
			"-",
			false,
			false,
		},
	}

	for _, v := range table {
		t.Run(v.name, func(t *testing.T) {
			buff, err := inputReader(v.path, v.dflag)

			if err != nil && !v.err {
				t.Errorf("Error thrown  %v", err)
			} else if err == nil && v.err {
				t.Error("Error was expected")
			} else if !v.err {
				_, ok := buff.(io.ByteReader)
				if !ok {
					t.Error("no reader returned")
				}
			}
		})
	}
}
