package ostree

import (
	"bytes"
	"os"

	"github.com/a8m/tree"
)

// FS uses the system filesystem
type FS struct{}

// Stat a path
func (f *FS) Stat(path string) (os.FileInfo, error) {
	return os.Lstat(path)
}

// ReadDir reads a directory
func (f *FS) ReadDir(path string) ([]string, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	names, err := dir.Readdirnames(-1)
	dir.Close()
	if err != nil {
		return nil, err
	}
	return names, nil
}

// Print a tree of the directory
func Print(dir string) string {
	b := new(bytes.Buffer)
	tr := tree.New(dir)
	opts := &tree.Options{
		Fs:      new(FS),
		OutFile: b,
	}
	tr.Visit(opts)
	tr.Print(opts)
	return b.String()
}
