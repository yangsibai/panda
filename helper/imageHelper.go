package helper

import (
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

func GetImageDimensions(filepath string) (width int, height int) {
	if reader, err := os.Open(filepath); err == nil {
		defer reader.Close()
		cfg, _, err := image.DecodeConfig(reader)
		if err == nil {
			width = cfg.Width
			height = cfg.Height
			return
		}
	}
	return -1, -1
}

func CreateThumbnail(filepath, extension, outpath string, width uint) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	var img image.Image
	if extension == ".png" {
		img, err = png.Decode(file)
	} else {
		img, err = jpeg.Decode(file)
	}
	if err != nil {
		return err
	}
	file.Close()

	m := resize.Resize(width, 0, img, resize.NearestNeighbor)

	out, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer out.Close()
	if extension == ".png" {
		png.Encode(out, m)
	} else {
		jpeg.Encode(out, m, nil)
	}
	return nil
}
