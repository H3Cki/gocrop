package gocrop_test

import (
	"gocrop"
	"image"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCropper_Rect loads all images from testdata/described
// parses their file names and compares their cropping rectangles with parsed points
// files in testdata/described should be named <anything>-<minX>-<minY>-<maxX>-<maxY>.png
func TestCropper_Rect(t *testing.T) {
	basicCropper, _ := gocrop.NewCropper()

	tests := []struct {
		cropper *gocrop.Cropper
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
			croppable, err := gocrop.LoadCroppable(path.Join("testdata/described", tt.fn))
			assert.NoError(t, err)
			assert.NotNil(t, croppable)

			croppingRect := tt.cropper.Rect(croppable.Cropper)
			assert.True(t, tt.exRect.Eq(croppingRect))
		})
	}
}

// TODO test only checks img size
func TestCropper_Crop(t *testing.T) {
	basicCropper, _ := gocrop.NewCropper()

	tests := []struct {
		cropper *gocrop.Cropper
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
			croppable, err := gocrop.LoadCroppable(path.Join("testdata/described", tt.fn))
			assert.NoError(t, err)
			assert.NotNil(t, croppable)

			cropped, ok := tt.cropper.Crop(croppable)
			assert.Equal(t, tt.ok, ok)

			expectedCrop, err := gocrop.LoadCroppable(path.Join("testdata/described/cropped", tt.fn))
			assert.NoError(t, err)

			assert.True(t, imagesEqual(expectedCrop.Cropper, cropped))
		})
	}
}

func imagesEqual(img1, img2 image.Image) bool {
	return img1.Bounds().Size().Eq(img2.Bounds().Size())
}
