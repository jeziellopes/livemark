package sections

import (
	"testing"
)

func TestRank_MergedFirst(t *testing.T) {
	if rank(contribution{merged: true}) != 0 {
		t.Error("merged should have rank 0")
	}
}

func TestRank_OpenSecond(t *testing.T) {
	if rank(contribution{open: true}) != 1 {
		t.Error("open should have rank 1")
	}
}

func TestRank_ClosedLast(t *testing.T) {
	if rank(contribution{}) != 2 {
		t.Error("closed should have rank 2")
	}
}

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
