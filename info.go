package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type info struct {
	os.FileInfo
	path  string
	depth int
	err   error
	infos []*info
}

//Think/Write pusodue code
//what is the right way to write this function...
func (inf *info) visit(opts *options) (dirs, files int) {
	fi, err := os.Stat(inf.path)
	if err != nil {
		inf.err = err
		return
	}
	inf.FileInfo = fi
	if !fi.IsDir() {
		return 0, 1
	}
	dir, err := os.Open(inf.path)
	if err != nil {
		inf.err = err
		return
	}
	names, err := dir.Readdirnames(-1)
	dir.Close()
	if err != nil {
		inf.err = err
		return
	}
	inf.infos = make([]*info, 0)
	for _, name := range names {
		// "all" option
		if !opts.all && strings.HasPrefix(name, ".") {
			continue
		}
		_inf := &info{
			path:  filepath.Join(inf.path, name),
			depth: inf.depth + 1,
		}
		d, f := _inf.visit(opts)
		// "dirs only" option
		if opts.dirsOnly && !_inf.IsDir() {
			continue
		}
		inf.infos = append(inf.infos, _inf)
		dirs, files = dirs+d, files+f
	}
	return dirs + 1, files
}

func (inf *info) print(indent string, opts *options) {
	if inf.err != nil {
		err := strings.Split(inf.err.Error(), ": ")[1]
		fmt.Printf("%s [%s]\n", inf.path, err)
		return
	}
	// Print name/path
	if inf.depth == 0 || opts.fullPath {
		fmt.Println(inf.path)
	} else {
		fmt.Println(inf.Name())
	}
	add := "│   "
	for i, _inf := range inf.infos {
		if i == len(inf.infos)-1 {
			fmt.Printf(indent + "└── ")
			add = "    "
		} else {
			fmt.Printf(indent + "├── ")
		}
		_inf.print(indent+add, opts)
	}
}
