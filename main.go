package main

import (
	"flag"
	"fmt"
)

type options struct {
	all      bool
	dirsOnly bool
}

type tree struct {
	opts        *options
	dirs, files int
}

var (
	a = flag.Bool("a", false, "")
	d = flag.Bool("d", false, "")
	f = flag.Bool("f", false, "")
	s = flag.Bool("s", false, "")
	h = flag.Bool("h", false, "")
)

func main() {
	flag.Parse()
	var dirs = []string{"."}
	// Make it work with leading dirs
	if args := flag.Args(); len(args) > 0 {
		dirs = args
	}
	tr := &tree{
		opts: &options{
			all:      *a,
			dirsOnly: *d,
		},
	}
	for _, dir := range dirs {
		inf := &info{path: dir}
		if d, f := inf.visit(tr.opts); inf.err == nil {
			tr.dirs, tr.files = tr.dirs+d-1, tr.files+f
		}
		inf.print("")
	}
	fmt.Printf("\n%d directories, %d files\n", tr.dirs, tr.files)
}
