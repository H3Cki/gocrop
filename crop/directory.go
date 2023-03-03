package crop

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
)

type Directory struct {
	Dirs      []string
	Prefix    string
	Suffix    string
	OutDir    string
	Regex     *regexp.Regexp
	Recursive bool
}

func (i *Directory) Crop() error {
	files := []string{}

	for _, dir := range i.Dirs {
		files = append(files, i.getDirFiles(dir)...)
	}

	imgCrop := Image{
		Paths:  files,
		Prefix: i.Prefix,
		Suffix: i.Suffix,
		OutDir: i.OutDir,
	}

	return imgCrop.Crop()
}

func (i *Directory) getDirFiles(dir string) []string {
	files := []string{}
	if i.Recursive {
		err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}

			fName := d.Name()
			if i.Regex == nil || (i.Regex != nil && i.Regex.Match([]byte(fName))) {
				files = append(files, p)
			}

			return nil
		})

		if err != nil {
			fmt.Printf("unable to read %s: %s\n", dir, err.Error())
			return []string{}
		}
	} else {
		fileInfos, err := ioutil.ReadDir(dir)
		if err != nil {
			fmt.Printf("unable to read %s: %s\n", dir, err.Error())
			return []string{}
		}

		for _, fi := range fileInfos {
			if fi.IsDir() {
				continue
			}

			fName := fi.Name()
			if i.Regex == nil || (i.Regex != nil && i.Regex.Match([]byte(fName))) {
				files = append(files, path.Join(dir, fName))
			}
		}
	}

	return files
}
