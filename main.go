package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Sqrl",
		Usage: "Expand or crop images into a square",
		Flags: []cli.Flag{
			&cli.UintFlag{
				Name:    "size",
				Aliases: []string{"s"},
				Value:   1000,
				Usage:   "pixel size of square",
			},
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Value:   "",
				Usage:   "specify a filename to prevent resizing the whole folder",
			},
			&cli.BoolFlag{
				Name:    "expand",
				Aliases: []string{"e"},
				Value:   false,
				Usage:   "square by adding whitespace",
			},
		},
		Action: func(c *cli.Context) error {
			if c.String("file") == "" {
				// resize the folder
				imgs := getImages()
				for _, i := range imgs {
					if c.Bool("expand") {
						expandImg(i, c.Uint("size"))
					} else {
						cropImg(i, c.Uint("size"))
					}
				}
			} else {
				// specific file
				if c.Bool("expand") {
					expandImg(c.String("file"), c.Uint("size"))
				} else {
					cropImg(c.String("file"), c.Uint("size"))
				}
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func folderize() {
	imgs := getImages()
	for _, i := range imgs {
		expandImg(i, 1000)
		cropImg(i, 1000)
	}
}

func getImages() []string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	var files []string

	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		found, _ := regexp.MatchString("^*(.jpg|.png|.jpeg|.gif)", path)
		ignore, _ := regexp.MatchString("^*(sqrl-)", path)
		if found && !ignore {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func cropImg(path string, size uint) {
	fimg, _ := os.Open(path)
	defer fimg.Close()
	img, _, _ := image.Decode(fimg)
	bounds := img.Bounds()
	imgW := bounds.Max.X
	imgH := bounds.Max.Y
	max := 0
	if imgW > imgH {
		max = imgH
	} else {
		max = imgW
	}

	croppedImg, err := cutter.Crop(img, cutter.Config{
		Width:   max,
		Height:  max,
		Mode:    cutter.Centered,
		Options: cutter.Ratio & cutter.Copy, // Copy is useless here
	})
	if err != nil {
		log.Fatal(err)
	}

	resizedImg := resize.Resize(size, 0, croppedImg, resize.Lanczos3)

	file := filepath.Base(path)
	f, _ := os.Create("sqrl-crop-" + file)
	defer f.Close()

	opt := jpeg.Options{
		Quality: 90,
	}
	err = jpeg.Encode(f, resizedImg, &opt)
	if err != nil {
		log.Fatal(err)
	}
}

func expandImg(path string, size uint) {
	// expands the image to make a square
	// not to be confused with enlarging
	fimg, _ := os.Open(path)
	defer fimg.Close()
	img, _, _ := image.Decode(fimg)
	bounds := img.Bounds()
	imgW := bounds.Max.X
	imgH := bounds.Max.Y
	max := 0
	if imgW > imgH {
		max = imgW
	} else {
		max = imgH
	}

	// create a new white background
	upLeft := image.Point{0, 0}
	lowRight := image.Point{max, max}
	offsetX := (max - imgW) / 2
	offsetY := (max - imgH) / 2

	bg := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	for x := 0; x < max; x++ {
		for y := 0; y < max; y++ {
			if x > offsetX && x < max-offsetX && y > offsetY && y < max-offsetY {
				color := img.At(x-offsetX, y-offsetY)
				bg.Set(x, y, color)
			} else {
				bg.Set(x, y, color.White)
			}
		}
	}
	resizedImg := resize.Resize(size, 0, bg, resize.NearestNeighbor)
	opt := jpeg.Options{
		Quality: 90,
	}
	file := filepath.Base(path)
	f, _ := os.Create("sqrl-expand-" + file)
	jpeg.Encode(f, resizedImg, &opt)
}
