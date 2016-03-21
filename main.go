package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/yangsibai/panda/helper"
	"github.com/yangsibai/panda/routes"
	"log"
	"net/http"
)

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins:   helper.Config.CorHosts,
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST"},
	})

	router := httprouter.New()
	router.GET("/", routes.Index)
	router.POST("/upload/img", routes.HandleSingleImageUpload)
	router.POST("/upload/imgs", routes.HandleMultipleImagesUpload)
	router.POST("/upload/file", routes.HandleSingleFileUpload)
	router.GET("/img/:id", routes.HandleFetchSingleImage)
	router.GET("/file/:id", routes.HandleFetchSingleFile)
	router.GET("/info/:id", routes.HandleGetInfo)

	hanlder := c.Handler(router)
	log.Printf("panda is listenning at %s", helper.Config.Addr)
	log.Fatal(http.ListenAndServe(helper.Config.Addr, hanlder))
}
