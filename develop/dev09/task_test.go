package main

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestWget_ExtractAttrs(t *testing.T) {
	r := strings.NewReader(`<a href="test1"></a><div></div><a href="test2"></a>`)
	tags := map[string]set[string]{"a": {"href": {}}}
	want := map[string]set[string]{"a": {"test1": {}, "test2": {}}}

	w := Wget{}
	if got := w.ExtractAttrs(r, tags); !reflect.DeepEqual(got, want) {
		t.Errorf("Wget.ExtractAttrs() = %v, want %v", got, want)
	}
}

func TestWget_ExtractLinks(t *testing.T) {
	w := Wget{}
	r := strings.NewReader("<a href=\"test1\"></a><div></div><a href=\"\t\"></a>")
	tags := map[string]set[string]{"a": {"href": {}}}
	want := map[string][]string{"a": {"http://example.com/test1"}}

	type args struct {
		r    io.Reader
		host string
		tags map[string]set[string]
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]string
		wantErr bool
	}{
		{"InvalidHost", args{nil, "\t", nil}, nil, true},
		{"ValidHost", args{r, "http://example.com", tags}, want, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := w.ExtractLinks(tt.args.r, tt.args.host, tt.args.tags)
			if (err != nil) != tt.wantErr {
				t.Errorf("Wget.ExtractLinks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Wget.ExtractLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}
