package node

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
	Yello
	Blue
	Magenta
	Cyan
	White
)

func ansiFormat(color int, s string) string {
	return fmt.Sprintf("%s[%d;%dm%s%s[%dm", Escape, Bold, color, s, Escape, Reset)
}

// ANSIColor
func ANSIColor(node os.FileInfo, s string) string {
	var color int
	switch ext := filepath.Ext(node.Name()); strings.ToLower(ext) {
	case ".bat", ".btm", ".cmd", ".com", ".dll", ".exe":
		color = Green
	case ".arj", ".bz2", ".deb", ".gz", ".lzh", ".rpm", ".tar", ".taz", ".tb2", ".tbz2",
		".tbz", ".tgz", ".tz", ".tz2", ".z", ".zip", ".zoo":
		color = Red
	case ".asf", ".avi", ".bmp", ".flac", ".gif", ".jpg", "jpeg", ".m2a", ".m2v", ".mov",
		".mp3", ".mpeg", ".mpg", ".ogg", ".ppm", ".rm", ".tga", ".tif", ".wav", ".wmv",
		".xbm", ".xpm":
		color = Magenta
	default:
		if node.IsDir() {
			color = Blue
		}
		if node.Mode()&os.ModeSymlink == os.ModeSymlink {
			color = Cyan
		}
	}
	return ansiFormat(color, s)
}

// TODO: HTMLColor
