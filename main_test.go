package main

import (
	"reflect"
	"testing"
)

func TestReadCsv(t *testing.T) {
	var tests = []struct {
		in  string
		out [][]float64
	}{
		{
			in: "./test/wave.csv",
			out: [][]float64{
				{126, 98, 108, 88, 113},   // first 5 value
				{117, 111, -115, 124, 10}, // last 5 value
			},
		},
		{
			in: "./test/interm.csv",
			out: [][]float64{
				{4, 3, 5, 1, 5},
				{3, 5, 4, 4, 6},
			},
		},
	}

	for _, tt := range tests {
		got, _ := ReadCsv(tt.in)
		want := [][]float64{
			got[0][:5],
			got[len(got)-1][len(got[0])-5:],
		}

		if !reflect.DeepEqual(want, tt.out) {
			for _, v := range want {
				t.Errorf("got: %f", v)
			}
			for _, v := range tt.out {
				t.Errorf("want: %f", v)
			}
		}
	}
}
