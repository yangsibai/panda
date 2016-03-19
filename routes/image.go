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

func getImagePath(info models.ImageInfo, width int) (imgPath string, err error) {
	if width == 0 || width >= info.Width {
		imgPath = info.Path
	} else {
		resize_key := fmt.Sprintf("w_%d", width)
		imgPath = fmt.Sprintf(info.Path+"_w_%d", width)
		if _, exists := info.Resizes[resize_key]; exists == false {
			originalAbsolutePath := filepath.Join(helper.Config.SaveDir, info.Path)
			newAbsolutePath := filepath.Join(helper.Config.SaveDir, imgPath)
			err = helper.CreateThumbnail(originalAbsolutePath, info.Extension, newAbsolutePath, uint(width))
			if err == nil {
				info.Resizes[resize_key] = imgPath
				err = updateImageInfoResize(info.ID, info.Resizes)
			}
		}
	}
	return
}

func getImageInfoById(id string) (info models.ImageInfo, err error) {
	session := db.GetSession()
	C := session.DB("resource").C("image")
	defer session.Close()

	oid := bson.ObjectIdHex(id)
	err = C.FindId(oid).One(&info)
	return
}

func updateImageInfoResize(oid bson.ObjectId, resizes map[string]string) (err error) {
	session := db.GetSession()
	C := session.DB("resource").C("image")
	defer session.Close()

	err = C.Update(bson.M{"_id": oid}, bson.M{"$set": bson.M{"resizes": resizes}})
	return
}

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

	info, err := getImageInfoById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	imgPath, err := getImagePath(info, width)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if helper.Config.ResourceServerBaseURL == "" {
		f, err := os.Open(filepath.Join(helper.Config.SaveDir, imgPath))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		io.Copy(w, f)
	} else {
		http.Redirect(w, r, helper.Config.ResourceServerBaseURL+imgPath, 301)
	}
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

	var hash models.HashInfo

	hash, err = helper.CalculateBasicHashes(savePath)

	if err != nil {
		return
	}

	URL := helper.Config.BaseURL + "/img/" + newID.Hex()

	info = models.ImageInfo{
		ID:        newID,
		Name:      part.FileName(),
		Extension: filepath.Ext(part.FileName()),
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
