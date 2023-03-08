# gocrop

gocrop provides a CLI and an API for cropping transparent images.

# Installation

```
go get github.com/H3Cki/gocrop
```

# CLI Examples

## 1. Crop specific images and output them with `_cropped` suffix:

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

## 2. Crop all images which file name matches a regex, in specific directories and all their subdirectories, output results into `images/cropped`:

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

## 1. Cropping single image

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

	"github.com/H3Cki/gocrop"
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


## 2. Cropping single image with 10-pixel padding and saving it in another directory

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

	"github.com/H3Cki/gocrop"
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


## 3. Cropping all images in all subdirectories and saving them with _cropped suffix in their original directory

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

	"github.com/H3Cki/gocrop"
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