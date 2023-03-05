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

type DirectoryLoaderOption func(*DirectoryLoader) error

func WithRegex(regex string) DirectoryLoaderOption {
	return func(d *DirectoryLoader) error {
		re, err := regexp.Compile(regex)
		if err != nil {
			return err
		}

		d.regex = re

		return nil
	}
}

func WithRecursive(enable bool) DirectoryLoaderOption {
	return func(d *DirectoryLoader) error {
		d.recursive = enable
		return nil
	}
}

type DirectoryLoader struct {
	regex     *regexp.Regexp
	recursive bool
}

func NewDirectoryLoader(options ...DirectoryLoaderOption) (*DirectoryLoader, error) {
	d := &DirectoryLoader{}

	for _, opt := range options {
		if err := opt(d); err != nil {
			return nil, err
		}
	}

	return d, nil
}

func (d *DirectoryLoader) LoadCroppablesIter(dirs []string) (*CroppableLoadIterator, error) {
	var loader func(string) ([]string, error)

	if d.recursive {
		loader = d.loadRecursive
	} else {
		loader = d.loadDir
	}

	paths := []string{}
	errs := []error{}

	for _, dir := range dirs {
		path, err := loader(dir)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		paths = append(paths, path...)
	}

	return newCroppableLoadIterator(paths), errors.Join(errs...)
}

func (d *DirectoryLoader) LoadCroppables(dirs []string) ([]*Croppable, error) {
	var loader func(string) ([]string, error)

	if d.recursive {
		loader = d.loadRecursive
	} else {
		loader = d.loadDir
	}

	croppables := []*Croppable{}
	errs := []error{}

	for _, dir := range dirs {
		paths, err := loader(dir)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		for _, path := range paths {
			croppable, err := LoadCroppable(path)
			if err != nil {
				fmt.Println(err)
				continue
			}

			croppables = append(croppables, croppable)
		}
	}

	return croppables, errors.Join(errs...)
}

func (i *DirectoryLoader) loadRecursive(dir string) ([]string, error) {
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

func (i *DirectoryLoader) loadDir(dir string) ([]string, error) {
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
