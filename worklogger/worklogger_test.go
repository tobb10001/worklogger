package worklogger_test

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"worklogger/worklogger"
)

func TestStartCreatesDirectoryIfNotExists(t *testing.T) {
	tempdir, err := os.MkdirTemp("", "testing")
	require.NoError(t, err)
	defer os.RemoveAll(tempdir)

	day := time.Date(2025, 1, 26, 0, 0, 0, 0, time.UTC)

	root_dir := path.Join(tempdir, "worklog")

	err = worklogger.Start(root_dir, day, "topic")
	assert.NoError(t, err)

	_, err = os.Stat(root_dir)

	assert.NoError(t, err)
}

func TestStartWorksIfDirAlreadyExists(t *testing.T) {

	tempdir, err := os.MkdirTemp("", "testing")
	require.NoError(t, err)
	defer os.RemoveAll(tempdir)

	day := time.Date(2025, 1, 26, 0, 0, 0, 0, time.UTC)

	root_dir := path.Join(tempdir, "worklog")
	err = os.Mkdir(root_dir, 0755)
	require.NoError(t, err)

	err = worklogger.Start(root_dir, day, "topic")
	assert.NoError(t, err)
}

func TestStartCreatesNestedDir(t *testing.T) {
	tempdir, err := os.MkdirTemp("", "testing")
	require.NoError(t, err)
	defer os.RemoveAll(tempdir)

	day := time.Date(2025, 1, 26, 0, 0, 0, 0, time.UTC)

	root_dir := path.Join(tempdir, "outer", "worklog")

	err = worklogger.Start(root_dir, day, "topic")
	assert.NoError(t, err)

	_, err = os.Stat(root_dir)
	assert.NoError(t, err)
}

func TestStartCreatesFileWithHeader(t *testing.T) {
	tempdir, err := os.MkdirTemp("", "testing")
	require.NoError(t, err)
	defer os.RemoveAll(tempdir)

	day := time.Date(2025, 1, 26, 1, 2, 3, 4, time.UTC)

	err = worklogger.Start(tempdir, day, "topic")
	assert.NoError(t, err)

	file_path := path.Join(tempdir, "2025-01-26.md")

	dat, err := os.ReadFile(file_path)
	require.NoError(t, err)

	assert.Equal(t,
		`# Worklog for 2025-01-26

- 01:02 - CURRENT topic
`, // trailing newline is important
		string(dat),
	)
}

func TestEnd(t *testing.T) {
	// Arrange
	tempdir, err := os.MkdirTemp("", "testing")
	require.NoError(t, err)
	defer os.RemoveAll(tempdir)

	day := time.Date(2025, 1, 26, 2, 3, 4, 5, time.UTC)

	file_path := path.Join(tempdir, "2025-01-26.md")

	err = os.WriteFile(file_path,
		[]byte(`# Worklog for 2025-01-26

- 01:02 - CURRENT topic
`),
		0644,
	)
	require.NoError(t, err)

	// Act
	err = worklogger.End(tempdir, day)
	require.NoError(t, err)

	// Assert
	dat, err := os.ReadFile(file_path)
	require.NoError(t, err)

	assert.Equal(t,
		`# Worklog for 2025-01-26

- 01:02 - 02:03 topic
`,
		string(dat),
	)
}
