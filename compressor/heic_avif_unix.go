// +build darwin linux

package compressor

import (
	"fmt"
	"image"
	"os"

	"github.com/vegidio/avif-go"
	"github.com/jdeng/goheif"
)

// decodeHeicAvif - Decode HEIC/AVIF images. This function tries to decode the image with gheif and go-avif. If both fail, it returns an error.
func DecodeHeicAvif(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file error: %w", err)
	}
	defer f.Close()

	if img, err := goheif.Decode(f); err == nil {
		return img, nil
	}

	if _, err := f.Seek(0, 0); err != nil {
		return nil, err
	}

	if img, err := avif.Decode(f); err == nil {
		return img, nil
	}

	return nil, fmt.Errorf("unsupported HEIC/HEIF/AVIF: %s", path)
}


