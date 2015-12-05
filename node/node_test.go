package node

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"
)

// Mock file/FileInfo
type file struct {
	name    string
	size    int64
	files   []*file
	lastMod time.Time
}

func (f *file) Name() string          { return f.name }
func (f *file) Size() int64           { return f.size }
func (f *file) Mode() (o os.FileMode) { return }
func (f *file) ModTime() time.Time    { return f.lastMod }
func (f *file) IsDir() bool           { return nil != f.files }
func (f *file) Sys() interface{} {
	var s *syscall.Stat_t
	return s
}

// Mock filesystem
type MockFs struct {
	files map[string]*file
}

func NewFs() *MockFs {
	return &MockFs{make(map[string]*file)}
}

func (f *MockFs) addFile(path string, file *file) *MockFs {
	f.files[path] = file
	return f
}

func (f *MockFs) Stat(path string) (os.FileInfo, error) {
	return f.files[path], nil
}
func (f *MockFs) ReadDir(path string) ([]string, error) {
	var names []string
	for _, file := range f.files[path].files {
		names = append(names, file.Name())
	}
	return names, nil
}

// Mock output file
type Out struct {
	str string
}

func (o *Out) Write(p []byte) (int, error) {
	o.str += string(p)
	return len(p), nil
}

// Tests
func TestSimple(t *testing.T) {
	root := &file{
		"root",
		0,
		[]*file{
			&file{"a", 0, nil, time.Now()},
			&file{"b", 0, nil, time.Now()},
		},
		time.Now(),
	}
	fs := NewFs().addFile("root", root)
	for _, f := range root.files {
		fs.addFile("root/"+f.name, f)
	}
	out := new(Out)
	opts := &Options{
		Fs:      fs,
		OutFile: out,
	}
	inf := New("root")
	inf.Visit(opts)
	inf.Print("", opts)
	fmt.Println(out.str)
}
