package compressor

import (
	"slices"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var supportedExt = []string{
	".jpg", ".jpeg", ".png", ".webp",
	".heic", ".heif", ".avif",
}

func isImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return slices.Contains(supportedExt, ext)
}

func CollectImages(dir string, recursive bool) ([]string, error) {
	var files []string

	if recursive {
		err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				return nil
			}
			if isImageFile(path) && !IsFileEmptyOrUnreadable(path) {
				files = append(files, path)
			}
			return nil
		})

		if err != nil {
			return nil, err
		}
	} else {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			path := filepath.Join(dir, e.Name())
			if isImageFile(path) && !IsFileEmptyOrUnreadable(path) {
				files = append(files, path)
			}
		}
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no image files found in %s", dir)
	}

	return files, nil
}