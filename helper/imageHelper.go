package helper

import (
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"os"
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
	isPNG := getFormat(file) == "png"

	if isPNG {
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
	if isPNG {
		png.Encode(out, m)
	} else {
		jpeg.Encode(out, m, nil)
	}
	return nil
}

func getFormat(file *os.File) string {
	bytes := make([]byte, 4)
	n, _ := file.ReadAt(bytes, 0)
	if n < 4 {
		return ""
	}
	if bytes[0] == 0x89 && bytes[1] == 0x50 && bytes[2] == 0x4E && bytes[3] == 0x47 {
		return "png"
	}
	if bytes[0] == 0xFF && bytes[1] == 0xD8 {
		return "jpg"
	}
	if bytes[0] == 0x47 && bytes[1] == 0x49 && bytes[2] == 0x46 && bytes[3] == 0x38 {
		return "gif"
	}
	if bytes[0] == 0x42 && bytes[1] == 0x4D {
		return "bmp"
	}
	return ""
}
