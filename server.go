package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"log"
	"os"
)

func main() {
	err, config := readConfig()
	if err != nil {
		log.Fatal(err)
		return
	}
	m := martini.Classic()

	m.Get("/", func() string {
		return "Yes, commander!"
	})

	m.Post("/upload/img", handleImageUpload)
	m.RunOnAddr(config.Addr)
}

type Configuration struct {
	Addr    string
	SaveDir string
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
