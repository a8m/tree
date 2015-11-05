package main

import (
	"flag"
	fmt "github.com/k0kubun/pp"
	"os"
	"path/filepath"
)

var (
	input = flag.String("i", "", "")
)

func main() {
	flag.Parse()
	fmt.Println("input ==>", *input)

	var infos []*Info
	var dirs = []string{"."}
	if args := flag.Args(); len(args) > 0 {
		dirs = args
	}
	// Loop over all root dirs
	for _, dir := range dirs {
		var info *Info
		if path, err := filepath.Abs(dir); err != nil {
			info = &Info{path: path, err: err}
		} else {
			fi, err := os.Stat(path)
			info = &Info{fi, path, err, nil}
		}
		infos = append(infos, info)
	}
	// Print infos
	for _, info := range infos {
		if info.err == nil {
			fmt.Println(info.path, info.Name(), info.IsDir())
		} else {
			fmt.Println(info.path, info.err.Error())
		}
	}
}
