package main

import (
	"os"
)

type Info struct {
	os.FileInfo
	path  string
	err   error
	infos []Info
}
