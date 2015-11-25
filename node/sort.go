package node

import (
	"os"
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
