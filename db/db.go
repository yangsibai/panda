package db

import (
	"github.com/yangsibai/panda/helper"
	"github.com/yangsibai/panda/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const DB_NAME string = "resource"
const COLLECTION_IMAGE_NAME string = "image"
const COLLECTION_FILE_NAME string = "file"

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

	C := session.DB(DB_NAME).C(COLLECTION_IMAGE_NAME)
	err = C.Insert(&info)
	return
}

func GetImage(id string) (info models.ImageInfo, err error) {
	session := GetSession()
	C := session.DB(DB_NAME).C(COLLECTION_IMAGE_NAME)
	defer session.Close()

	oid := bson.ObjectIdHex(id)
	err = C.FindId(oid).One(&info)
	return
}

func StoreResource(info *models.FileInfo) (err error) {
	session := GetSession()
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	C := session.DB(DB_NAME).C(COLLECTION_FILE_NAME)
	err = C.Insert(&info)
	return
}

func GetFile(id string) (info models.FileInfo, err error) {
	session := GetSession()
	C := session.DB(DB_NAME).C(COLLECTION_FILE_NAME)
	defer session.Close()

	oid := bson.ObjectIdHex(id)
	err = C.FindId(oid).One(&info)
	return
}
