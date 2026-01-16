package compressor

import (
	"image"
	"math"

	"github.com/nfnt/resize"
)

/*
Smart Auto Mode Logic
---------------------
1. Inspect input image:
   - Dimensions
   - Aspect ratio
   - Average brightness (sampled for speed)

2. If file >6MB:
      Downscale (0.55–0.75x depending on aspect ratio)

3. Adaptive quality:
      Dark images - lower quality (dark compress worse)
      Bright images - higher quality

4. Apply resize & JPEG encode.

This attempts to produce the best compression/quality tradeoff with ZERO user config.
*/

func CompressAuto(in, out string) error {
	size := fileSize(in)

	img, err := loadImageAnyFormat(in)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	aspect := float64(w) / float64(h)
	brightness := avgBrightnessSample(img)

	if size > 6*1024*1024 {
		scale := 0.6

		if aspect > 2.1 {
			scale = 0.75
		}

		newW := uint(float64(w) * scale)
		resized := resize.Resize(newW, 0, img, resize.Lanczos3)
		q := autoQuality(brightness)

		return SaveJPEG(out, resized, q)
	}

	return SaveJPEG(out, img, autoQuality(brightness))
}

// histogram brightness (my python brain could come up with this only)
func avgBrightnessSample(img image.Image) float64 {
	b := img.Bounds()

	var sum float64
	var count float64

	step := 50

	for y := b.Min.Y; y < b.Max.Y; y += step {
		for x := b.Min.X; x < b.Max.X; x += step {
			r, g, bl, _ := img.At(x, y).RGBA()
			avg := (float64(r) + float64(g) + float64(bl)) / (3 * 65535)
			sum += math.Sqrt(avg)
			count++
		}
	}

	if count == 0 {
		return 0.5
	}

	return sum / count
}

// select quality
func autoQuality(br float64) int {
	switch {
	case br < 0.25:
		return 100
	case br < 0.45:
		return 80
	default:
		return 85
	}
}
