package status

import (
	"fmt"
	"strings"
)

func ConvertToUnit(b uint64, e string) float64 {
	if strings.Contains(e, "i") {
		return ByteCountIECConvert(b, e)
	}
	return ByteCountSIConvert(b, e)
}

func ByteCountSIConvert(b uint64, e string) float64 {
	const unit = 1000
	if e == "B" {
		return float64(b)
	}
	div, exp := uint64(unit), 0
	ex := fmt.Sprintf("%cB", "kMGTPE"[exp])
	for n := b / unit; e != ex && exp < 6; n /= unit {
		div *= unit
		exp++
		ex = fmt.Sprintf("%cB", "kMGTPE"[exp])
	}
	return float64(b) / float64(div)
}

func ByteCountIECConvert(b uint64, e string) float64 {
	const unit = 1024
	if e == "B" {
		return float64(b)
	}
	div, exp := uint64(unit), 0
	ex := fmt.Sprintf("%ciB", "KMGTPE"[exp])
	for n := b / unit; e != ex && exp < 6; n /= unit {
		div *= unit
		exp++
		ex = fmt.Sprintf("%ciB", "KMGTPE"[exp])
	}
	return float64(b) / float64(div)
}
