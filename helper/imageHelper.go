package helper

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func GetImageDimensions(filepath string) (width int, height int) {
	if reader, err := os.Open(filepath); err == nil {
		defer reader.Close()
		im, _, err := image.DecodeConfig(reader)
		if err == nil {
			width = im.Width
			height = im.Height
			return
		}
	}
	return -1, -1
}
