package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/kaustubha-chaturvedi/yst-compress-img/compressor"
)

var (
	quality  int
	maxSize  string
	lossless bool
	width    int
	height   int
	autoMode bool
	output   string

	batchDir  string
	recursive bool
	parallel  int
)

var compressCmd = &cobra.Command{
	Use:   "compress <image | directory>",
	Short: "Compress a single image or an entire directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 && batchDir == "" {
			return fmt.Errorf("missing image or directory path")
		}

		if batchDir != "" {
			return compressor.BatchCompress(
				batchDir,
				quality,
				maxSize,
				lossless,
				width,
				height,
				autoMode,
				recursive,
				parallel,
			)
		}

		in := args[0]

		if _, err := os.Stat(in); err != nil {
			return fmt.Errorf("file not found: %s", in)
		}

		if output == "" {
			ext := filepath.Ext(in)
			output = strings.TrimSuffix(in, ext) + "_compressed" + ext
		}

		switch {
			case autoMode:
				return compressor.CompressAuto(in, output)

			case lossless:
				return compressor.CompressLossless(in, output)

			case maxSize != "":
				sz, err := compressor.ParseSize(maxSize)
				if err != nil {
					return err
				}
				return compressor.CompressToMaxSize(in, output, sz)

			case width > 0 || height > 0:
				return compressor.CompressResize(in, output, width, height, quality)

			default:
				return compressor.CompressQuality(in, output, quality)
		}
	},
}

func init() {
	compressCmd.Flags().IntVarP(&quality, "quality", "q", 85, "JPEG/WebP quality (1-100)")
	compressCmd.Flags().StringVarP(&maxSize, "max-size", "ms", "", "target max size (e.g. 500kb, 1mb)")
	compressCmd.Flags().BoolVarP(&lossless, "lossless", "l", true, "lossless mode")
	compressCmd.Flags().IntVarP(&width, "width", "w", 0, "resize width")
	compressCmd.Flags().IntVarP(&height, "height", "h", 0, "resize height")
	compressCmd.Flags().BoolVarP(&autoMode, "auto", "a", false, "smart mode")
	compressCmd.Flags().StringVarP(&output, "output", "o", "", "output file path (default: <input>_compressed.<ext>)")

	compressCmd.Flags().StringVarP(&batchDir, "batch", "b", "", "directory to compress recursively or non-recursively")
	compressCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "scan directory recursively in batch mode")
	compressCmd.Flags().IntVarP(&parallel, "parallel", "p", 4, "parallel workers for batch mode")
}
