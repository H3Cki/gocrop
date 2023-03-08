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

type Cropper struct {
	threshold     uint32
	outPrefix     string
	outSuffix     string
	outDir        string
	skipUnchanged bool
	padding       int
}

func NewCropper(options ...CropperOption) (*Cropper, error) {
	c := &Cropper{}

	for _, opt := range options {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Crop takes a *Croppable and crops it. If cropped result differs from *Croppable source image
// cropped image is returned as non-nil image.Image with a true bool flag. If cropped result
// is the same as *Croppable source image function returns (nil, false)
func (i *Cropper) Crop(croppable *Croppable) (image.Image, bool) {
	cropped, ok := i.crop(croppable)

	return cropped, ok
}

// CropAndSave takes a *Croppable, calls *Cropper.Crop(*Croppable) and attempts to save the cropped result.
// If cropped result is the same as *Croppable source image and WithSkipUnchanged option was set to true,
// image will not be saved.
func (i *Cropper) CropAndSave(croppable *Croppable) error {
	if i.outDir != "" {
		if err := os.MkdirAll(i.outDir, os.ModePerm); err != nil {
			return err
		}
	}

	cropped, ok := i.crop(croppable)
	if !ok && i.skipUnchanged {
		return nil
	}

	if err := i.save(croppable, cropped); err != nil {
		return err
	}

	return nil
}

func (i *Cropper) crop(c *Croppable) (image.Image, bool) {
	rect := i.Rect(c.Cropper)

	var cropped image.Image

	// if rect is small enough it's possible to extend it and crop the image with padded rect
	if i.padding != 0 &&
		rect.Min.X >= i.padding &&
		rect.Min.Y >= i.padding &&
		c.Cropper.Bounds().Dx()-rect.Max.X >= i.padding &&
		c.Cropper.Bounds().Dy()-rect.Max.X >= i.padding {
		rect.Min.X -= i.padding
		rect.Min.Y -= i.padding
		rect.Max.X += i.padding
		rect.Max.Y += i.padding

		cropped = c.Cropper.SubImage(rect)
	} else if i.padding != 0 {
		// if rect is too small create new empty image with proper size and draw the cropped image onto it
		cropped = image.NewRGBA(image.Rect(0, 0, rect.Dx()+(2*i.padding), rect.Dy()+(2*i.padding)))
		croppedRect := image.Rect(i.padding, i.padding, i.padding+rect.Dx(), i.padding+rect.Dy())

		draw.Draw(cropped.(draw.Image), croppedRect, c.Cropper.SubImage(rect), image.Point{rect.Min.X, rect.Min.Y}, draw.Src)
	} else {
		if rect.Size().Eq(c.Cropper.Bounds().Size()) {
			return c.Cropper, false
		}

		cropped = c.Cropper.SubImage(rect)
	}

	return cropped, true
}

func (i *Cropper) save(c *Croppable, img image.Image) error {
	dir := c.Dir
	if i.outDir != "" {
		dir = i.outDir
	}

	name := i.outPrefix + c.Name + i.outSuffix + "." + c.Format
	outPath := path.Join(dir, name)

	if err := saveImage(outPath, img, c.Encode); err != nil {
		return fmt.Errorf("error saving %s: %w", outPath, err)
	}

	return nil
}

func (i *Cropper) Rect(img image.Image) image.Rectangle {
	rect := img.Bounds()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	top, left := rect.Max.Y, rect.Max.X

	// Seek min
	go func() {
		defer wg.Done()

		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for x := rect.Min.X; x < rect.Max.X; x++ {
				pixel := img.At(x, y)

				_, _, _, alpha := pixel.RGBA()

				if alpha > i.threshold {
					if x < left {
						left = x
					}

					if y < top {
						top = y
					}

					break
				}
			}
		}
	}()

	bottom, right := 0, 0

	// Seek max
	go func() {
		defer wg.Done()

		for y := rect.Max.Y; y > 0; y-- {
			for x := rect.Max.X; x > 0; x-- {
				pixel := img.At(x, y)

				_, _, _, alpha := pixel.RGBA()

				if alpha > i.threshold {
					if x > right {
						right = x + 1
					}

					if y > bottom {
						bottom = y + 1
					}

					break
				}
			}
		}
	}()

	wg.Wait()

	if right > rect.Max.X {
		right = rect.Max.X
	}

	if bottom > rect.Max.Y {
		bottom = rect.Max.Y
	}

	return image.Rectangle{
		Min: image.Point{
			X: left,
			Y: top,
		},
		Max: image.Point{
			X: right,
			Y: bottom,
		},
	}
}

type Croppable struct {
	Dir, Name, Format string
	Image             image.Image
	Cropper           CroppableImage
	Encode            func(w io.Writer, m image.Image) error
}

type CroppableImage interface {
	image.Image
	SubImage(r image.Rectangle) image.Image
}

type CropperOption func(*Cropper) error

func WithThreshold(threshold uint32) CropperOption {
	return func(c *Cropper) error {
		c.threshold = threshold
		return nil
	}
}

func WithPadding(padding int) CropperOption {
	return func(c *Cropper) error {
		c.padding = padding
		return nil
	}
}

func WithOutPrefix(prefix string) CropperOption {
	return func(c *Cropper) error {
		c.outPrefix = prefix
		return nil
	}
}

func WithOutSuffix(suffix string) CropperOption {
	return func(c *Cropper) error {
		c.outSuffix = suffix
		return nil
	}
}

func WithOutDir(dir string) CropperOption {
	return func(c *Cropper) error {
		c.outDir = dir
		return nil
	}
}

func WithSkipUnchanged(skip bool) CropperOption {
	return func(c *Cropper) error {
		c.skipUnchanged = skip
		return nil
	}
}
