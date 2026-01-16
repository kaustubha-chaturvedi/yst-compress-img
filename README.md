# yst-compress-img

YST plugin for image compression.

## Usage with yeast:
  `yst compress-img <image | -b directory> [flags]`

## Supported Flags:
  -a, --auto              smart mode
  -b, --batch string      directory to compress recursively or non-recursively
  -h, --height int        resize height
      --help              show help
  -l, --lossless          lossless mode
  -m, --max-size string   target max size (e.g. 500kb, 1mb)
  -o, --output string     output file path (default: <file name>_compressed.<ext>)
  -p, --parallel int      parallel workers for batch mode (default 4)
  -q, --quality int       JPEG/WebP quality (1-100) (default 85)
  -r, --recursive         scan directory recursively in batch mode
  -w, --width int         resize width

## Disclaimer

Some file may rather get bigger in auto mode compression i would love if someone want to impreove my auto mode algo