package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/yangsibai/panda/db"
	"github.com/yangsibai/panda/helper"
	"github.com/yangsibai/panda/models"
	"gopkg.in/mgo.v2/bson"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func getSimpleContentTypeByFileName(filename string) string {
	if filepath.Ext(filename) == ".mp3" {
		return "audio/mp3"
	}
	return ""
}

// save a single image
func handleSaveSingleFile(part *multipart.Part) (info models.FileInfo, err error) {
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

	var hash models.HashInfo

	hash, err = helper.CalculateBasicHashes(savePath)

	if err != nil {
		return
	}

	URL := helper.Config.BaseURL + "/file/" + newID.Hex()

	info = models.FileInfo{
		ID:          newID,
		Name:        part.FileName(),
		Extension:   filepath.Ext(part.FileName()),
		Path:        path,
		URL:         URL,
		Hash:        hash,
		Size:        bytes,
		CreatedAt:   time.Now(),
		ContentType: getSimpleContentTypeByFileName(part.FileName()),
	}
	err = db.StoreResource(&info)
	if err != nil {
		return
	}
	return info, nil
}

func HandleSingleFileUpload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	f, err := handleSaveSingleFile(part)
	if err != nil {
		helper.WriteErrorResponse(w, err)
		return
	}
	helper.WriteResponse(w, f)
}

// get single file
func HandleFetchSingleFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	info, err := db.GetFile(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if helper.Config.ResourceServerBaseURL == "" {
		f, err := os.Open(filepath.Join(helper.Config.SaveDir, info.Path))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		w.Header().Set("Content-Type", info.ContentType)
		io.Copy(w, f)
	} else {
		http.Redirect(w, r, helper.Config.ResourceServerBaseURL+info.Path, 301)
	}
}
