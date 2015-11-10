package main

import (
	"flag"
	fmt "github.com/k0kubun/pp"
	"path/filepath"
)

type options struct {
	all bool
}

type tree struct {
	opts        *options
	infos       []*info
	dirs, files int
}

// Call visit
func (t *tree) visit() {
	for _, inf := range t.infos {
		if inf.err == nil {
			d, f := inf.visit(t.opts)
			t.dirs, t.files = t.dirs+d, t.files+f
		}
	}
	fmt.Println(t.dirs, t.files)
}

func (t *tree) print() {
	fmt.Println(t)
}

var (
	all = flag.Bool("a", false, "")
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
			all: *all,
		},
		dirs:  -1,
		infos: make([]*info, len(dirs)),
	}
	for i, dir := range dirs {
		path, err := filepath.Abs(dir)
		if err != nil {
			tr.infos[i] = &info{path: dir, err: err}
		}
		tr.infos[i] = &info{path: path, depth: 0}
	}
	tr.visit() // visit all infos
	//tr.print() // print based on options format
}
