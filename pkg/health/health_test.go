package health

import "testing"

func TestHealth_Check(t *testing.T) {
	type fields struct {
		checks []func() error
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Health{
				checks: tt.fields.checks,
			}
			if got := h.Check(); got != tt.want {
				t.Errorf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}
