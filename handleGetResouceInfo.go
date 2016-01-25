package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func GetResouceInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	session := getSession()
	C := session.DB("resource").C("image")
	defer session.Close()

	info := ImageInfo{}

	oid := bson.ObjectIdHex(id)
	if err := C.FindId(oid).One(&info); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jo, err := json.Marshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", jo)
}
