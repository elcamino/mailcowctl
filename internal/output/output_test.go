package output

import "testing"

func TestMaskSecret(t *testing.T) {
	cases := []struct {
		short, full, want string
	}{
		{"abc...", "abcdef", "abc..."},
		{"", "abcdef", "***"},
		{"", "", ""},
	}
	for _, c := range cases {
		if got := maskSecret(c.short, c.full); got != c.want {
			t.Fatalf("maskSecret(%q,%q) = %q, want %q", c.short, c.full, got, c.want)
		}
	}
}
