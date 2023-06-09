package gocropper_test

import (
	"image"
	"path"
	"sync"
	"testing"

	"github.com/H3Cki/gocrop/gocropper"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		fn    string
		exErr error
	}{
		{"2squares.png", nil},
		{"blank.png", nil},
		{"circle.png", nil},
		{"line.gif", nil},
		{"white.jpg", gocropper.ErrUnsupportedFormat},
		{"rect.png", nil},
	}

	for _, tt := range tests {
		t.Run(tt.fn, func(t *testing.T) {
			croppable, err := gocropper.Load(path.Join("testdata", tt.fn))
			assert.ErrorIs(t, tt.exErr, err)
			assert.Equal(t, croppable != nil, tt.exErr == nil)
		})
	}
}

func TestCropper_Rect(t *testing.T) {
	basicCropper, _ := gocropper.NewCropper()

	tests := []struct {
		cropper *gocropper.Cropper
		fn      string
		exRect  image.Rectangle
	}{
		{
			cropper: basicCropper,
			fn:      "white-0-0-100-100.png",
			exRect:  image.Rect(0, 0, 100, 100),
		},
		{
			cropper: basicCropper,
			fn:      "circle-25-25-75-75.png",
			exRect:  image.Rect(25, 25, 75, 75),
		},
		{
			cropper: basicCropper,
			fn:      "rect-25-30-75-70.png",
			exRect:  image.Rect(25, 30, 75, 70),
		},
		{
			cropper: basicCropper,
			fn:      "recthollow-25-30-75-70.png",
			exRect:  image.Rect(25, 30, 75, 70),
		},
		{
			cropper: basicCropper,
			fn:      "line1px-49-0-50-100.png",
			exRect:  image.Rect(49, 0, 50, 100),
		},
		{
			cropper: basicCropper,
			fn:      "line1px-49-0-50-100.gif",
			exRect:  image.Rect(49, 0, 50, 100),
		},
	}

	for _, tt := range tests {
		t.Run(tt.fn, func(t *testing.T) {
			croppable, err := gocropper.Load(path.Join("testdata/described", tt.fn))
			assert.NoError(t, err)
			assert.NotNil(t, croppable)

			croppingRect := tt.cropper.Rect(croppable.Image)
			assert.True(t, tt.exRect.Eq(croppingRect))
		})
	}
}

// TODO test only checks img size
func TestCropper_Crop(t *testing.T) {
	basicCropper, _ := gocropper.NewCropper()

	tests := []struct {
		cropper *gocropper.Cropper
		fn      string
		exFn    string
		ok      bool
	}{
		{
			cropper: basicCropper,
			fn:      "white-0-0-100-100.png",
			exFn:    "white-0-0-100-100.png",
			ok:      false,
		},
		{
			cropper: basicCropper,
			fn:      "circle-25-25-75-75.png",
			exFn:    "circle-25-25-75-75.png",
			ok:      true,
		},
		{
			cropper: basicCropper,
			fn:      "rect-25-30-75-70.png",
			exFn:    "rect-25-30-75-70.png",
			ok:      true,
		},
		{
			cropper: basicCropper,
			fn:      "recthollow-25-30-75-70.png",
			exFn:    "recthollow-25-30-75-70.png",
			ok:      true,
		},
		{
			cropper: basicCropper,
			fn:      "line1px-49-0-50-100.png",
			exFn:    "line1px-49-0-50-100.png",
			ok:      true,
		},
		{
			cropper: basicCropper,
			fn:      "line1px-49-0-50-100.gif",
			exFn:    "line1px-49-0-50-100.gif",
			ok:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.fn, func(t *testing.T) {
			croppable, err := gocropper.Load(path.Join("testdata/described", tt.fn))
			assert.NoError(t, err)
			assert.NotNil(t, croppable)

			cropped, ok := tt.cropper.Crop(croppable)
			assert.Equal(t, tt.ok, ok)

			if !ok {
				return
			}

			expectedCrop, err := gocropper.Load(path.Join("testdata/described/cropped", tt.fn))
			assert.NoError(t, err)

			assert.True(t, imagesEqual(expectedCrop.Image, cropped.Image))
		})
	}
}

func imagesEqual(img1, img2 image.Image) bool {
	return img1.Bounds().Size().Eq(img2.Bounds().Size())
}

func BenchmarkCropper_Crop(b *testing.B) {
	cropper, _ := gocropper.NewCropper(gocropper.WithOutDir("testoutput"))
	cpbl, err := gocropper.Load("testdata/bigblank.png")
	assert.NoError(b, err)

	n := 10

	b.Run("x", func(b *testing.B) {
		wg := &sync.WaitGroup{}

		for i := 0; i < n; i++ {
			wg.Add(1)

			go func(ii int) {
				defer wg.Done()
				cropper.Crop(cpbl)
			}(i)
		}

		wg.Wait()
	})

}
