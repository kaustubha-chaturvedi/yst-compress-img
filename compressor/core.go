package compressor

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
)

// image loading (support for JPEG/PNG/WebP)
func loadImageAnyFormat(path string) (image.Image, error) {
	img, _, err := LoadImage(path)
	if err == nil {
		return img, nil
	}

	return nil, err
}

// lossy quality compression 
func CompressQuality(in, out string, q int) error {
	img, err := loadImageAnyFormat(in)
	if err != nil {
		return err
	}
	return SaveJPEG(out, img, q)
}

// lossless: PNG or WebP
func CompressLossless(in, out string) error {
	img, format, err := LoadImage(in)
	if err == nil {
		if format == "png" {
			f, _ := os.Create(out)
			defer f.Close()
			enc := png.Encoder{CompressionLevel: png.BestCompression}
			return enc.Encode(f, img)
		}

		return webp.Save(out, img, &webp.Options{Lossless: true})
	}

	return webp.Save(out, img, &webp.Options{Lossless: true})
}

// resize
func CompressResize(in, out string, w, h, q int) error {
	img, err := loadImageAnyFormat(in)
	if err != nil {
		return err
	}

	resized := resize.Resize(
		uint(w),
		uint(h),
		img,
		resize.Lanczos3,
	)

	return SaveJPEG(out, resized, q)
}

// save jpeg
func SaveJPEG(path string, img image.Image, q int) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	return jpeg.Encode(out, img, &jpeg.Options{Quality: q})
}

// max size mode (degrade quality first, then degrade resolution)
func CompressToMaxSize(in, out string, maxBytes int64) error {
	img, err := loadImageAnyFormat(in)
	if err != nil {
		return err
	}

	temp := out + ".tmp.jpg"
	quality := 85
	scale := 1.0

	for {
		w := int(float64(img.Bounds().Dx()) * scale)
		h := int(float64(img.Bounds().Dy()) * scale)

		resized := resize.Resize(
			uint(w),
			uint(h),
			img,
			resize.Lanczos3,
		)

		if err := SaveJPEG(temp, resized, quality); err != nil {
			return err
		}

		if fileSize(temp) <= maxBytes {
			break
		}

		if quality > 35 {
			quality -= 5
			continue
		}

		scale *= 0.9
		if scale < 0.35 {
			break
		}
	}

	return os.Rename(temp, out)
}