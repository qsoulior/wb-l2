package main

import (
	"reflect"
	"testing"
)

func TestSearchAnagramms(t *testing.T) {
	words := []string{"листок", "пятак", "пятка", "столик", "тяпка", "слиток", "тест"}
	want := map[string][]string{"листок": {"листок", "слиток", "столик"}, "пятак": {"пятак", "пятка", "тяпка"}}

	if got := SearchAnagramms(words); !reflect.DeepEqual(got, want) {
		t.Errorf("SearchAnagramms() = %v, want %v", got, want)
	}
}
