package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type ImageInfo struct {
	ID        bson.ObjectId     `json:"id" bson:"_id"`
	Name      string            `json:"name" bson:"name"`
	BaseDir   string            `json:"-" bson: "baseDir"`
	Path      string            `json:"path" bson:"path" `
	Extension string            `json:"extension" bson:"extension"`
	Width     int               `json:"width" bson:"width"`
	Height    int               `json:"height" bson:"height"`
	URL       string            `json:"URL" bson:"URL"`
	Resizes   map[string]string `json:"resizes" bson:"resizes"`
	Hash      HashInfo          `json:"-" bson:"hash"`
	Size      int64             `json:"size" bson:"size"`
	CreatedAt time.Time         `json:"created" bson: "created"`
}
