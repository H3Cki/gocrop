package gocrop

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

// CroppableFinder is a tool for loading all valid images from given directories.
type CroppableFinder struct {
	regex     *regexp.Regexp
	recursive bool
}

// NewCroppableFinder is a constructor for *DirectoryLoader, accpts DirectoryLoaderOptions,
// returns an error when any option fails.
func NewCroppableFinder(options ...CroppableFinderOptions) (*CroppableFinder, error) {
	d := &CroppableFinder{}

	for _, opt := range options {
		if err := opt(d); err != nil {
			return nil, err
		}
	}

	return d, nil
}

// LoadCroppablesIter takes a slice of directory paths and returns an iterator for all valid images
// found in those directories. Can go through all subdirectories with WithRecursive(bool) option set to true.
func (d *CroppableFinder) Find(dirs []string) ([]string, error) {
	var loader func(string) ([]string, error)

	if d.recursive {
		loader = d.findRecursive
	} else {
		loader = d.findInDir
	}

	paths := []string{}
	errs := []error{}

	for _, dir := range dirs {
		p, err := loader(dir)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		paths = append(paths, p...)
	}

	return paths, errors.Join(errs...)
}

func (i *CroppableFinder) findRecursive(dir string) ([]string, error) {
	paths := []string{}

	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if i.regex == nil || (i.regex != nil && i.regex.MatchString(d.Name())) {
			paths = append(paths, p)
		}

		return nil
	})

	if err != nil {
		return []string{}, err
	}

	return paths, nil
}

func (i *CroppableFinder) findInDir(dir string) ([]string, error) {
	paths := []string{}

	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("unable to read %s: %s\n", dir, err.Error())
		return []string{}, err
	}

	for _, fi := range fileInfos {
		if fi.IsDir() {
			continue
		}

		fName := fi.Name()
		if i.regex == nil || (i.regex != nil && i.regex.MatchString(fName)) {
			paths = append(paths, path.Join(dir, fName))
		}
	}

	return paths, nil
}

type CroppableFinderOptions func(*CroppableFinder) error

// WithRegex compiles given regex string and sets *DirectoryLoader's
// regex to compiled result.
func WithRegex(regex string) CroppableFinderOptions {
	return func(d *CroppableFinder) error {
		re, err := regexp.Compile(regex)
		if err != nil {
			return err
		}

		d.regex = re

		return nil
	}
}

// WithRecursive sets the recursive flag of *DirectoryLoader
func WithRecursive(enable bool) CroppableFinderOptions {
	return func(d *CroppableFinder) error {
		d.recursive = enable
		return nil
	}
}
