package health

import (
	"errors"
	"testing"
)

type checkFn func() error

func ok() error {
	return nil
}
func er() error {
	return errors.New("some error")
}

func TestHealth_Check(t *testing.T) {
	tests := []struct {
		name   string
		checks []checkFn
		want   bool
	}{
		{"#1 all ok", []checkFn{ok, ok, ok, ok, ok, ok}, true},
		{"#2 all ok", []checkFn{ok, er, ok, ok, ok, ok}, false},
		{"#3 all ok", []checkFn{er, er, er, er, er, er}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hc := New()
			for _, c := range tt.checks {
				hc.Watch(c)
			}

			if got := hc.Check(); got != tt.want {
				t.Errorf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}
