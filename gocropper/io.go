package gocropper

import (
	"errors"
	"image"
	"image/gif"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/tiff"
)

var ErrUnsupportedFormat = errors.New("unsupported format")
var ErrImageUncroppable = errors.New("image does not support cropping")
var ErrImageLoadFailed = errors.New("unable to load image")

type imageCoder struct {
	decode func(r io.Reader) (image.Image, error)
	encode func(w io.Writer, m image.Image) error
}

var imageCoders = map[string]imageCoder{
	".png": {
		decode: png.Decode,
		encode: png.Encode,
	},
	".gif": {
		decode: gif.Decode,
		encode: func(w io.Writer, m image.Image) error {
			return gif.Encode(w, m, nil)
		},
	},
	".tiff": {
		decode: tiff.Decode,
		encode: func(w io.Writer, m image.Image) error {
			return tiff.Encode(w, m, nil)
		},
	},
}

func saveImage(fp string, img image.Image, encode func(w io.Writer, m image.Image) error) error {
	fd, err := os.Create(fp)
	if err != nil {
		return err
	}

	defer fd.Close()

	return encode(fd, img)
}

func dirFileExt(fp string) (dir, name, ext string) {
	dir = filepath.Dir(fp)
	ext = filepath.Ext(fp)
	name = strings.TrimSuffix(filepath.Base(fp), ext)

	return
}

func fileExt(fileName string) (name, ext string) {
	ext = filepath.Ext(fileName)
	name = strings.TrimSuffix(filepath.Base(fileName), ext)

	return
}
