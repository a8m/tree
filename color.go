package tree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const Escape = "\x1b"
const (
	Reset int = 0
	Bold  int = 1
	Black int = iota + 28
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// ANSIColor
func ANSIColor(node *Node, s string) string {
	var color int
	var mode = node.Mode()
	var ext = filepath.Ext(node.Name())
	switch {
	case contains([]string{".bat", ".btm", ".cmd", ".com", ".dll", ".exe"}, ext):
		color = Green
	case contains([]string{".arj", ".bz2", ".deb", ".gz", ".lzh", ".rpm",
		".tar", ".taz", ".tb2", ".tbz2", ".tbz", ".tgz", ".tz", ".tz2", ".z",
		".zip", ".zoo"}, ext):
		color = Red
	case contains([]string{".asf", ".avi", ".bmp", ".flac", ".gif", ".jpg",
		"jpeg", ".m2a", ".m2v", ".mov", ".mp3", ".mpeg", ".mpg", ".ogg", ".ppm",
		".rm", ".tga", ".tif", ".wav", ".wmv",
		".xbm", ".xpm"}, ext):
		color = Magenta
	case node.IsDir() || mode&os.ModeDir != 0:
		color = Blue
	default:
		// IsSymlink
		if node.Mode()&os.ModeSymlink == os.ModeSymlink {
			// IsOrphan
			if _, err := filepath.EvalSymlinks(node.path); err != nil {
				// Error link color
				return fmt.Sprintf("%s[40;%d;%dm%s%s[%dm", Escape, Bold, Red, s, Escape, Reset)
			} else {
				color = Cyan
			}
		}
		// IsSocket
		if node.Mode()&os.ModeSocket == os.ModeSocket {
			return fmt.Sprintf("%s[40;%d;%dm%s%s[%dm", Escape, Bold, Magenta, s, Escape, Reset)
		}
		// IsFifo
		if node.Mode()&os.ModeNamedPipe == os.ModeNamedPipe {
			return fmt.Sprintf("%s[40;%dm%s%s[%dm", Escape, Yellow, s, Escape, Reset)
		}
		// IsBlk - a block special file (a device like a disk)
		if node.Mode()&os.ModeDevice == os.ModeDevice {
			return fmt.Sprintf("%s[40;%d;01m%s%s[%dm", Escape, Yellow, s, Escape, Reset)
		}
		// IsChr
		if node.Mode()&os.ModeCharDevice == os.ModeCharDevice {
			return fmt.Sprintf("%s[40;%d;01m%s%s[%dm", Escape, Yellow, s, Escape, Reset)
		}
		// IsExec
		// Refactor after write some tests
		// http://stackoverflow.com/questions/13098620/using-stat-to-check-if-a-file-is-executable-in-c
		//if node.Mode()&(syscall.S_IXUSR|syscall.S_IXGRP|syscall.S_IXOTH) != 0 {
		//	return fmt.Sprintf("%s[01;%dm%s%s[%dm", Escape, Green, s, Escape, Reset)
		//}
	}
	return fmt.Sprintf("%s[%d;%dm%s%s[%dm", Escape, Bold, color, s, Escape, Reset)
}

// case-insensitive contains helper
func contains(slice []string, str string) bool {
	for _, val := range slice {
		if val == strings.ToLower(str) {
			return true
		}
	}
	return false
}

// TODO: HTMLColor
