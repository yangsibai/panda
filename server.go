package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "hello")
}

var config Configuration

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins:   config.CorHosts,
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST"},
	})

	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/upload/img", handleSingleImageUpload)
	router.POST("/upload/imgs", handleMultipleImagesUpload)
	router.GET("/img/:id", handleFetchSingleImage)
	router.GET("/info/:id", handleGetInfo)

	hanlder := c.Handler(router)
	log.Fatal(http.ListenAndServe(config.Addr, hanlder))
}

type Configuration struct {
	Addr     string   `json: "addr"`
	SaveDir  string   `json: "saveDir"`
	BaseURL  string   `json: "baseURL"`
	MongoURL string   `json: "mongo"`
	CorHosts []string `json: "corHosts"`
	MaxSize  int64    `json:"maxSize"`
}

func readConfig() (error, Configuration) {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
		return err, configuration
	}
	return nil, configuration
}

func init() {
	var err error
	err, config = readConfig()
	if err != nil {
		log.Fatal(err)
	}
}
