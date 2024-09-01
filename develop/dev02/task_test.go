package main

import "testing"

func Test_Unpack(t *testing.T) {
	type args struct {
		inp string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"ValidString#01", args{"a4bc2d5e"}, "aaaabccddddde", false},
		{"ValidString#02", args{"abcd"}, "abcd", false},
		{"InvalidString#01", args{"45"}, "", true},
		{"ValidString#03", args{""}, "", false},
		{"ValidString#04", args{`qwe\4\5`}, "qwe45", false},
		{"ValidString#05", args{`qwe\45`}, "qwe44444", false},
		{"ValidString#06", args{`qwe\\5`}, `qwe\\\\\`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unpack(tt.args.inp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unpack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Unpack() = %v, want %v", got, tt.want)
			}
		})
	}
}
