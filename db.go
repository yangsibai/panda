package main

import (
	"gopkg.in/mgo.v2"
)

func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial(config.MongoURL)

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	return s
}

// store image to db
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
