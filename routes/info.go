package routes

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/yangsibai/panda/db"
	"github.com/yangsibai/panda/models"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

// get image info by id
func HandleGetInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	session := db.GetSession()
	C := session.DB("resource").C("image")
	defer session.Close()

	info := models.ImageInfo{}

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jo)
}
