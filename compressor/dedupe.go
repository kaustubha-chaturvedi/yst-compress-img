package compressor

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"image"
	"os"
	"sync"
)

// image duplicate detection with hashig
var (
	seenHashes = make(map[string]string)
	hashMu     sync.Mutex
)

func HashImage(img image.Image) (string, error) {
	b := img.Bounds()
	w := b.Dx()
	h := b.Dy()

	buf := make([]byte, w*h*4)
	idx := 0

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl, a := img.At(x, y).RGBA()

			buf[idx+0] = byte(r >> 8)
			buf[idx+1] = byte(g >> 8)
			buf[idx+2] = byte(bl >> 8)
			buf[idx+3] = byte(a >> 8)

			idx += 4
		}
	}

	hh := sha1.Sum(buf)
	return hex.EncodeToString(hh[:]), nil
}

func IsDuplicateImage(path string) (bool, string, error) {
	img, err := loadImageAnyFormat(path)
	if err != nil {
		return false, "", err
	}

	hash, err := HashImage(img)
	if err != nil {
		return false, "", err
	}

	hashMu.Lock()
	defer hashMu.Unlock()

	if orig, exists := seenHashes[hash]; exists {
		fmt.Printf("[dedupe] skipped duplicate: %s (original: %s)\n", path, orig)
		return true, orig, nil
	}

	seenHashes[hash] = path
	return false, "", nil
}

func ResetDedupe() {
	hashMu.Lock()
	defer hashMu.Unlock()
	seenHashes = make(map[string]string)
}

func HashFileIfPossible(path string) (string, error) {
	img, err := loadImageAnyFormat(path)
	if err != nil {
		return "", err
	}
	return HashImage(img)
}

func IsFileEmptyOrUnreadable(path string) bool {
	info, err := os.Stat(path)
	if err != nil || info.Size() == 0 {
		return true
	}
	return false
}
