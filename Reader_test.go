package main

import (
	"reflect"
	"testing"
)

func Test_kMain(t *testing.T) {
	type args struct {
		in []CompareResult
	}
	tests := []struct {
		name    string
		args    args
		wantOut []CompareResult
	}{
		{name: "1", args: args{in: []CompareResult{
			{Index: int64(1), A: 0, B: 0},
			{Index: int64(2), A: 0, B: 0},
			{Index: int64(3), A: 0, B: 0},
			{Index: int64(5), A: 0, B: 0},
			{Index: int64(6), A: 0, B: 0},
			{Index: int64(7), A: 0, B: 0},
			{Index: int64(9), A: 0, B: 0},
			{Index: int64(11), A: 0, B: 0},
		}}, wantOut: []CompareResult{
			{Index: int64(1), A: 0, B: 0},
			{Index: int64(2), A: 0, B: 0},
			{Index: int64(3), A: 0, B: 0},
			{Index: int64(5), A: 0, B: 0},
			{Index: int64(6), A: 0, B: 0},
			{Index: int64(7), A: 0, B: 0},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := kMain(tt.args.in); !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("kMain() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
