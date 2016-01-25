package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	_ "github.com/nfnt/resize"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"panda/helper"
	"path/filepath"
	"strconv"
)

func HandleSingleImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	session := getSession()
	C := session.DB("resource").C("image")
	defer session.Close()

	info := ImageInfo{}

	oid := bson.ObjectIdHex(id)
	if err := C.FindId(oid).One(&info); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resize_key := fmt.Sprintf("w_%d", width)
	if val, ok := info.Resizes[resize_key]; ok {
		http.Redirect(w, r, val, 301)
		return
	}

	newPath := fmt.Sprintf(info.Path+"_w_%d", width)
	newAbsolutePath := filepath.Join(config.SaveDir, newPath)
	err = helper.CreateThumbnail(filepath.Join(config.SaveDir, info.Path), info.Extension, newAbsolutePath, uint(width))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if info.Resizes == nil {
		info.Resizes = map[string]string{}
	}
	info.Resizes[resize_key] = config.BaseURL + newPath
	change := bson.M{"$set": bson.M{"resizes": info.Resizes}}
	err = C.Update(bson.M{"_id": oid}, change)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, config.BaseURL+newPath, 301)
}
