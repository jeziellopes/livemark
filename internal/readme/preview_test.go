package readme

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWritePreview_ContainsUsername(t *testing.T) {
	path := filepath.Join(t.TempDir(), "preview.md")
	zones := []PreviewZone{{Name: "PROJECTS", Content: "### Projects"}}
	if err := WritePreview(path, "testuser", zones); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := readFile(t, path)
	if !strings.Contains(got, "testuser") {
		t.Error("preview does not contain username")
	}
}

func TestWritePreview_ContainsZoneHeadings(t *testing.T) {
	path := filepath.Join(t.TempDir(), "preview.md")
	zones := []PreviewZone{
		{Name: "PROJECTS", Content: "### Projects"},
		{Name: "OSS", Content: "### OSS"},
	}
	if err := WritePreview(path, "u", zones); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := readFile(t, path)
	if !strings.Contains(got, "## PROJECTS zone") {
		t.Error("missing PROJECTS heading")
	}
	if !strings.Contains(got, "## OSS zone") {
		t.Error("missing OSS heading")
	}
}

func TestWritePreview_ContainsZoneContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "preview.md")
	zones := []PreviewZone{{Name: "PROJECTS", Content: "### My projects content"}}
	if err := WritePreview(path, "u", zones); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := readFile(t, path)
	if !strings.Contains(got, "### My projects content") {
		t.Error("zone content missing from preview")
	}
}

func TestWritePreview_ContainsTimestamp(t *testing.T) {
	path := filepath.Join(t.TempDir(), "preview.md")
	if err := WritePreview(path, "u", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := readFile(t, path)
	if !strings.Contains(got, "Generated:") {
		t.Error("missing timestamp in preview")
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
