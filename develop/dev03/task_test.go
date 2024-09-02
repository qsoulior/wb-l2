package main

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestParseNumeric(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"InvalidString", args{"abcabc"}, 0},
		{"ValidString", args{"abc10abc"}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseNumeric(tt.args.s); got != tt.want {
				t.Errorf("ParseNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompareNumeric(t *testing.T) {
	type args struct {
		a string
		b string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"InvalidStrings", args{"abcabc", "cbacba"}, 0},
		{"ValidStrings", args{"abc11abc", "cba10cba"}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareNumeric(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("CompareNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseNumericWithSuffix(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"InvalidString", args{"abcabc"}, 0},
		{"ValidString", args{"abc10Kabc"}, 10000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseNumericWithSuffix(tt.args.s); got != tt.want {
				t.Errorf("ParseNumericWithSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompareNumericWithSuffix(t *testing.T) {
	type args struct {
		a string
		b string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"InvalidStrings", args{"abcabc", "cbacba"}, 0},
		{"ValidStrings", args{"abc1Mabc", "cba400Kcba"}, 6e5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareNumericWithSuffix(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("CompareNumericWithSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseMonth(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"InvalidString", args{"abcabc"}, 0},
		{"ValidString", args{"abcJUNabc"}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseMonth(tt.args.s); got != tt.want {
				t.Errorf("ParseMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompareMonth(t *testing.T) {
	type args struct {
		a string
		b string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"InvalidStrings", args{"abcabc", "cbacba"}, 0},
		{"ValidStrings", args{"abcJUNabc", "cbaFEBcba"}, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareMonth(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("CompareMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	type args struct {
		opts SortOptions
	}
	tests := []struct {
		name string
		args args
		want CompareFunc
	}{
		{"Compare", args{SortOptions{}}, strings.Compare},
		{"CompareNumeric", args{SortOptions{numeric: true}}, CompareNumeric},
		{"CompareHumanNumeric", args{SortOptions{humanNumeric: true}}, CompareNumericWithSuffix},
		{"CompareMonth", args{SortOptions{month: true}}, CompareMonth},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Compare(tt.args.opts); fmt.Sprintf("%p", got) != fmt.Sprintf("%p", tt.want) {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIgnoreBlanks(t *testing.T) {
	in := []string{"   abc   ", " bcd", "def "}
	IgnoreBlanks(in)
	out := []string{"abc", "bcd", "def"}

	if !reflect.DeepEqual(in, out) {
		t.Errorf("got = %v, want %v", in, out)
	}
}

func TestIsSorted(t *testing.T) {
	type args struct {
		in   []string
		opts SortOptions
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Sorted", args{[]string{"abc", "bca", "bca"}, SortOptions{}}, true},
		{"Unsorted", args{[]string{"bca", "abc", "bca"}, SortOptions{}}, false},
		{"SortedReverse", args{[]string{"bca", "abc", "abc"}, SortOptions{reverse: true}}, true},
		{"SortedUnique", args{[]string{"abc", "bca"}, SortOptions{unique: true}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSorted(tt.args.in, tt.args.opts); got != tt.want {
				t.Errorf("IsSorted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSort(t *testing.T) {
	type args struct {
		in   []string
		opts SortOptions
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"Sort", args{[]string{"bca", "abc", "bca"}, SortOptions{}}, []string{"abc", "bca", "bca"}},
		{"SortReverse", args{[]string{"bca", "abc", "bca"}, SortOptions{reverse: true}}, []string{"bca", "bca", "abc"}},
		{"SortUnique", args{[]string{"bca", "abc", "bca"}, SortOptions{unique: true}}, []string{"abc", "bca"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Sort(&tt.args.in, tt.args.opts)
			if !reflect.DeepEqual(tt.args.in, tt.want) {
				t.Errorf("got = %v, want %v", tt.args.in, tt.want)
			}
		})
	}
}

func TestReadStrings(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"ValidReader", args{strings.NewReader("abc\nbca")}, []string{"abc", "bca"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadStrings(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadStrings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}
