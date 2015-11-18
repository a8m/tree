package main

import (
	"flag"
	"fmt"
	"os"
)

type options struct {
	fs       Fs
	all      bool
	dirsOnly bool
	fullPath bool
	byteSize bool
	unitSize bool
}

var (
	a = flag.Bool("a", false, "")
	d = flag.Bool("d", false, "")
	f = flag.Bool("f", false, "")
	s = flag.Bool("s", false, "")
	h = flag.Bool("h", false, "")
)

type fs struct{}

func (f *fs) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func (f *fs) ReadDir(path string) ([]string, error) {
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

func main() {
	var nd, nf int
	var dirs = []string{"."}
	flag.Parse()
	// Make it work with leading dirs
	if args := flag.Args(); len(args) > 0 {
		dirs = args
	}
	opts := &options{
		fs:       new(fs),
		all:      *a,
		dirsOnly: *d,
		fullPath: *f,
		byteSize: *s,
		unitSize: *h,
	}
	for _, dir := range dirs {
		inf := &info{path: dir}
		if d, f := inf.visit(opts); inf.err == nil {
			nd, nf = nd+d-1, nf+f
		}
		inf.print("", opts)
	}
	// print footer
	footer := fmt.Sprintf("\n%d directories", nd)
	if !opts.dirsOnly {
		footer += fmt.Sprintf(", %d files", nf)
	}
	fmt.Println(footer)
}
