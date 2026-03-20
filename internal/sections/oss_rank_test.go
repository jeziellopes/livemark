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

func TestOSSSort_MergedBeforeOpenBeforeClosed(t *testing.T) {
	input := []contribution{
		{title: "closed", status: "❌ Closed"},
		{title: "open", status: "🔄 Open", open: true},
		{title: "merged", status: "✅ Merged", merged: true},
		{title: "closed2", status: "❌ Closed"},
		{title: "open2", status: "🔄 Open", open: true},
	}

	// Apply the same sort used in BuildOSS
	sortContribs(input)

	wantOrder := []string{"merged", "open", "open2", "closed", "closed2"}
	for i, want := range wantOrder {
		if input[i].title != want {
			t.Errorf("position %d: want %q, got %q", i, want, input[i].title)
		}
	}
}
