package routes

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/yangsibai/panda/db"
	"github.com/yangsibai/panda/helper"
	"github.com/yangsibai/panda/models"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func getImgFilePath(path, ext string, width int) (imgPath string, err error) {
	if width == 0 {
		return path, nil
	}
	imgPath = fmt.Sprintf(path+"_w_%d", width)

	originalAbsolutePath := filepath.Join(helper.Config.SaveDir, path)
	newAbsolutePath := filepath.Join(helper.Config.SaveDir, imgPath)

	if _, err := os.Stat(newAbsolutePath); os.IsNotExist(err) {
		err = helper.CreateThumbnail(originalAbsolutePath, ext, newAbsolutePath, uint(width))
		return newAbsolutePath, err
	}
	return newAbsolutePath, nil
}

// get single file
func HandleFetchSingleFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	info, err := db.GetFile(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println(info.Extension, info.Extension == ".png")
	if strings.ToLower(info.Extension) == ".png" || strings.ToLower(info.Extension) == ".jpg" {
		width := getWidth(r)
		log.Println("width", width)
		imgPath, err := getImgFilePath(info.Path, info.Extension, width)
		log.Println("error", err, "img path", imgPath)
		if err != nil {
			f, err := os.Open(imgPath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer f.Close()

			w.Header().Set("Content-Type", info.ContentType)
			io.Copy(w, f)
			return
		}
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

func getWidth(r *http.Request) int {
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
	return width
}
