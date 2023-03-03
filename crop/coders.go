package crop

import (
	"fmt"
	"image"
)

func cropImage(img image.Image) (image.Image, error) {
	rect := croppingRect(img)

	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	// img is an Image interface. This checks if the underlying value has a
	// method called SubImage. If it does, then we can use SubImage to crop the
	// image.
	simg, ok := img.(subImager)
	if !ok {
		return nil, fmt.Errorf("image does not support cropping")
	}

	return simg.SubImage(rect), nil
}

func croppingRect(img image.Image) image.Rectangle {
	rect := img.Bounds()

	cropRect := image.Rectangle{
		Min: rect.Max,
		Max: rect.Min,
	}

	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			pixel := img.At(x, y)

			_, _, _, alpha := pixel.RGBA()

			if alpha > 0 {
				if x < cropRect.Min.X {
					cropRect.Min.X = x
				}
				if y < cropRect.Min.Y {
					cropRect.Min.Y = y
				}
				if x > cropRect.Max.X {
					cropRect.Max.X = x
				}
				if y > cropRect.Max.Y {
					cropRect.Max.Y = y
				}
			}
		}
	}

	return cropRect
}
