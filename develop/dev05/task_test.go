package main

import (
	"reflect"
	"testing"
)

func TestNewFixedMatcher(t *testing.T) {
	type args struct {
		pattern string
		opts    MatchOptions
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"WithoutIgnoreCase", args{"Abc", MatchOptions{ignoreCase: false}}, false},
		{"WithIgnoreCase", args{"Abc", MatchOptions{ignoreCase: true}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFixedMatcher(tt.args.pattern, tt.args.opts).Match("abc"); got != tt.want {
				t.Errorf("NewFixedMatcher().Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fixedMatcher_Match(t *testing.T) {
	m := NewFixedMatcher("abc", MatchOptions{ignoreCase: true})
	if got := m.Match("Abc"); !got {
		t.Errorf("fixedMatcher.Match() = %v, want %v", got, true)
	}
}

func TestNewRegexpMatcher(t *testing.T) {
	type args struct {
		pattern string
		opts    MatchOptions
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"WithoutIgnoreCase", args{"Abc", MatchOptions{ignoreCase: false}}, false, false},
		{"WithIgnoreCase", args{"Abc", MatchOptions{ignoreCase: true}}, true, false},
		{"InvalidPattern", args{"\\", MatchOptions{}}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewRegexpMatcher(tt.args.pattern, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRegexpMatcher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if got := m.Match("bca abc bca"); got != tt.want {
					t.Errorf("NewRegexpMatcher().Match() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_regexpMatcher_Match(t *testing.T) {
	m, _ := NewRegexpMatcher("Abc", MatchOptions{ignoreCase: true})
	if got := m.Match("bca abc bca"); !got {
		t.Errorf("regexpMatcher.Match() = %v, want %v", got, true)
	}
}

func TestGrep(t *testing.T) {
	input := []string{"bca", "def", "abc", "efg"}
	matcher := NewFixedMatcher("abc", MatchOptions{})
	type args struct {
		opts GrepOptions
	}
	tests := []struct {
		name string
		args args
		want []bool
	}{
		{"NoOption", args{GrepOptions{}}, []bool{false, false, true, false}},
		{"InvertOption", args{GrepOptions{invert: true}}, []bool{true, true, false, true}},
		{"AfterOption", args{GrepOptions{after: 1}}, []bool{false, false, true, true}},
		{"BeforeOption", args{GrepOptions{before: 1}}, []bool{false, true, true, false}},
		{"ContextOption", args{GrepOptions{context: 1}}, []bool{false, true, true, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Grep(input, matcher, tt.args.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Grep() = %v, want %v", got, tt.want)
			}
		})
	}
}
