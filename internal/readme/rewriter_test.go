package readme

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sampleReadme = `# Hello

<!-- PROJECTS_START -->
### Old content
<!-- PROJECTS_END -->

<!-- OSS_START -->
### Old OSS
<!-- OSS_END -->
`

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "README.md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestRewrite_ChangesContent(t *testing.T) {
	path := writeTemp(t, sampleReadme)
	changed, err := Rewrite(path, "PROJECTS", "### New content")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Fatal("expected changed=true")
	}
	got, _ := os.ReadFile(path)
	if !strings.Contains(string(got), "### New content") {
		t.Error("new content not written to file")
	}
	if strings.Contains(string(got), "### Old content") {
		t.Error("old content still present")
	}
}

func TestRewrite_NoChangeWhenIdentical(t *testing.T) {
	const body = "### Exact content"
	readme := "# Hi\n\n<!-- PROJECTS_START -->\n" + body + "\n<!-- PROJECTS_END -->\n"
	path := writeTemp(t, readme)
	changed, err := Rewrite(path, "PROJECTS", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changed {
		t.Fatal("expected changed=false when content is identical")
	}
}

func TestRewrite_ErrorWhenZoneMissing(t *testing.T) {
	path := writeTemp(t, "# No zones here\n")
	_, err := Rewrite(path, "PROJECTS", "content")
	if err == nil {
		t.Fatal("expected error for missing zone")
	}
}

func TestRewrite_ErrorWhenFileMissing(t *testing.T) {
	_, err := Rewrite(filepath.Join(t.TempDir(), "nonexistent.md"), "PROJECTS", "x")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRewrite_OnlyTargetedZoneChanges(t *testing.T) {
	path := writeTemp(t, sampleReadme)
	_, err := Rewrite(path, "PROJECTS", "### New projects")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := os.ReadFile(path)
	// OSS zone must be untouched
	if !strings.Contains(string(got), "### Old OSS") {
		t.Error("OSS zone was unexpectedly modified")
	}
}
