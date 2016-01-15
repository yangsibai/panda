package main

import (
	"errors"
	"fmt"
	"github.com/twinj/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type ImageInfo struct {
	Name   string
	URL    string
	Width  int
	Height int
}

func handleSaveSingleImage(part *multipart.Part) (info ImageInfo, err error) {
	if part.FileName() == "" {
		err = errors.New("filename is empty")
		return
	}
	savePath := getNewSaveName(part.FileName())
	dst, err := os.Create(savePath)

	defer dst.Close()

	if err != nil {
		return
	}

	if _, err = io.Copy(dst, part); err != nil {
		return
	}

	width, height := getImageDimensions(savePath)

	info = ImageInfo{
		Name:   part.FileName(),
		URL:    savePath,
		Width:  width,
		Height: height,
	}
	return info, nil
}

func handleImageUpload(res http.ResponseWriter, req *http.Request) {
	reader, err := req.MultipartReader()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		info, err := handleSaveSingleImage(part)
		fmt.Println(info.Name, info.Width, info.Height)
	}
	io.WriteString(res, "ok")
	res.WriteHeader(200)
}

func getNewSaveName(filename string) string {
	_, config := readConfig()
	newName := uuid.NewV4().String() + filepath.Ext(filename)
	return filepath.Join(config.SaveDir, newName)
}
