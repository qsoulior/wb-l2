package main

import (
	"reflect"
	"testing"
)

func TestCut(t *testing.T) {
	type args struct {
		fields    []int
		delim     string
		separated bool
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{"-f 2,4 -d ':'", args{[]int{1, 3}, ":", false}, " white: frost", true},
		{"-f 5 -d ':'", args{[]int{4}, ":", false}, "", true},
		{"-f 1 -d ','", args{[]int{1}, ",", false}, "Winter: white: snow: frost", true},
		{"-f 1 -d ',' -s", args{[]int{1}, ",", true}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Cut("Winter: white: snow: frost", tt.args.fields, tt.args.delim, tt.args.separated)
			if got != tt.want {
				t.Errorf("Cut() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Cut() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestFields_String(t *testing.T) {
	f := Fields{0, 1}
	want := "[0 1]"

	if got := f.String(); !reflect.DeepEqual(got, want) {
		t.Errorf("Fields.String() = %v, want %v", got, want)
	}
}

func TestFields_Set(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		f       Fields
		args    args
		wantErr bool
	}{
		{"ValidString", Fields{}, args{"1,2"}, false},
		{"InvalidString", Fields{}, args{"0,1"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.Set(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("Fields.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
