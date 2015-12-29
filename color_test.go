package tree

import (
	"testing"
)

var extsTests = []struct {
	name     string
	expected string
}{
	{"foo.jpg", "\x1b[1;35mfoo.jpg\x1b[0m"},
	{"bar.tar", "\x1b[1;31mbar.tar\x1b[0m"},
	{"baz.exe", "\x1b[1;32mbaz.exe\x1b[0m"},
}

func TestExtension(t *testing.T) {
	for _, test := range extsTests {
		f := &file{name: test.name}
		n := &Node{FileInfo: f}
		if actual := ANSIColor(n, f.name); actual != test.expected {
			t.Errorf("\ngot:\n%+v\nexpected:\n%+v", actual, test.expected)
		}
	}
}

var modeTests = []struct {
	name     string
	expected string
	mode     uint
}{
	{"dir", "dir", 122},
}
