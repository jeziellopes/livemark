package readme

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// PreviewZone holds a named zone and its generated markdown content.
type PreviewZone struct {
	Name    string
	Content string
}

// WritePreview writes a self-contained markdown preview file at path,
// showing each zone's generated content under its own section heading.
// The file is intended for local inspection only and should not be committed.
func WritePreview(path, username string, zones []PreviewZone) error {
	var sb strings.Builder

	fmt.Fprintf(&sb, "# livemark preview — %s\n", username)
	fmt.Fprintf(&sb, "_Generated: %s_\n", time.Now().UTC().Format("2006-01-02 15:04 UTC"))

	for _, z := range zones {
		sb.WriteString("\n---\n\n")
		fmt.Fprintf(&sb, "## %s zone\n\n", z.Name)
		sb.WriteString(z.Content)
		sb.WriteString("\n")
	}

	return os.WriteFile(path, []byte(sb.String()), 0644)
}
