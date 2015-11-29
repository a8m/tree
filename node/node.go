package node

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

type Node struct {
	os.FileInfo
	path  string
	depth int
	err   error
	nodes Nodes
}

type Nodes []*Node

func (n Nodes) Len() int      { return len(n) }
func (n Nodes) Swap(i, j int) { n[i], n[j] = n[j], n[i] }

type Fs interface {
	Stat(path string) (os.FileInfo, error)
	ReadDir(path string) ([]string, error)
}

type Options struct {
	Fs Fs
	// List
	All       bool
	DirsOnly  bool
	FullPath  bool
	DeepLevel int
	// File
	ByteSize bool
	UnitSize bool
	FileMode bool
	ShowUid  bool
	ShowGid  bool
	LastMod  bool
	Quotes   bool
	Inodes   bool
	Device   bool
	// Sort
	NoSort    bool
	VerSort   bool
	ModSort   bool
	DirSort   bool
	NameSort  bool
	SizeSort  bool
	CTimeSort bool
	ReverSort bool
	// Graphics
	NoIndent bool
}

// New get path and create new node
func New(path string) *Node {
	return &Node{path: path}
}

// Visit fn
func (node *Node) Visit(opts *Options) (dirs, files int) {
	fi, err := opts.Fs.Stat(node.path)
	if err != nil {
		node.err = err
		return
	}
	node.FileInfo = fi
	if !fi.IsDir() {
		return 0, 1
	}
	names, err := opts.Fs.ReadDir(node.path)
	if err != nil {
		node.err = err
		return
	}
	node.nodes = make(Nodes, 0)
	// DeepLevel option
	if opts.DeepLevel > 0 && opts.DeepLevel <= node.depth {
		return 1, 0
	}
	for _, name := range names {
		// "all" option
		if !opts.All && strings.HasPrefix(name, ".") {
			continue
		}
		nnode := &Node{
			path:  filepath.Join(node.path, name),
			depth: node.depth + 1,
		}
		d, f := nnode.Visit(opts)
		// "dirs only" option
		if opts.DirsOnly && !nnode.IsDir() {
			continue
		}
		node.nodes = append(node.nodes, nnode)
		dirs, files = dirs+d, files+f
	}
	// Sorting
	if !opts.NoSort {
		node.sort(opts)
	}
	return dirs + 1, files
}

func (node *Node) sort(opts *Options) {
	var fn SortFunc
	switch {
	case opts.ModSort:
		fn = ModSort
	case opts.CTimeSort:
		fn = CTimeSort
	case opts.DirSort:
		fn = DirSort
	case opts.VerSort:
		fn = VerSort
	case opts.SizeSort:
		fn = SizeSort
	case opts.NameSort:
		fn = NameSort
	}
	if fn != nil {
		if opts.ReverSort {
			sort.Sort(sort.Reverse(ByFunc{node.nodes, fn}))
		} else {
			sort.Sort(ByFunc{node.nodes, fn})
		}
	}
}

func (node *Node) Print(indent string, opts *Options) {
	if node.err != nil {
		err := strings.Split(node.err.Error(), ": ")[1]
		fmt.Printf("%s [%s]\n", node.path, err)
		return
	}
	if !node.IsDir() {
		var props []string
		var stat = node.Sys().(*syscall.Stat_t)
		// inodes
		if opts.Inodes {
			props = append(props, fmt.Sprintf("%d", stat.Ino))
		}
		// device
		if opts.Device {
			props = append(props, fmt.Sprintf("%3d", stat.Dev))
		}
		// Mode
		if opts.FileMode {
			props = append(props, node.Mode().String())
		}
		// Owner/Uid
		if opts.ShowUid {
			uid := strconv.Itoa(int(stat.Uid))
			if u, err := user.LookupId(uid); err != nil {
				props = append(props, fmt.Sprintf("%-8s", uid))
			} else {
				props = append(props, fmt.Sprintf("%-8s", u.Username))
			}
		}
		// Gorup/Gid
		// TODO: support groupname
		if opts.ShowGid {
			gid := strconv.Itoa(int(stat.Gid))
			props = append(props, fmt.Sprintf("%-4s", gid))
		}
		// Size
		if opts.ByteSize || opts.UnitSize {
			var size string
			if opts.UnitSize {
				size = fmt.Sprintf("%4s", formatBytes(node.Size()))
			} else {
				size = fmt.Sprintf("%11d", node.Size())
			}
			props = append(props, size)
		}
		// Last modification
		if opts.LastMod {
			props = append(props, node.ModTime().Format("Jan 02 15:04"))
		}
		// Print properties
		if len(props) > 0 {
			fmt.Printf("[%s]  ", strings.Join(props, " "))
		}
	}
	// name/path
	var name string
	if node.depth == 0 || opts.FullPath {
		name = node.path
	} else {
		name = node.Name()
	}
	// Quotes
	if opts.Quotes {
		name = fmt.Sprintf("\"%s\"", name)
	}
	// Print file details
	fmt.Println(name)
	add := "│   "
	for i, nnode := range node.nodes {
		if opts.NoIndent {
			add = ""
		} else {
			if i == len(node.nodes)-1 {
				fmt.Printf(indent + "└── ")
				add = "    "
			} else {
				fmt.Printf(indent + "├── ")
			}
		}
		nnode.Print(indent+add, opts)
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
