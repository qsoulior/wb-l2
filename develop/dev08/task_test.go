package main

import (
	"reflect"
	"testing"
)

func TestCD_Execute(t *testing.T) {
	c := new(CD)
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"ValidArgument", args{nil}, false},
		{"InvalidArgument", args{[]string{""}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := c.Execute(tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("CD.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPWD_Execute(t *testing.T) {
	c := new(PWD)
	_, err := c.Execute()
	if (err != nil) != false {
		t.Errorf("PWD.Execute() error = %v, wantErr %v", err, false)
		return
	}
}

func TestEcho_Execute(t *testing.T) {
	c := new(Echo)
	got, _ := c.Execute("abc", "bcd")
	want := "abc bcd"
	if got != want {
		t.Errorf("Echo.Execute() = %v, want %v", got, want)
	}
}

func TestKill_Execute(t *testing.T) {
	c := new(Kill)
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"EmptyArgument", args{}, "", false},
		{"InvalidArgument", args{[]string{""}}, "", true},
		{"InvalidArgumentPID", args{[]string{"-1"}}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Execute(tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kill.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Kill.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPS_Execute(t *testing.T) {
	c := new(PS)
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"EmptyArgument", args{}, false},
		{"ValidArgument", args{[]string{"10"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := c.Execute(tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("PS.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestExec_Execute(t *testing.T) {
	c := new(Exec)
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"EmptyArgument", args{}, false},
		{"InvalidArgument", args{[]string{""}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := c.Execute(tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exec.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestParseLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    Command
		want1   []string
		wantErr bool
	}{
		{"ValidLine", args{"echo 123"}, new(Echo), []string{"123"}, false},
		{"InvalidLine", args{""}, nil, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParseLine(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseLine() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ParseLine() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestExecuteLines(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"NoLines", args{nil}, "", false},
		{"ValidFirstLine", args{[]string{"echo 123"}}, "123", false},
		{"InvalidFirstLine", args{[]string{""}}, "", true},
		{"InvalidFirstLineArgument", args{[]string{"cd \\0"}}, "", true},
		{"ValidLines", args{[]string{"echo 123", "echo"}}, "123", false},
		{"InvalidLines", args{[]string{"echo 123", ""}}, "", true},
		{"InvalidLinesArguments", args{[]string{"echo \\0", "cd"}}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExecuteLines(tt.args.lines)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExecuteLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
