package gocropper

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

type pathTest struct {
	fp   string
	dir  string
	name string
	ext  string
}

func TestDirNameExt(t *testing.T) {
	var tests []pathTest

	if runtime.GOOS == "windows" {
		tests = []pathTest{
			{"test.png", ".", "test", ".png"},
			{"test.png.jpg", ".", "test.png", ".jpg"},
			{"test", ".", "test", ""},
			{"foo\\bar\\test", "foo\\bar", "test", ""},
			{".\\foo\\bar\\test", "foo\\bar", "test", ""},
			{"C:\\foo\\test", "C:\\foo", "test", ""},
		}
	} else {
		tests = []pathTest{
			{"test.png", ".", "test", ".png"},
			{"test.png.jpg", ".", "test.png", ".jpg"},
			{"test", ".", "test", ""},
			{"foo/bar/test", "foo/bar", "test", ""},
			{"./foo/bar/test", "foo/bar", "test", ""},
			{"C:/foo/test", "C:/foo", "test", ""},
		}
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
		{"test.png", "test", ".png"},
		{"test.png.jpg", "test.png", ".jpg"},
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
