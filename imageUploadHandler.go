package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/panda/helper"
	"gopkg.in/mgo.v2/bson"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// save a single image
func handleSaveSingleImage(part *multipart.Part) (info ImageInfo, err error) {
	newID := bson.NewObjectId()
	date := time.Now().Format("20060102")

	err = helper.CreateDirIfNotExists(filepath.Join(config.SaveDir, date))
	if err != nil {
		return
	}
	path := filepath.Join(date, newID.Hex())
	savePath := filepath.Join(config.SaveDir, path)

	dst, err := os.Create(savePath)

	if err != nil {
		return
	}

	defer dst.Close()

	var bytes int64
	if bytes, err = io.Copy(dst, part); err != nil {
		return
	}

	width, height := helper.GetImageDimensions(savePath)

	URL := config.BaseURL + path

	var hash helper.HashInfo

	hash, err = helper.CalculateBasicHashes(savePath)

	if err != nil {
		return
	}

	info = ImageInfo{
		ID:        newID,
		Name:      part.FileName(),
		Extension: filepath.Ext(part.FileName()),
		BaseDir:   config.SaveDir,
		Path:      path,
		Width:     width,
		Height:    height,
		URL:       URL,
		Resizes:   map[string]string{},
		Hash:      hash,
		Size:      bytes,
		CreatedAt: time.Now(),
	}
	err = storeImage(&info)
	if err != nil {
		return
	}
	return info, nil
}

// upload multiple images
func handleMultipleImagesUpload(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if req.ContentLength > config.MaxSize {
		http.Error(res, "file too large", http.StatusRequestEntityTooLarge)
		return
	}

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
	helper.WriteResponse(res, imgs)
}

// upload single image
func handleSingleImageUpload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.ContentLength > config.MaxSize {
		http.Error(w, "file too large", http.StatusRequestEntityTooLarge)
		return
	}
	reader, err := r.MultipartReader()
	if err != nil {
		helper.WriteErrorResponse(w, err)
		return
	}
	part, err := reader.NextPart()
	if err != nil {
		helper.WriteErrorResponse(w, err)
		return
	}
	img, err := handleSaveSingleImage(part)
	if err != nil {
		helper.WriteErrorResponse(w, err)
		return
	}
	helper.WriteResponse(w, img)
}
