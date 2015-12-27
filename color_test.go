package tree

import (
	"fmt"
	"testing"
)

func TestANSIColor(t *testing.T) {
	var f1 file
	f1.name = "hello"
	n1 := &Node{FileInfo: f1}
	fmt.Printf(ANSIColor(n1, f1.name) + "\n")
}
