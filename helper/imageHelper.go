package helper

import (
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

func GetImageDimensions(filepath string) (width int, height int) {
	if file, err := os.Open(filepath); err == nil {
		defer file.Close()
		cfg, _, err := image.DecodeConfig(file)
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
	if isPNG(extension) {
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
	if isPNG(extension) {
		png.Encode(out, m)
	} else {
		jpeg.Encode(out, m, nil)
	}
	return nil
}

func isPNG(extension string) bool {
	return strings.ToLower(extension) == ".png"
}
