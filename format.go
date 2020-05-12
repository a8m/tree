package tree

import (
	"fmt"
	"strings"
)

// KB = 1000 bytes
const (
	KB = 1 * 1000
	MB = KB * 1000
	GB = MB * 1000
	TB = GB * 1000
	PB = TB * 1000
	EB = PB * 1000
)

// Do it using KiB format, so KB = 1024 ...
const (
	_         = iota // ignore first value by assigning to blank identifier
	KiB int64 = 1 << (10 * iota)
	MiB
	GiB
	TiB
	PiB
	EiB
)

// round use like so: "%.1f", round(f, 0.1) or "%.0f", round(f, 1)
// Otherwise 9.9999 is < 10 but "%.1f" will give "10.0"
func round(x, unit float64) float64 {
	return float64(int64(x/unit+0.5)) * unit
}

// Convert bytes to human readable string. Like a 2 MB, 64.2 KB, 52 B
func formatBytes(i int64) (result string) {
	var n float64
	sFmt, eFmt := "%.01f", ""
	switch {
	case i >= EB:
		eFmt = "E"
		n = float64(i) / float64(EB)
	case i >= PB:
		eFmt = "P"
		n = float64(i) / float64(PB)
	case i >= TB:
		eFmt = "T"
		n = float64(i) / float64(TB)
	case i >= GB:
		eFmt = "G"
		n = float64(i) / float64(GB)
	case i >= MB:
		eFmt = "M"
		n = float64(i) / float64(MB)
	case i >= KB:
		eFmt = "K"
		n = float64(i) / float64(KB)
	default:
		sFmt = "%.0f"
		n = float64(i)
	}
	if eFmt != "" && round(n, 0.1) >= 10 {
		sFmt = "%.0f"
	}
	result = fmt.Sprintf(sFmt+eFmt, n)
	result = strings.Trim(result, " ")
	return
}

// Convert bytes to human readable string. Like a 2 MB, 64.2 KB, 52 B
func formatBytesKiB(i int64) (result string) {
	var n float64
	sFmt, eFmt := "%.01f", ""
	switch {
	case i >= EiB:
		eFmt = "E"
		n = float64(i) / float64(EiB)
	case i >= PiB:
		eFmt = "P"
		n = float64(i) / float64(PiB)
	case i >= TiB:
		eFmt = "T"
		n = float64(i) / float64(TiB)
	case i >= GiB:
		eFmt = "G"
		n = float64(i) / float64(GiB)
	case i >= MiB:
		eFmt = "M"
		n = float64(i) / float64(MiB)
	case i >= KiB:
		eFmt = "K"
		n = float64(i) / float64(KiB)
	default:
		sFmt = "%.0f"
		n = float64(i)
	}
	if eFmt != "" && round(n, 0.1) >= 10 {
		sFmt = "%.0f"
	}
	result = fmt.Sprintf(sFmt+eFmt, n)
	result = strings.Trim(result, " ")
	return
}
