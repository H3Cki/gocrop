package gocropper

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

// Finder is a tool for loading all valid images from given directories.
// It can apply regex filtering and descend recursively through the directory tree.
type Finder struct {
	regex     *regexp.Regexp
	recursive bool
}

// NewFinder is a constructor for *Finder, accepts FinderOptions,
// returns an error if any option fails.
func NewFinder(options ...FinderOptions) (*Finder, error) {
	d := &Finder{}

	for _, opt := range options {
		if err := opt(d); err != nil {
			return nil, err
		}
	}

	return d, nil
}

// Find takes a slice of directory paths, searches those directories for images in supported formats
// and returns a list of croppables, ready to be loaded. Does not load images to check if they can be cropped.
func (d *Finder) Find(dirs []string) ([]*Croppable, error) {
	var loader func(string) ([]*Croppable, error)

	if d.recursive {
		loader = d.findRecursive
	} else {
		loader = d.findInDir
	}

	crops := []*Croppable{}
	errs := []error{}

	for _, dir := range dirs {
		p, err := loader(dir)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		crops = append(crops, p...)
	}

	return crops, errors.Join(errs...)
}

func (i *Finder) findRecursive(dir string) ([]*Croppable, error) {
	crops := []*Croppable{}

	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		_, _, ext := dirFileExt(p)
		coder, ok := imageCoders[ext]
		if !ok {
			return nil
		}

		if i.regex == nil || (i.regex != nil && i.regex.MatchString(d.Name())) {
			crops = append(crops, &Croppable{
				Path:   p,
				Image:  nil,
				Decode: coder.decode,
				Encode: coder.encode,
			})
		}

		return nil
	})

	if err != nil {
		return []*Croppable{}, err
	}

	return crops, nil
}

func (i *Finder) findInDir(dir string) ([]*Croppable, error) {
	crops := []*Croppable{}

	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("unable to read %s: %s\n", dir, err.Error())
		return []*Croppable{}, err
	}

	for _, fi := range fileInfos {
		if fi.IsDir() {
			continue
		}

		_, ext := fileExt(fi.Name())
		coder, ok := imageCoders[ext]
		if !ok {
			continue
		}

		if i.regex == nil || (i.regex != nil && i.regex.MatchString(fi.Name())) {
			crops = append(crops, &Croppable{
				Path:   path.Join(dir, fi.Name()),
				Image:  nil,
				Decode: coder.decode,
				Encode: coder.encode,
			})
		}
	}

	return crops, nil
}

type FinderOptions func(*Finder) error

// WithRegex attempts to compile given regex string
// if successful enables Finder regex filtering.
func WithRegex(regex string) FinderOptions {
	return func(d *Finder) error {
		re, err := regexp.Compile(regex)
		if err != nil {
			return err
		}

		d.regex = re

		return nil
	}
}

// WithRecursive enables Finder to traverse all subdirectories.
func WithRecursive(enable bool) FinderOptions {
	return func(d *Finder) error {
		d.recursive = enable
		return nil
	}
}
