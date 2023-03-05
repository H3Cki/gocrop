package gocrop_test

import (
	"gocrop"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadCroppable(t *testing.T) {
	tests := []struct {
		fn    string
		exErr error
	}{
		{"2squares.png", nil},
		{"blank.png", nil},
		{"circle.png", nil},
		{"line.gif", nil},
		{"white.jpg", gocrop.ErrUnsupportedFormat},
		{"rect.png", nil},
	}

	for _, tt := range tests {
		t.Run(tt.fn, func(t *testing.T) {
			croppable, err := gocrop.LoadCroppable(path.Join("testdata", tt.fn))
			assert.ErrorIs(t, tt.exErr, err)
			assert.Equal(t, croppable != nil, tt.exErr == nil)
		})
	}
}

// func TestDirNameExt(t *testing.T) {
// 	tests := []struct {
// 		fp   string
// 		dir  string
// 		name string
// 		ext  string
// 	}{
// 		{"test.png", ".", "test", "png"},
// 		{"test.png.jpg", ".", "test.png", "jpg"},
// 		{"test", ".", "test", ""},
// 		{"foo/bar/test", "foo\\bar", "test", ""},
// 		{"./foo/bar/test", "foo\\bar", "test", ""},
// 		{"C:\\foo\\test", "C:\\foo", "test", ""},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.fp, func(t *testing.T) {
// 			d, n, e := dirFileExt(tt.fp)
// 			assert.Equal(t, tt.dir, d)
// 			assert.Equal(t, tt.name, n)
// 			assert.Equal(t, tt.ext, e)
// 		})
// 	}
// }

// func TestFileExt(t *testing.T) {
// 	tests := []struct {
// 		fn   string
// 		name string
// 		ext  string
// 	}{
// 		{"test.png", "test", "png"},
// 		{"test.png.jpg", "test.png", "jpg"},
// 		{"test", "test", ""},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.fn, func(t *testing.T) {
// 			n, e := fileExt(tt.fn)
// 			assert.Equal(t, tt.name, n)
// 			assert.Equal(t, tt.ext, e)
// 		})
// 	}
// }
