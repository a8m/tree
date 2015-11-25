package node

import (
	"os"
)

type ByFunc struct {
	Nodes
	Fn LessFn
}

func (b ByFunc) Less(i, j int) bool {
	return b.Fn(b.Nodes[i].FileInfo, b.Nodes[j].FileInfo)
}

type LessFn func(f1, f2 os.FileInfo) bool
