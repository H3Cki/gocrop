package main

import (
	"errors"
	"fmt"
	"gocrop/crop"
	"log"
	"os"
	"regexp"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "image",
				Aliases: []string{"img", "i"},
				Usage:   "load image",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() == 0 {
						return errors.New("no images specified")
					}

					imgCrop := crop.Image{
						Paths:  cCtx.Args().Slice(),
						Suffix: cCtx.String("suffix"),
						Prefix: cCtx.String("prefix"),
						OutDir: cCtx.String("out_dir"),
					}

					return imgCrop.Crop()
				},
				Flags: []cli.Flag{
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
				},
			},
			{
				Name:    "directory",
				Aliases: []string{"dir", "d"},
				Usage:   "crop images in directory",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() == 0 {
						return errors.New("no directories specified")
					}

					dirs := []string{}
					for _, fp := range cCtx.Args().Slice() {
						if _, err := os.Stat(fp); err != nil {
							fmt.Printf("%s is not a valid directory\n", fp)
							continue
						}

						dirs = append(dirs, fp)
					}

					dirCrop := crop.Directory{
						Dirs:      dirs,
						Suffix:    cCtx.String("suffix"),
						Prefix:    cCtx.String("prefix"),
						OutDir:    cCtx.String("out_dir"),
						Recursive: cCtx.Bool("recursive"),
					}

					if cCtx.IsSet("regex") {
						re, err := regexp.Compile(cCtx.String("regex"))
						if err != nil {
							return fmt.Errorf("invalid regex: %w", err)
						}

						dirCrop.Regex = re
					}

					return dirCrop.Crop()
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "out_dir",
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
					}, &cli.BoolFlag{
						Name:  "recursive",
						Usage: "Enables recursive mode, cropping will be attempted in all subdirectories",
					}, &cli.StringFlag{
						Name:  "regex",
						Usage: "Sets regex for filtering directories",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
