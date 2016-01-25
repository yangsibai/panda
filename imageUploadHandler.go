package main

import (
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"panda/helper"
	"path/filepath"
	"time"
)

type ImageInfo struct {
	ID        bson.ObjectId     `json:"id" bson:"_id"`
	Name      string            `json:"name" bson:"name"`
	BaseDir   string            `json:"-" bson: "baseDir"`
	Path      string            `json:"path" bson:"path" `
	Extension string            `json:"extension" bson:"extension"`
	Width     int               `json:"width" bson:"width"`
	Height    int               `json:"height" bson:"height"`
	URL       string            `json:"URL" bson:"URL"`
	Resizes   map[string]string `json:"resizes" bson:"resizes"`
}

func storeImage(info *ImageInfo) (err error) {
	session := getSession()
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	C := session.DB("resource").C("image")
	if err != nil {
		return
	}
	err = C.Insert(&info)
	return
}

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

	if _, err = io.Copy(dst, part); err != nil {
		return
	}

	width, height := helper.GetImageDimensions(savePath)

	URL := config.BaseURL + path
	info = ImageInfo{
		ID:        newID,
		Name:      part.FileName(),
		Extension: filepath.Ext(part.FileName()),
		BaseDir:   config.SaveDir,
		Path:      path,
		Width:     width,
		Height:    height,
		URL:       URL,
		Resizes: map[string]string{
			"w_0": URL,
		},
	}
	err = storeImage(&info)
	if err != nil {
		return
	}
	return info, nil
}

func handleImageUpload(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
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
