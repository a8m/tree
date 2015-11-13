package main

import (
	"flag"
	"fmt"
)

type options struct {
	all      bool
	dirsOnly bool
}

var (
	a = flag.Bool("a", false, "")
	d = flag.Bool("d", false, "")
	f = flag.Bool("f", false, "")
	s = flag.Bool("s", false, "")
	h = flag.Bool("h", false, "")
)

func main() {
	var nd, nf int
	var dirs = []string{"."}
	flag.Parse()
	// Make it work with leading dirs
	if args := flag.Args(); len(args) > 0 {
		dirs = args
	}
	opts := &options{
		all:      *a,
		dirsOnly: *d,
	}
	for _, dir := range dirs {
		inf := &info{path: dir}
		if d, f := inf.visit(opts); inf.err == nil {
			nd, nf = nd+d-1, nf+f
		}
		inf.print("")
	}
	fmt.Printf("\n%d directories, %d files\n", nd, nf)
}
