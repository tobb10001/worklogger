package worklogger

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

// Start working on something.
func Start(root_dir string, calltime time.Time, topic string) error {
	err := os.MkdirAll(root_dir, 0755)
	if err != nil {
		return err
	}

	isodate := calltime.Format(time.DateOnly)
	isotime := calltime.Format("15:04")
	dayfile := path.Join(root_dir, isodate+".md")

	var content string
	if _, err = os.Stat(dayfile); err != nil {
		content = fmt.Sprintf("# Worklog for %s\n\n", isodate)
	} else {
		content_bytes, err := os.ReadFile(dayfile)
		if err != nil {
			return err
		}
		content = string(content_bytes)
	}

	content = content + fmt.Sprintf("- %s - CURRENT %s\n", isotime, topic)

	err = os.WriteFile(dayfile, []byte(content), 0644)

	return err
}

// End working on something
func End(root_dir string, calltime time.Time) error {
	_, isotime, dayfile := reusables(root_dir, calltime)

	content_bytes, err := os.ReadFile(dayfile)
	if err != nil {
		return err
	}

	content := string(content_bytes)
	lines := strings.Split(content, "\n")
	lines[len(lines)-2] = strings.Replace(lines[len(lines)-2], "CURRENT", isotime, 1)
	content = strings.Join(lines, "\n")

	err = os.WriteFile(dayfile, []byte(content), 0644)

	return err
}

// Returns
// - the formatted date
// - the formatted time
// - the filename of the dayfile
func reusables(root_dir string, calltime time.Time) (string, string, string) {
	isodate := calltime.Format(time.DateOnly)
	isotime := calltime.Format(time.TimeOnly)[:5]
	dayfile := path.Join(root_dir, isodate+".md")
	return isodate, isotime, dayfile
}
