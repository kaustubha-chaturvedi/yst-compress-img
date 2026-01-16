package compressor

import (
	"fmt"
	"image"
	"os"
	"strconv"
	"strings"
)

// parse size
func ParseSize(s string) (int64, error) {
	s = strings.ToLower(strings.TrimSpace(s))

	switch {
	case strings.HasSuffix(s, "kb"):
		n, err := strconv.Atoi(strings.TrimSuffix(s, "kb"))
		if err != nil {
			return 0, err
		}
		return int64(n) * 1024, nil

	case strings.HasSuffix(s, "mb"):
		n, err := strconv.Atoi(strings.TrimSuffix(s, "mb"))
		if err != nil {
			return 0, err
		}
		return int64(n) * 1024 * 1024, nil

	default:
		return 0, fmt.Errorf("invalid size format: use 500kb or 1mb")
	}
}

// file size
func fileSize(path string) int64 {
	info, _ := os.Stat(path)
	if info == nil {
		return 0
	}
	return info.Size()
}

// loadImage - Wraps image.Decode + format detection.
func LoadImage(path string) (image.Image, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	img, format, err := image.Decode(f)
	return img, format, err
}