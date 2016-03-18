package routes

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	_ "github.com/nfnt/resize"
	"github.com/yangsibai/panda/db"
	"github.com/yangsibai/panda/helper"
	"github.com/yangsibai/panda/models"
	"gopkg.in/mgo.v2/bson"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// get single image by id and width
func HandleFetchSingleImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	widths := r.URL.Query()["w"]
	var width int
	var err error
	if len(widths) == 0 || widths[0] == "" {
		width = 0
	} else {
		width, err = strconv.Atoi(widths[0])
		if err != nil {
			width = 0
		}
	}

	session := db.GetSession()
	C := session.DB("resource").C("image")
	defer session.Close()

	info := models.ImageInfo{}

	oid := bson.ObjectIdHex(id)
	if err := C.FindId(oid).One(&info); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if width == 0 || width >= info.Width {
		http.Redirect(w, r, info.URL, 301)
		return
	}

	resize_key := fmt.Sprintf("w_%d", width)
	if val, ok := info.Resizes[resize_key]; ok {
		http.Redirect(w, r, val, 301)
		return
	}

	newPath := fmt.Sprintf(info.Path+"_w_%d", width)
	newAbsolutePath := filepath.Join(helper.Config.SaveDir, newPath)
	err = helper.CreateThumbnail(filepath.Join(helper.Config.SaveDir, info.Path), info.Extension, newAbsolutePath, uint(width))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if info.Resizes == nil {
		info.Resizes = map[string]string{}
	}
	info.Resizes[resize_key] = helper.Config.BaseURL + newPath
	change := bson.M{"$set": bson.M{"resizes": info.Resizes}}
	err = C.Update(bson.M{"_id": oid}, change)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, helper.Config.BaseURL+newPath, 301)
}

// save a single image
func handleSaveSingleImage(part *multipart.Part) (info models.ImageInfo, err error) {
	newID := bson.NewObjectId()
	date := time.Now().Format("20060102")

	err = helper.CreateDirIfNotExists(filepath.Join(helper.Config.SaveDir, date))
	if err != nil {
		return
	}
	path := filepath.Join(date, newID.Hex())
	savePath := filepath.Join(helper.Config.SaveDir, path)

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

	URL := helper.Config.BaseURL + path

	var hash models.HashInfo

	hash, err = helper.CalculateBasicHashes(savePath)

	if err != nil {
		return
	}

	info = models.ImageInfo{
		ID:        newID,
		Name:      part.FileName(),
		Extension: filepath.Ext(part.FileName()),
		BaseDir:   helper.Config.SaveDir,
		Path:      path,
		Width:     width,
		Height:    height,
		URL:       URL,
		Resizes:   map[string]string{},
		Hash:      hash,
		Size:      bytes,
		CreatedAt: time.Now(),
	}
	err = db.StoreImage(&info)
	if err != nil {
		return
	}
	return info, nil
}

// upload multiple images
func HandleMultipleImagesUpload(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if req.ContentLength > helper.Config.MaxSize {
		http.Error(res, "file too large", http.StatusRequestEntityTooLarge)
		return
	}

	reader, err := req.MultipartReader()
	if err != nil {
		helper.WriteErrorResponse(res, err)
		return
	}
	var imgs []models.ImageInfo
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
func HandleSingleImageUpload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.ContentLength > helper.Config.MaxSize {
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
