package gocrop

import (
	"fmt"
	"image"
	"image/draw"
	"io"
	"os"
	"path"
	"sync"
)

// Cropper crops and saves images.
type Cropper struct {
	threshold     uint32
	outPrefix     string
	outSuffix     string
	outDir        string
	skipUnchanged bool
	padding       int
	enumerate     bool
	num           int
	numMu         sync.Mutex
}

// NewCropper creates an instance of *Cropped with provided options,
// returns error if any option fails.
//
// Default Cropper with no options:
//
// - has alpha threshold of 0 and no padding
//
// - saves images under the same name in the same directory as the source image (file will be overwritten).
// If cropping made no changes it still saves the result.
func NewCropper(options ...CropperOption) (*Cropper, error) {
	c := &Cropper{}

	for _, opt := range options {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Crop takes a *Croppable and returns a cropped version of it and a bool flag indicating if cropping was done.
// If cropping would make no changes to given *Croppable, the provided *Croppable is returned back with false flag.
func (i *Cropper) Crop(croppable *Croppable) (*Croppable, bool) {
	rect := i.Rect(croppable.Image)

	if i.padding == 0 {
		if rect.Size().Eq(croppable.Image.Bounds().Size()) {
			return croppable, false
		}

		return croppable.With(croppable.Image.SubImage(rect).(CroppableImage)), true
	}

	// if rect cuts deep enough it's possible to extend it and crop the image with padded rect
	if rect.Min.X >= i.padding &&
		rect.Min.Y >= i.padding &&
		croppable.Image.Bounds().Dx()-rect.Max.X >= i.padding &&
		croppable.Image.Bounds().Dy()-rect.Max.X >= i.padding {
		rect.Min.X -= i.padding
		rect.Min.Y -= i.padding
		rect.Max.X += i.padding
		rect.Max.Y += i.padding

		return croppable.With(croppable.Image.SubImage(rect).(CroppableImage)), true
	}

	// if rect is too small create new empty image with proper size and draw the cropped image onto it
	bg := image.NewRGBA(image.Rect(0, 0, rect.Dx()+(2*i.padding), rect.Dy()+(2*i.padding)))
	croppedRect := image.Rect(i.padding, i.padding, i.padding+rect.Dx(), i.padding+rect.Dy())

	draw.Draw(bg, croppedRect, croppable.Image.SubImage(rect), image.Point{rect.Min.X, rect.Min.Y}, draw.Src)

	return croppable.With(bg), true
}

// Save saves the croppable, creates a directory if it doesn't exist.
func (i *Cropper) Save(c *Croppable) error {
	if i.outDir != "" {
		if err := os.MkdirAll(i.outDir, os.ModePerm); err != nil {
			return err
		}
	}

	return i.save(c)
}

func (i *Cropper) save(c *Croppable) error {
	dir, name, ext := dirFileExt(c.Path)

	if i.outDir != "" {
		dir = i.outDir
	}

	num := ""
	if i.enumerate {
		num = fmt.Sprintf("_%d", i.enum())
	}

	name = i.outPrefix + name + num + i.outSuffix + "." + ext
	outPath := path.Join(dir, name)

	if err := saveImage(outPath, c.Image, c.Encode); err != nil {
		return err
	}

	return nil
}

func (i *Cropper) enum() int {
	i.numMu.Lock()
	defer i.numMu.Unlock()
	n := i.num
	i.num += 1

	return n
}

// CropAndSave crops an image and saves it, output directory will be created if it does not exist.
// Error will be returned if the image was not saved successfully.
func (i *Cropper) CropAndSave(croppable *Croppable) error {
	if i.outDir != "" {
		if err := os.MkdirAll(i.outDir, os.ModePerm); err != nil {
			return err
		}
	}

	cropped, ok := i.Crop(croppable)
	if !ok && i.skipUnchanged {
		return nil
	}

	if err := i.save(cropped); err != nil {
		return fmt.Errorf("error saving image: %w", err)
	}

	return nil
}

// Rect returns the cropping rectangle of the image, does not include padding.
func (i *Cropper) Rect(img image.Image) image.Rectangle {
	rect := img.Bounds()

	min := image.Point{-1, -1}
	max := image.Point{-1, -1}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		for y := 0; y < rect.Dy(); y++ {
			for x := 0; x < rect.Dx(); x++ {
				pixel := img.At(x, y)

				_, _, _, alpha := pixel.RGBA()

				if alpha > i.threshold {
					if min.X == -1 || x < min.X {
						min.X = x
					}

					if min.Y == -1 || y < min.Y {
						min.Y = y
					}

					break
				}
			}
		}
	}()

	go func() {
		defer wg.Done()
		for y := rect.Dy() - 1; y >= 0; y-- {
			for x := rect.Dx() - 1; x >= 0; x-- {
				pixel := img.At(x, y)

				_, _, _, alpha := pixel.RGBA()

				if alpha > i.threshold {
					if x > max.X {
						max.X = x + 1
					}

					if y > max.Y {
						max.Y = y + 1
					}

					break
				}
			}
		}
	}()

	wg.Wait()

	if min.Eq(image.Point{-1, -1}) {
		min = image.Point{}
	}

	if max.Eq(image.Point{-1, -1}) {
		max = image.Point{rect.Dx(), rect.Dy()}
	}

	return image.Rectangle{
		Min: min,
		Max: max,
	}
}

// Croppable holds the directory of the image, it's name and format in separate string fields (for easier manipulation of the output destination),
// the core image itself and a proper encoder function for encoding the image.
type Croppable struct {
	Path   string
	Image  CroppableImage
	Decode func(r io.Reader) (image.Image, error)
	Encode func(w io.Writer, m image.Image) error
}

// Load validates if given image format is supported, if so
// creates a *Croppable and calls it's Load() method.
func Load(path string) (*Croppable, error) {
	_, _, ext := dirFileExt(path)

	coder, ok := imageCoders[ext]
	if !ok {
		return nil, ErrUnsupportedFormat
	}

	c := &Croppable{
		Path:   path,
		Decode: coder.decode,
		Encode: coder.encode,
	}

	if err := c.Load(); err != nil {
		return nil, err
	}

	return c, nil
}

// Load loads the image of the croppable using it's decoder
// returns an error if image was not successfully decoded or image is not croppable.
func (c *Croppable) Load() error {
	file, err := os.Open(c.Path)
	if err != nil {
		return err
	}

	defer file.Close()

	img, err := c.Decode(file)
	if err != nil {
		return fmt.Errorf("%s: %w", err.Error(), ErrImageLoadFailed)
	}

	croppableImg, ok := img.(CroppableImage)
	if !ok {
		return ErrImageUncroppable
	}

	c.Image = croppableImg

	return nil
}

// With returns a copy of current croppable with Image set to provided image.
func (c *Croppable) With(ci CroppableImage) *Croppable {
	return &Croppable{
		Path:   c.Path,
		Image:  ci,
		Decode: c.Decode,
		Encode: c.Encode,
	}
}

// Croppable image is an extension of image.Image interface to ensure the image is croppable.
type CroppableImage interface {
	image.Image
	SubImage(r image.Rectangle) image.Image
}

type CropperOption func(*Cropper) error

// WithThreshold sets the alpha channel threshold for cropping.
// Only pixels that satisfy the condition pixelAlpha > threshold will be used in the process of finding a cropping rectangle.
func WithThreshold(threshold uint32) CropperOption {
	return func(c *Cropper) error {
		c.threshold = threshold
		return nil
	}
}

// WithPadding sets the number of pixels to add in each direction around the cropped image.
// If the cropped output is the size 25x25px, with 5px of padding it will be 35x35px with the cropped element centered.
func WithPadding(padding int) CropperOption {
	return func(c *Cropper) error {
		c.padding = padding
		return nil
	}
}

// WithOutPrefix adds a prefix to an image name.
// Given prefix "cropped_" and image file name "image1.png", the output image will be named "cropped_image1.png".
func WithOutPrefix(prefix string) CropperOption {
	return func(c *Cropper) error {
		c.outPrefix = prefix
		return nil
	}
}

// WithOutSuffix adds a suffix to an image name.
// Given suffix "_cropped" and image file name "image1.png", the output image will be named "image1_cropped.png".
func WithOutSuffix(suffix string) CropperOption {
	return func(c *Cropper) error {
		c.outSuffix = suffix
		return nil
	}
}

// WithOutDir sets the output directory for cropped images.
func WithOutDir(dir string) CropperOption {
	return func(c *Cropper) error {
		c.outDir = dir
		return nil
	}
}

// WithSkipUnchanged if set to true output image won't be saved if cropping made no changes to original image.
func WithSkipUnchanged(skip bool) CropperOption {
	return func(c *Cropper) error {
		c.skipUnchanged = skip
		return nil
	}
}

// WithEnumerate enables enumeration of output images.
// "_n" is appended to it's name (before the suffix), n is an integer incremented by 1 each time an image is saved.
// Given n of 5, image file name of "image1.png", the output image will be named "image1_5.png".
func WithEnumerate(enumerate bool) CropperOption {
	return func(c *Cropper) error {
		c.enumerate = enumerate
		return nil
	}
}
