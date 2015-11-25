package node

import (
	"os"
	"syscall"
)

type ByFunc struct {
	Nodes
	Fn SortFunc
}

func (b ByFunc) Less(i, j int) bool {
	return b.Fn(b.Nodes[i].FileInfo, b.Nodes[j].FileInfo)
}

type SortFunc func(f1, f2 os.FileInfo) bool

func ModSort(f1, f2 os.FileInfo) bool {
	return f1.ModTime().Before(f2.ModTime())
}

func CTimeSort(f1, f2 os.FileInfo) bool {
	s1, s2 := f1.Sys().(*syscall.Stat_t), f2.Sys().(*syscall.Stat_t)
	return s1.Ctimespec.Sec < s2.Ctimespec.Sec
}

func DirSort(f1, f2 os.FileInfo) bool {
	return f1.IsDir() && !f2.IsDir()
}
