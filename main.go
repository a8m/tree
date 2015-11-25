package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/a8m/tree/node"
)

var (
	// Files
	a      = flag.Bool("a", false, "")
	d      = flag.Bool("d", false, "")
	f      = flag.Bool("f", false, "")
	s      = flag.Bool("s", false, "")
	h      = flag.Bool("h", false, "")
	p      = flag.Bool("p", false, "")
	u      = flag.Bool("u", false, "")
	g      = flag.Bool("g", false, "")
	Q      = flag.Bool("Q", false, "")
	D      = flag.Bool("D", false, "")
	inodes = flag.Bool("inodes", false, "")
	device = flag.Bool("device", false, "")
	// Sort
	U = flag.Bool("U", false, "")
	t = flag.Bool("t", false, "")
	c = flag.Bool("c", false, "")
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
	opts := &node.Options{
		Fs: new(fs),
		// Files
		All:      *a,
		DirsOnly: *d,
		FullPath: *f,
		ByteSize: *s,
		UnitSize: *h,
		FileMode: *p,
		ShowUid:  *u,
		ShowGid:  *g,
		LastMod:  *D,
		Quotes:   *Q,
		Inodes:   *inodes,
		Device:   *device,
		// Sort
		NoSort:    *U,
		ModSort:   *t,
		CTimeSort: *c,
	}
	for _, dir := range dirs {
		inf := node.New(dir)
		if d, f := inf.Visit(opts); f != 0 {
			nd, nf = nd+d-1, nf+f
		}
		inf.Print("", opts)
	}
	// print footer
	footer := fmt.Sprintf("\n%d directories", nd)
	if !opts.DirsOnly {
		footer += fmt.Sprintf(", %d files", nf)
	}
	fmt.Println(footer)
}
