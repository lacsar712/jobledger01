package precopy

import "testing"

// TestClonePrefixRegression guards the independence contract of ClonePrefix:
// the returned slice must not share a backing array with src, and boundary
// handling for n must remain correct after the fix.
func TestClonePrefixRegression(t *testing.T) {
	// Mutating any element of the returned prefix must leave src untouched.
	src := []byte("abcdef")
	got := ClonePrefix(src, 4)
	if want := "abcd"; string(got) != want {
		t.Fatalf("prefix=%q want %q", got, want)
	}
	for i := range got {
		got[i] = 'X'
	}
	if src[0] != 'a' || src[3] != 'd' {
		t.Fatalf("ClonePrefix aliased src after mutation: %q", src)
	}
	if cap(got) == cap(src) && len(got) > 0 && &got[0] == &src[0] {
		t.Fatalf("ClonePrefix shares backing array with src")
	}
}

func TestClonePrefixBoundaries(t *testing.T) {
	src := []byte("abcdef")
	cases := []struct {
		n    int
		want string
	}{
		{0, ""},
		{-3, ""},
		{6, "abcdef"},
		{100, "abcdef"},
	}
	for _, c := range cases {
		got := ClonePrefix(src, c.n)
		if string(got) != c.want {
			t.Fatalf("n=%d: prefix=%q want %q", c.n, got, c.want)
		}
		// Even an empty/clamped result must be safe to mutate without
		// touching src.
		for i := range got {
			got[i] = 'Z'
		}
		if string(src) != "abcdef" {
			t.Fatalf("ClonePrefix(n=%d) aliased src: %q", c.n, src)
		}
	}
}
