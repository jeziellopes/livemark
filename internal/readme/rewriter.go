package readme

import (
	"fmt"
	"os"
	"regexp"
)

// Rewrite reads the file at path, replaces the content between
// <!-- ZONE_START --> and <!-- ZONE_END --> markers with body,
// and writes the result back only if it changed.
// Returns (changed, error).
func Rewrite(path, zone, body string) (bool, error) {
	original, err := os.ReadFile(path)
	if err != nil {
		return false, fmt.Errorf("reading %s: %w", path, err)
	}

	pattern := regexp.MustCompile(`(?s)(<!-- ` + zone + `_START -->)\n.*?\n(<!-- ` + zone + `_END -->)`)
	if !pattern.Match(original) {
		return false, fmt.Errorf("zone %s not found in %s", zone, path)
	}

	updated := pattern.ReplaceAllString(string(original),
		"<!-- "+zone+"_START -->\n"+body+"\n<!-- "+zone+"_END -->",
	)

	if updated == string(original) {
		return false, nil
	}

	if err := os.WriteFile(path, []byte(updated), 0644); err != nil {
		return false, fmt.Errorf("writing %s: %w", path, err)
	}
	return true, nil
}
