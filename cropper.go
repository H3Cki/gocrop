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
	threshold uint32
	outPrefix string
	outSuffix string
	outDir    string
	padding   int
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

func (i *Cropper) CropAndSave(croppables []*Croppable) error {
	if i.outDir != "" {
		if err := os.MkdirAll(i.outDir, os.ModePerm); err != nil {
			return err
		}
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(croppables))

	for _, croppable := range croppables {
		go func(c *Croppable) {
			defer wg.Done()
			cropped, ok := i.crop(c)
			if !ok {
				return
			}

			err := i.save(c, cropped)
			if err != nil {
				fmt.Println(err)
			}
		}(croppable)
	}

	wg.Wait()

	return nil
}

func (i *Cropper) Crop(croppables []*Croppable) ([]image.Image, error) {
	if i.outDir != "" {
		if err := os.MkdirAll(i.outDir, os.ModePerm); err != nil {
			return []image.Image{}, err
		}
	}

	crops := []image.Image{}

	for _, croppable := range croppables {
		cropped, _ := i.crop(croppable)
		crops = append(crops, cropped)
	}

	return crops, nil
}

type CroppableIterator interface {
	Reset()
	Next()
	Current() (*Croppable, error)
	Valid() bool
}

func (i *Cropper) CropIter(iter CroppableIterator) error {
	if i.outDir != "" {
		if err := os.MkdirAll(i.outDir, os.ModePerm); err != nil {
			return err
		}
	}

	wg := &sync.WaitGroup{}

	for iter.Reset(); iter.Valid(); iter.Next() {
		croppable, err := iter.Current()
		if err != nil {
			fmt.Println(err)
			continue
		}

		wg.Add(1)
		go func() {
			i.crop(croppable)
		}()
	}

	wg.Wait()

	return nil
}

func (i *Cropper) crop(c *Croppable) (image.Image, bool) {
	rect := i.cropperRect(c.Cropper)

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
			return nil, false
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

func (i *Cropper) cropperRect(img image.Image) image.Rectangle {
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
	Cropper           CroppableImage
	Encode            func(w io.Writer, m image.Image) error
}

func (c *Croppable) Path() string {
	return fmt.Sprintf("%s/%s.%s", c.Dir, c.Name, c.Format)
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
