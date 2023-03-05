package gocrop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadCroppable(t *testing.T) {
	tests := []struct {
		path        string
		expectedErr error
	}{
		{"testdata/t.png", nil},
		{"testdata/e.png", nil},
		{"testdata/testdir1/testimg1.png", ErrImageLoadFailed},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			cpb, err := LoadCroppable(tt.path)
			assert.ErrorIs(t, tt.expectedErr, err)
			assert.True(t, (cpb != nil) == (tt.expectedErr == nil))
		})
	}
}

func TestDirNameExt(t *testing.T) {
	tests := []struct {
		fp   string
		dir  string
		name string
		ext  string
	}{
		{"test.png", ".", "test", "png"},
		{"test.png.jpg", ".", "test.png", "jpg"},
		{"test", ".", "test", ""},
		{"foo/bar/test", "foo\\bar", "test", ""},
		{"./foo/bar/test", "foo\\bar", "test", ""},
		{"C:\\foo\\test", "C:\\foo", "test", ""},
	}

	for _, tt := range tests {
		t.Run(tt.fp, func(t *testing.T) {
			d, n, e := dirFileExt(tt.fp)
			assert.Equal(t, tt.dir, d)
			assert.Equal(t, tt.name, n)
			assert.Equal(t, tt.ext, e)
		})
	}
}

func TestFileExt(t *testing.T) {
	tests := []struct {
		fn   string
		name string
		ext  string
	}{
		{"test.png", "test", "png"},
		{"test.png.jpg", "test.png", "jpg"},
		{"test", "test", ""},
	}

	for _, tt := range tests {
		t.Run(tt.fn, func(t *testing.T) {
			n, e := fileExt(tt.fn)
			assert.Equal(t, tt.name, n)
			assert.Equal(t, tt.ext, e)
		})
	}
}
