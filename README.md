# gocrop

gocrop provides a CLI and an API for cropping transparent images.

[![GoDoc](https://godoc.org/github.com/H3Cki/gocrop?status.svg)](https://godoc.org/github.com/H3Cki/gocrop)
[![Go Report Card](https://goreportcard.com/badge/github.com/H3Cki/gocrop)](https://goreportcard.com/report/github.com/H3Cki/gocrop)
[![codecov](https://codecov.io/gh/H3Cki/gocrop/branch/master/graph/badge.svg)](https://codecov.io/gh/H3Cki/gocrop)
[![build](https://github.com/H3Cki/gocrop/actions/workflows/build.yaml/badge.svg?branch=master)](https://github.com/H3Cki/gocrop/actions/workflows/build.yaml)

# Installation
In order to use the API:
```
go get github.com/H3Cki/gocrop
```

In order to install the CLI:
```
go install github.com/H3Cki/gocrop@latest
```

# How it works
All transparent pixels (or at least those with alpha value higher than the provided threshold) are discarded in all directions, left, top, right, bottom. If padding option was used the cropped image will be extended by the provided padding amount, equally in all directions. Examples (with hacky border to visualize the image size better, I suggest opening the image anyways):

Original image:

<kbd>
<img src="https://hecki.codes/gocrop/circle.png">
</kbd>

Cropped image:

<kbd>
<img src="https://hecki.codes/gocrop/circle_cropped.png">
</kbd>

Cropped and padded image:

<kbd>
<img src="https://hecki.codes/gocrop/circle_cropped_padded.png">
</kbd>

# CLI Examples

### 1. Crop specific images and output them with `_cropped` suffix:

Directory tree before:
```   
├─ img1.png
├─ img2.png
├─ img3.png

```

Directory tree after:
```       
├─ img1.png
├─ img1_cropped.png
├─ img2.png
├─ img2_cropped.png
├─ img3.png

```

```cli
gocrop image --suffix _cropped img1.png img2.png
```

### 2. Crop all images which file name matches a regex, in specific directories and all their subdirectories, output results into `images/cropped`:

Directory tree before:
```
dir1/
├─ img1.png
├─ img2.gif
dir2/
├─ img3.gif

```

Directory tree after:
```
dir1/
├─ img1.png
├─ img2.gif
dir2/
├─ img3.gif  
images/
├─ cropped/
│  ├─ img2.gif
│  ├─ img3.gif

```

```cli
gocrop directory --out_dir images/cropped --regex ^.*gif.*$ --recursive dir1 dir2
```

# API Examples

### 1. Cropping single image

Directory tree before:
```
images/         
├─ test.png

```

Directory tree after:
```
images/         
├─ test.png // image is overwritten

```


```go
import (
	"fmt"

	"github.com/H3Cki/gocrop/gocrop"
)

func main() {
    // Load image, doing it this way assures the image is croppable 
    // and it's format is supported
	croppable, err := gocrop.Load("images/test.png")
	if err != nil {
		fmt.Println(err)
		return
	}

    // Create cropper with no options (source image will be overwritten)
	cropper, _ := gocrop.NewCropper()
	if err != nil {
		fmt.Println(err)
		return
	}

    // Crop the image
	cropped, ok := cropper.Crop(croppable)
    // If ok is false we can skip saving the image because no changes were made
	if !ok {
		fmt.Println("cropping would make no difference to target image")
		return
	}

    // Encode and save the image at "images/test.png"
	err = cropper.Save(cropped)
	if err != nil {
		fmt.Println(err)
		return
	}
}
```


### 2. Cropping single image with 10-pixel padding and saving it in another directory

Directory tree before:
```
images/         
├─ test.png

```

Directory tree after:
```
images/
├─ test.png
cropped/
├─ test.png

```


```go
package main

import (
	"fmt"

	"github.com/H3Cki/gocrop/gocrop"
)

func main() {
    // Load image, doing it this way assures the image is croppable 
    // and it's format is supported
	croppable, err := gocrop.Load("images/test.png")
	if err != nil {
		fmt.Println(err)
		return
	}

    // Create a cropper that saves the images in "images/cropped" directory (will be created if doesn't exist) and
    // applies a 10px padding to the cropped image.
	cropper, err := gocrop.NewCropper(
		gocrop.WithOutDir("images/cropped"),
		gocrop.WithPadding(10),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

    // This time we used CropAndSave method for convenience.
	err = cropper.CropAndSave(croppable)
	if err != nil {
		fmt.Println(err)
		return
	}
}
```


### 3. Cropping all images in all subdirectories and saving them with _cropped suffix in their original directory

Directory tree before:
```
images/
├─ image1.png
├─ avatars/
│  ├─ avatar1.png
│  ├─ avatar2.png
│  ├─ icons/
│  │  ├─ icon.png
```

Directory tree after:
```
images/
├─ image1.png
├─ image1_cropped.png
├─ avatars/
│  ├─ avatar1.png
│  ├─ avatar1_cropped.png
│  ├─ avatar2.png
│  ├─ avatar2_cropped.png
│  ├─ icons/
│  │  ├─ icon.png
│  │  ├─ icon_cropped.png
```


```go
package main

import (
	"fmt"

	"github.com/H3Cki/gocrop/gocrop"
)

func main() {
	// Create a finder with recursive option to traverse all subdirectories
	finder, err := gocrop.NewFinder(gocrop.WithRecursive(true))
	if err != nil {
		fmt.Println(err)
		return
	}

	directories := []string{"."}

	// Find all images in supported formats and wrap them in Croppable
	croppables, err := finder.Find(directories)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a cropper with OutSuffix option that will append "_cropped" to the end of cropped images when saving them
	cropper, err := gocrop.NewCropper(gocrop.WithOutSuffix("_cropped"))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Iterate over found croppables
	for _, croppable := range croppables {
		// Load the image (loads from croppable.Path and decodes it using croppable.Decode)
		if err := croppable.Load(); err != nil {
			continue
		}

		// Crop and save it
		if err := cropper.CropAndSave(croppable); err != nil {
			fmt.Printf("error cropping %s: %s\n", croppable.Path, err.Error())
		}
	}
}
```
