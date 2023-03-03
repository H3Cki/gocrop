package crop

import (
	"fmt"
	"os"
	"path"
	"sync"
)

type Image struct {
	Paths  []string
	Prefix string
	Suffix string
	OutDir string
}

func (i *Image) Crop() error {
	if i.OutDir != "" {
		if err := os.MkdirAll(i.OutDir, os.ModePerm); err != nil {
			return err
		}
	}

	wg := &sync.WaitGroup{}

	wg.Add(len(i.Paths))

	for _, fp := range i.Paths {
		go i.handle(fp, wg)
	}

	wg.Wait()

	return nil
}

func (i *Image) handle(fp string, wg *sync.WaitGroup) {
	defer wg.Done()

	dir, name, ext := dirFileExt(fp)

	coder, ok := imgCoders[ext]
	if !ok {
		return
	}

	img, err := loadImage(fp, coder)
	if err != nil {
		fmt.Printf("error loading %s: %s\n", fp, err.Error())
		return
	}

	img, err = cropImage(img)
	if err != nil {
		fmt.Printf("error cropping %s: %s\n", fp, err.Error())
		return
	}

	if i.OutDir != "" {
		dir = i.OutDir
	}

	name = i.Prefix + name + i.Suffix + "." + ext
	outPath := path.Join(dir, name)

	err = saveImage(outPath, img, coder)
	if err != nil {
		fmt.Printf("error saving %s: %s\n", fp, err.Error())
		return
	}
}
