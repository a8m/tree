package node

import (
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

func (fs *MockFs) addFile(path string, file *file) *MockFs {
	fs.files[path] = file
	if file.IsDir() {
		for _, f := range file.files {
			fs.addFile(path+"/"+f.name, f)
		}
	}
	return fs
}

func (fs *MockFs) Stat(path string) (os.FileInfo, error) {
	return fs.files[path], nil
}
func (fs *MockFs) ReadDir(path string) ([]string, error) {
	var names []string
	for _, file := range fs.files[path].files {
		names = append(names, file.Name())
	}
	return names, nil
}

// Mock output file
type Out struct {
	str string
}

func (o *Out) equal(s string) bool {
	return o.str == s
}

func (o *Out) Write(p []byte) (int, error) {
	o.str += string(p)
	return len(p), nil
}

func (o *Out) clear() {
	o.str = ""
}

// Mock files and file-system
var (
	root = &file{
		"root",
		200,
		[]*file{
			&file{"a", 50, nil, time.Now()},
			&file{"b", 50, nil, time.Now()},
			&file{
				"c",
				100,
				[]*file{
					&file{"d", 50, nil, time.Now()},
					&file{"e", 50, nil, time.Now()},
					&file{".f", 0, nil, time.Now()},
				},
				time.Now()},
		},
		time.Now(),
	}
	fs  = NewFs().addFile(root.name, root)
	out = new(Out)
)

type treeTest struct {
	name     string
	opts     *Options
	expected string
}

var tests = []treeTest{
	{"basic", &Options{Fs: fs, OutFile: out}, `root
├── a
├── b
└── c
    ├── d
    └── e
`},
	{"all", &Options{Fs: fs, OutFile: out, All: true}, `root
├── a
├── b
└── c
    ├── d
    ├── e
    └── .f
`},
	{"dirs", &Options{Fs: fs, OutFile: out, DirsOnly: true}, `root
└── c
`}, {"fullPath", &Options{Fs: fs, OutFile: out, FullPath: true}, `root
├── root/a
├── root/b
└── root/c
    ├── root/c/d
    └── root/c/e
`}, {"deepLevel", &Options{Fs: fs, OutFile: out, DeepLevel: 1}, `root
├── a
├── b
└── c
`}}

// Tests
func TestSimple(t *testing.T) {
	for _, test := range tests {
		inf := New(root.name)
		inf.Visit(test.opts)
		inf.Print("", test.opts)
		if !out.equal(test.expected) {
			t.Errorf("%s:\ngot:\n%+v\nexpected:\n%+v", test.name, out.str, test.expected)
		}
		out.clear()
	}
}
