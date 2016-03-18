package db

import (
	"github.com/yangsibai/panda/helper"
	"github.com/yangsibai/panda/models"
	"gopkg.in/mgo.v2"
)

func GetSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial(helper.Config.MongoURL)

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	return s
}

// store image to db
func StoreImage(info *models.ImageInfo) (err error) {
	session := GetSession()
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	C := session.DB("resource").C("image")
	err = C.Insert(&info)
	return
}
