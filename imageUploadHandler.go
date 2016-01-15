package main

import (
	"encoding/json"
	"github.com/twinj/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"panda/helper"
	"path/filepath"
)

type ImageInfo struct {
	Name   string `json:"name"`
	URL    string
	Width  int `json:"width"`
	Height int `json:"height"`
}

func handleSaveSingleImage(part *multipart.Part) (info ImageInfo, err error) {
	savePath, URL := getSaveNameAndURL(part.FileName())
	dst, err := os.Create(savePath)

	defer dst.Close()

	if err != nil {
		return
	}

	if _, err = io.Copy(dst, part); err != nil {
		return
	}

	width, height := helper.GetImageDimensions(savePath)

	info = ImageInfo{
		Name:   part.FileName(),
		URL:    URL,
		Width:  width,
		Height: height,
	}
	return info, nil
}

func handleImageUpload(res http.ResponseWriter, req *http.Request) {
	reader, err := req.MultipartReader()
	if err != nil {
		helper.WriteErrorResponse(res, err)
		return
	}
	var imgs []ImageInfo
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if part.FileName() == "" {
			continue
		}
		info, err := handleSaveSingleImage(part)
		imgs = append(imgs, info)
	}
	bytes, err := json.Marshal(imgs)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(bytes)
}

func getNewSaveName(filename string) string {
	_, config := readConfig()
	newName := uuid.NewV4().String() + filepath.Ext(filename)
	return filepath.Join(config.SaveDir, newName)
}

func getSaveNameAndURL(filename string) (savename string, URL string) {
	_, config := readConfig()
	newName := uuid.NewV4().String() + filepath.Ext(filename)
	savename = filepath.Join(config.SaveDir, newName)
	URL = filepath.Join(config.BaseURL, newName)
	return
}
