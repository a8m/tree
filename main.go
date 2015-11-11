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
	infos       []*info
	dirs, files int
}

// Call visit
func (t *tree) visit() {
	for _, inf := range t.infos {
		d, f := inf.visit(t.opts)
		t.dirs, t.files = t.dirs+d-1, t.files+f
	}
}

func (t *tree) print() {
	for _, inf := range t.infos {
		inf.print("")
	}
	fmt.Printf("\n%d directories, %d files\n", t.dirs, t.files)
}

var (
	a = flag.Bool("a", false, "")
	d = flag.Bool("d", false, "")
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
		infos: make([]*info, len(dirs)),
	}
	for i, dir := range dirs {
		tr.infos[i] = &info{path: dir}
	}
	tr.visit() // visit all infos
	tr.print() // print based on options format
}
