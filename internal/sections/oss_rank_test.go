package sections

import (
	"testing"
)

func TestOSSSort_NewestFirst(t *testing.T) {
	input := []contribution{
		{title: "oldest", createdAt: "2026-03-15T10:00:00Z"},
		{title: "newest", createdAt: "2026-03-21T02:00:00Z"},
		{title: "middle", createdAt: "2026-03-18T15:00:00Z"},
	}

	sortContribs(input)

	wantOrder := []string{"newest", "middle", "oldest"}
	for i, want := range wantOrder {
		if input[i].title != want {
			t.Errorf("position %d: want %q, got %q", i, want, input[i].title)
		}
	}
}
