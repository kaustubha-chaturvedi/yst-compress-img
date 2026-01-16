package compressor

import (
	"fmt"
	"os"
	"sync"
)

func BatchCompress(
	dir string,
	quality int,
	maxSize string,
	lossless bool,
	width int,
	height int,
	auto bool,
	recursive bool,
	workers int,
) error {

	paths, err := CollectImages(dir, recursive)
	if err != nil {
		return err
	}

	fmt.Printf("[batch] found %d images in %s\n", len(paths), dir)

	var totalBytes int64
	for _, f := range paths {
		info, err := os.Stat(f)
		if err == nil {
			totalBytes += info.Size()
		}
	}

	ResetDedupe()

	progress := NewProgress(len(paths))

	jobs := make(chan string)
	var wg sync.WaitGroup

	for i := range workers {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for path := range jobs {
				dup, orig, err := IsDuplicateImage(path)
				if err == nil && dup {
					progress.Update()
					fmt.Printf("[dup] %s (original %s)\n", path, orig)
					continue
				}

				err = processSingle(
					path,
					quality,
					maxSize,
					lossless,
					width,
					height,
					auto,
				)
				if err != nil {
					fmt.Printf("[err] %s: %v\n", path, err)
				}

				progress.Update()
			}
		}(i)
	}

	go func() {
		for _, f := range paths {
			jobs <- f
		}
		close(jobs)
	}()

	wg.Wait()

	fmt.Println("[batch] done")
	return nil
}

func processSingle(
	path string,
	quality int,
	maxSize string,
	lossless bool,
	width int,
	height int,
	auto bool,
) error {

	out := outputPathFor(path)

	switch {
	case auto:
		return CompressAuto(path, out)

	case lossless:
		return CompressLossless(path, out)

	case maxSize != "":
		sz, err := ParseSize(maxSize)
		if err != nil {
			return err
		}
		return CompressToMaxSize(path, out, sz)

	case width > 0 || height > 0:
		return CompressResize(path, out, width, height, quality)

	default:
		return CompressQuality(path, out, quality)
	}
}