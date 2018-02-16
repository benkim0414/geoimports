package imports

import (
	"testing"
	"time"
)

func TestImportTypeAvailable(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		typ  *ImportType
		want bool
	}{
		{
			name: "today",
			typ:  &ImportType{AvailableVersionSpottedAt: now.Format(time.RFC3339)},
			want: true,
		},
		{
			name: "yesterday",
			typ:  &ImportType{AvailableVersionSpottedAt: now.Add(-24 * time.Hour).Format(time.RFC3339)},
			want: false,
		},
	}
	for _, tt := range tests {
		if got, want := tt.typ.Available(), tt.want; got != want {
			t.Errorf("Available (%q) = %v, want %v", tt.name, got, want)
		}
	}
}
