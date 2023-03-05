package main

import (
	"errors"
	"fmt"
	"gocrop"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var imageFlags = []cli.Flag{
	&cli.Int64Flag{
		Name:  "threshold",
		Value: 0,
		Usage: "Sets the alpha threshold for cropping, default is 0. Alpha value is an integer in range of 0-255",
	}, &cli.Int64Flag{
		Name:  "padding",
		Value: 0,
		Usage: "Sets the number of transparent pixels that will surround the min cropped rectangle",
	},
	&cli.StringFlag{
		Name:  "out",
		Usage: "Sets the output directory for cropped images.",
		Value: "",
		Action: func(ctx *cli.Context, s string) error {
			if s == "" {
				return errors.New("out_dir value cannot be empty")
			}

			return nil
		},
	},
	&cli.StringFlag{
		Name:  "suffix",
		Usage: "Sets suffix of the cropped image. The suffix is placed before file extension: filename[suffix].jpg",
		Value: "",
		Action: func(ctx *cli.Context, s string) error {
			if s == "" {
				return errors.New("suffix cannot be empty")
			}

			return nil
		},
	},
	&cli.StringFlag{
		Name:  "prefix",
		Usage: "Sets suffix of the cropped image. The prefix is placed before file name: [prefix]filename.jpg",
		Value: "",
		Action: func(ctx *cli.Context, s string) error {
			if s == "" {
				return errors.New("prefix cannot be empty")
			}

			return nil
		},
	},
}

var directoryFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:  "recursive",
		Usage: "Enables recursive mode, cropping will be attempted in all subdirectories",
	}, &cli.StringFlag{
		Name:  "regex",
		Usage: "Sets regex for filtering directories",
		Action: func(ctx *cli.Context, s string) error {
			if s == "" {
				return errors.New("regex cannot be empty")
			}

			return nil
		},
	},
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "image",
				Aliases: []string{"img", "i"},
				Usage:   "crop selected images",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() == 0 {
						return errors.New("no images specified")
					}

					cropper, err := gocrop.NewCropper(
						gocrop.WithThreshold(uint32(cCtx.Int64("threshold"))),
						gocrop.WithPadding(cCtx.Int("padding")),
						gocrop.WithOutDir(cCtx.String("out_dir")),
						gocrop.WithOutPrefix(cCtx.String("prefix")),
						gocrop.WithOutSuffix(cCtx.String("suffix")),
					)
					if err != nil {
						return err
					}

					croppables := []*gocrop.Croppable{}

					for _, arg := range cCtx.Args().Slice() {
						croppable, err := gocrop.LoadCroppable(arg)
						if err != nil {
							fmt.Println(err)
							continue
						}

						croppables = append(croppables, croppable)
					}

					return cropper.Crop(croppables)
				},
				Flags: imageFlags,
			},
			{
				Name:    "directory",
				Aliases: []string{"dir", "d"},
				Usage:   "crop images in a directory",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() == 0 {
						return errors.New("no directories specified")
					}

					cropper, err := gocrop.NewCropper(
						gocrop.WithThreshold(uint32(cCtx.Int64("threshold"))),
						gocrop.WithPadding(cCtx.Int("padding")),
						gocrop.WithOutDir(cCtx.String("out_dir")),
						gocrop.WithOutPrefix(cCtx.String("prefix")),
						gocrop.WithOutSuffix(cCtx.String("suffix")),
					)
					if err != nil {
						return err
					}

					opts := []gocrop.DirectoryLoaderOption{
						gocrop.WithRecursive(cCtx.Bool("recursive")),
					}

					if cCtx.IsSet("regex") {
						opts = append(opts, gocrop.WithRegex(cCtx.String("regex")))
					}

					loader, err := gocrop.NewDirectoryLoader(opts...)
					if err != nil {
						return err
					}

					croppables, err := loader.LoadCroppablesIter(cCtx.Args().Slice())
					if err != nil {
						return err
					}

					return cropper.CropIter(croppables)
				},
				Flags: append(imageFlags, directoryFlags...),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
