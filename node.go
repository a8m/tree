package tree

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/semaphore"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Node represent some node in the tree
// contains FileInfo, and its childs
type Node struct {
	os.FileInfo
	path   string
	depth  int
	err    error
	nodes  Nodes
	vpaths map[string]bool
}

// List of nodes
type Nodes []*Node

// To use this package programmatically, you must implement this
// interface.
// For example: PTAL on 'cmd/tree/tree.go'
type Fs interface {
	Stat(path string) (os.FileInfo, error)
	ReadDir(path string) ([]string, error)
}

// Options store the configuration for specific tree.
// Note, that 'Fs', and 'OutFile' are required (OutFile can be os.Stdout).
type Options struct {
	Fs      Fs
	OutFile io.Writer
	// List
	All        bool
	DirsOnly   bool
	FullPath   bool
	IgnoreCase bool
	FollowLink bool
	DeepLevel  int
	Pattern    string
	IPattern   string
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
	Colorize bool
	// Color defaults to ANSIColor()
	Color func(*Node, string) string

	// Internal data so we can do readdir()/stat() in parallel.
	sem *semaphore.Weighted
	res chan workerResult
}

func (opts *Options) color(node *Node, s string) string {
	f := opts.Color
	if f == nil {
		f = ANSIColor
	}
	return f(node, s)
}

// workerResult for go-ness
type workerResult struct {
	p *Node
	n *Node
	d int
	f int
}

// New get path and create new node(root).
func New(path string) *Node {
	return &Node{path: path, vpaths: make(map[string]bool)}
}

func newSubNode(opts *Options, node *Node, name string) (nnode *Node, dirs, files int) {
	nnode = &Node{
		path:   filepath.Join(node.path, name),
		depth:  node.depth + 1,
		vpaths: node.vpaths,
	}
	d, f := nnode.Visit(opts)
	if nnode.err == nil && !nnode.IsDir() {
		// "dirs only" option
		if opts.DirsOnly {
			return nil, 0, 0
		}
		var rePrefix string
		if opts.IgnoreCase {
			rePrefix = "(?i)"
		}
		// Pattern matching
		if opts.Pattern != "" {
			re, err := regexp.Compile(rePrefix + opts.Pattern)
			if err == nil && !re.MatchString(name) {
				return nil, 0, 0
			}
		}
		// IPattern matching
		if opts.IPattern != "" {
			re, err := regexp.Compile(rePrefix + opts.IPattern)
			if err == nil && re.MatchString(name) {
				return nil, 0, 0
			}
		}
	}

	return nnode, d, f
}

const semWeight = 128

// Visit all files under the given node.
func (node *Node) Visit(opts *Options) (dirs, files int) {
	goProcs := !opts.FollowLink && (semWeight > 0)

	// visited paths
	if !opts.FollowLink {
		node.vpaths = nil
	} else if path, err := filepath.Abs(node.path); err == nil {
		path = filepath.Clean(path)
		node.vpaths[path] = true
	}
	// stat
	fi, err := opts.Fs.Stat(node.path)
	if err != nil {
		node.err = err
		return
	}
	node.FileInfo = fi
	if !fi.IsDir() {
		return 0, 1
	}
	// increase dirs only if it's a dir, but not the root.
	if node.depth != 0 {
		dirs++
	}
	// DeepLevel option
	if opts.DeepLevel > 0 && opts.DeepLevel <= node.depth {
		return
	}
	names, err := opts.Fs.ReadDir(node.path)
	if err != nil {
		node.err = err
		return
	}
	node.nodes = make(Nodes, 0)
	var rwg sync.WaitGroup
	var fin chan workerResult
	if goProcs && node.depth == 0 {
		opts.sem = semaphore.NewWeighted(semWeight)
		opts.res = make(chan workerResult)
		rwg.Add(1)
		fin = make(chan workerResult)
		go func() {
			defer rwg.Done()
			defer close(fin)
			mdirs := 0
			mfiles := 0
			for val := range opts.res {
				val.p.nodes = append(val.p.nodes, val.n)
				mdirs, mfiles = mdirs+val.d, mfiles+val.f
			}
			fin <- workerResult{nil, node, mdirs, mfiles}
		}()
	}
	for i := range names {
		name := names[i]
		// fmt.Println("JDBG: beg:", name)
		// "all" option
		if !opts.All && strings.HasPrefix(name, ".") {
			continue
		}
		if goProcs && node.depth != 0 {
			if opts.sem.TryAcquire(2) {
				go func() {
					defer opts.sem.Release(2)
					nnode, d, f := newSubNode(opts, node, name)
					if nnode == nil {
						return
					}
					opts.res <- workerResult{node, nnode, d, f}
				}()
				continue
			}
		}
		// We ran out of semaphores, so just process the new node directly...
		nnode, d, f := newSubNode(opts, node, name)
		if nnode == nil {
			continue
		}
		if goProcs && node.depth != 0 {
			opts.res <- workerResult{node, nnode, d, f}
			continue
		}
		node.nodes = append(node.nodes, nnode)
		dirs, files = dirs+d, files+f
	}
	if goProcs && node.depth == 0 {
		opts.sem.Acquire(context.Background(), semWeight)
		close(opts.res)
		val := <-fin
		dirs += val.d
		files += val.f
		rwg.Wait()
	}

	// Sorting
	if !opts.NoSort {
		node.sort(opts)
	}
	return
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
	default:
		fn = NameSort // Default should be sorted, not unsorted.
	}
	if fn != nil {
		if opts.ReverSort {
			sort.Sort(sort.Reverse(ByFunc{node.nodes, fn}))
		} else {
			sort.Sort(ByFunc{node.nodes, fn})
		}
	}
}

// Path returns the Node's absolute path
func (node *Node) Path() string {
	return node.path
}

// Print nodes based on the given configuration.
func (node *Node) Print(opts *Options) { node.print("", opts) }

func dirRecursiveSize(opts *Options, node *Node) (size int64, err error) {
	if opts.DeepLevel > 0 && node.depth >= opts.DeepLevel {
		err = errors.New("Depth too high")
	}

	for _, nnode := range node.nodes {
		if nnode.err != nil {
			err = nnode.err
			continue
		}

		if !nnode.IsDir() {
			size += nnode.Size()
		} else {
			nsize, e := dirRecursiveSize(opts, nnode)
			size += nsize
			if e != nil {
				err = e
			}
		}
	}
	return
}

func (node *Node) print(indent string, opts *Options) {
	if node.err != nil {
		err := node.err.Error()
		if msgs := strings.Split(err, ": "); len(msgs) > 1 {
			err = msgs[1]
		}
		fmt.Printf("%s [%s]\n", node.path, err)
		return
	}
	if !node.IsDir() {
		var props []string
		ok, inode, device, uid, gid := getStat(node)
		// inodes
		if ok && opts.Inodes {
			props = append(props, fmt.Sprintf("%d", inode))
		}
		// device
		if ok && opts.Device {
			props = append(props, fmt.Sprintf("%3d", device))
		}
		// Mode
		if opts.FileMode {
			props = append(props, node.Mode().String())
		}
		// Owner/Uid
		if ok && opts.ShowUid {
			uidStr := strconv.Itoa(int(uid))
			if u, err := user.LookupId(uidStr); err != nil {
				props = append(props, fmt.Sprintf("%-8s", uidStr))
			} else {
				props = append(props, fmt.Sprintf("%-8s", u.Username))
			}
		}
		// Gorup/Gid
		// TODO: support groupname
		if ok && opts.ShowGid {
			gidStr := strconv.Itoa(int(gid))
			props = append(props, fmt.Sprintf("%-4s", gidStr))
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
			fmt.Fprintf(opts.OutFile, "[%s]  ", strings.Join(props, " "))
		}
	} else {
		var props []string
		// Size
		if opts.ByteSize || opts.UnitSize {
			var size string
			rsize, err := dirRecursiveSize(opts, node)
			if err != nil && rsize <= 0 {
				if opts.UnitSize {
					size = "????"
				} else {
					size = "???????????"
				}
			} else if opts.UnitSize {
				size = fmt.Sprintf("%4s", formatBytes(rsize))
			} else {
				size = fmt.Sprintf("%11d", rsize)
			}
			props = append(props, size)
		}
		// Print properties
		if len(props) > 0 {
			fmt.Fprintf(opts.OutFile, "[%s]  ", strings.Join(props, " "))
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
	// Colorize
	if opts.Colorize {
		name = opts.color(node, name)
	}
	// IsSymlink
	if node.Mode()&os.ModeSymlink == os.ModeSymlink {
		vtarget, err := os.Readlink(node.path)
		if err != nil {
			vtarget = node.path
		}
		targetPath, err := filepath.EvalSymlinks(node.path)
		if err != nil {
			targetPath = vtarget
		}
		fi, err := opts.Fs.Stat(targetPath)
		if opts.Colorize && fi != nil {
			vtarget = opts.color(&Node{FileInfo: fi, path: vtarget}, vtarget)
		}
		name = fmt.Sprintf("%s -> %s", name, vtarget)
		// Follow symbolic links like directories
		if opts.FollowLink {
			path, err := filepath.Abs(targetPath)
			if err == nil && fi != nil && fi.IsDir() {
				if _, ok := node.vpaths[filepath.Clean(path)]; !ok {
					inf := &Node{FileInfo: fi, path: targetPath}
					inf.vpaths = node.vpaths
					inf.Visit(opts)
					node.nodes = inf.nodes
				} else {
					name += " [recursive, not followed]"
				}
			}
		}
	}
	// Print file details
	// the main idea of the print logic came from here: github.com/campoy/tools/tree
	fmt.Fprintln(opts.OutFile, name)
	add := "│   "
	for i, nnode := range node.nodes {
		if opts.NoIndent {
			add = ""
		} else {
			if i == len(node.nodes)-1 {
				fmt.Fprintf(opts.OutFile, indent+"└── ")
				add = "    "
			} else {
				fmt.Fprintf(opts.OutFile, indent+"├── ")
			}
		}
		nnode.print(indent+add, opts)
	}
}

const (
	_        = iota // ignore first value by assigning to blank identifier
	KB int64 = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
)

// Convert bytes to human readable string. Like a 2 MB, 64.2 KB, 52 B
func formatBytes(i int64) (result string) {
	var n float64
	sFmt, eFmt := "%.01f", ""
	switch {
	case i > EB:
		eFmt = "E"
		n = float64(i) / float64(EB)
	case i > PB:
		eFmt = "P"
		n = float64(i) / float64(PB)
	case i > TB:
		eFmt = "T"
		n = float64(i) / float64(TB)
	case i > GB:
		eFmt = "G"
		n = float64(i) / float64(GB)
	case i > MB:
		eFmt = "M"
		n = float64(i) / float64(MB)
	case i > KB:
		eFmt = "K"
		n = float64(i) / float64(KB)
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
