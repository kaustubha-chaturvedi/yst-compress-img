// +build windows

package compressor

import (
	"fmt"
	"image"
)

func DecodeHeicAvif(path string) (image.Image, error) {
	return nil, fmt.Errorf("HEIC/AVIF decoding is not supported on Windows builds")
}
