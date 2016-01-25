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
