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