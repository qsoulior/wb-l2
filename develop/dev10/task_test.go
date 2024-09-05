package main

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
)

func TestRedirectStream(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	want := "abc bca"
	r := io.NopCloser(strings.NewReader(want))
	w := &bytes.Buffer{}

	RedirectStream(ctx, r, w)
	if got := w.String(); got != want {
		t.Errorf("RedirectStream() = %s, want %s", got, want)
	}
}
