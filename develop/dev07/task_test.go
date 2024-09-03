package main

import (
	"testing"
	"time"
)

func Test_or(t *testing.T) {
	sig := func(after time.Duration) <-chan any {
		c := make(chan any)
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	t.Run("NoChannels", func(t *testing.T) {
		if got := or(); got != nil {
			t.Errorf("or() = %v, want %v", got, nil)
		}
	})

	tests := []struct {
		name     string
		channels []<-chan any
	}{
		{"FirstChannel", []<-chan any{sig(0), sig(1 * time.Minute)}},
		{"SecondChannel", []<-chan any{sig(1 * time.Minute), sig(0)}},
		{"ThirdChannel", []<-chan any{sig(1 * time.Minute), sig(1 * time.Minute), sig(0)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, got := <-or(tt.channels...); got {
				t.Errorf("<-or() = %v, want %v", got, false)
			}
		})
	}
}
