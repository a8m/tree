package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

type info struct {
	os.FileInfo
	path  string
	depth int
	err   error
	infos []*info
}

type Fs interface {
	Stat(path string) (os.FileInfo, error)
	ReadDir(path string) ([]string, error)
}

// Visit fn
func (inf *info) visit(opts *options) (dirs, files int) {
	fi, err := opts.fs.Stat(inf.path)
	if err != nil {
		inf.err = err
		return
	}
	inf.FileInfo = fi
	if !fi.IsDir() {
		return 0, 1
	}
	names, err := opts.fs.ReadDir(inf.path)
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
	if !inf.IsDir() {
		var props []string
		var stat = inf.Sys().(*syscall.Stat_t)
		// inodes
		if opts.inodes {
			props = append(props, fmt.Sprintf("%d", stat.Ino))
		}
		// device
		if opts.device {
			props = append(props, fmt.Sprintf("%3d", stat.Dev))
		}
		// Mode
		if opts.fileMode {
			props = append(props, inf.Mode().String())
		}
		// Owner/Uid
		if opts.showUid {
			uid := strconv.Itoa(int(stat.Uid))
			if u, err := user.LookupId(uid); err != nil {
				props = append(props, fmt.Sprintf("%-8s", uid))
			} else {
				props = append(props, fmt.Sprintf("%-8s", u.Username))
			}
		}
		// Gorup/Gid
		// TODO: support groupname
		if opts.showGid {
			gid := strconv.Itoa(int(stat.Gid))
			props = append(props, fmt.Sprintf("%-4s", gid))
		}
		// Size
		if opts.byteSize || opts.unitSize {
			var size string
			if opts.unitSize {
				size = fmt.Sprintf("%4s", formatBytes(inf.Size()))
			} else {
				size = fmt.Sprintf("%11d", inf.Size())
			}
			props = append(props, size)
		}
		// Last modification
		if opts.lastMod {
			props = append(props, inf.ModTime().Format("Jan 02 15:04"))
		}
		// Print properties
		if len(props) > 0 {
			fmt.Printf("[%s]  ", strings.Join(props, " "))
		}
	}
	// name/path
	var name string
	if inf.depth == 0 || opts.fullPath {
		name = inf.path
	} else {
		name = inf.Name()
	}
	// Quotes
	if opts.quotes {
		name = fmt.Sprintf("\"%s\"", name)
	}
	// Print file details
	fmt.Println(name)
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

// Convert bytes to human readable string. Like a 2 MB, 64.2 KB, 52 B
func formatBytes(i int64) (result string) {
	var n float64
	sFmt, eFmt := "%.01f", ""
	switch {
	case i > (1024 * 1024 * 1024):
		eFmt = "G"
		n = float64(i) / 1024 / 1024 / 1024
	case i > (1024 * 1024):
		eFmt = "M"
		n = float64(i) / 1024 / 1024
	case i > 1024:
		eFmt = "K"
		n = float64(i) / 1024
	default:
		sFmt = "%.0f"
		n = float64(i)
	}
	if eFmt != "" && n >= 10 {
		sFmt = "%.0f"
	}
	result = fmt.Sprintf(sFmt+eFmt, n)
	result = strings.Trim(result, " ")
	return
}
