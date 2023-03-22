package gocrop

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
	"png": {
		decode: png.Decode,
		encode: png.Encode,
	},
	"gif": {
		decode: gif.Decode,
		encode: func(w io.Writer, m image.Image) error {
			return gif.Encode(w, m, nil)
		},
	},
	"tiff": {
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
	name, ext = fileExt(filepath.Base(fp))

	return
}

func fileExt(fileName string) (name, extension string) {
	split := strings.Split(fileName, ".")
	if len(split) < 2 {
		return fileName, ""
	}

	return strings.Join(split[0:len(split)-1], "."), split[len(split)-1]
}
