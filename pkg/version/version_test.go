package version

import "testing"

func TestShortSHA(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"full sha is truncated", "8bac50fdfda3a1b2c3d4e5f60718293a4b5c6d7e", "8bac50f"},
		{"exactly seven is unchanged", "8bac50f", "8bac50f"},
		{"shorter than seven is unchanged", "8bac", "8bac"},
		{"empty is unchanged", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shortSHA(tt.in); got != tt.want {
				t.Errorf("shortSHA(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// The build info fallback runs in init, so the best we can assert here is that
// it never leaves a field empty — `a2 version` must always print something.
func TestVersionFieldsArePopulated(t *testing.T) {
	fields := map[string]string{
		"Version":   Version,
		"GitSHA":    GitSHA,
		"BuildDate": BuildDate,
	}
	for name, value := range fields {
		if value == "" {
			t.Errorf("%s is empty", name)
		}
	}
}
