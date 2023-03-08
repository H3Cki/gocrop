# gocrop

gocrop provides a CLI and an API for cropping transparent images.

# Examples

Crop specific images using CLI and output them with `_cropped` suffix:
```cli
gocrop image --suffix _cropped img1.png img2.png
```

Crop all images which file name matches a regex, in specific directories and all their subdirectories, output all images into `images/cropped` directory:
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
    
}
```