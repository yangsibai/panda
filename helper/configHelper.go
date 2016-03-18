package helper

import (
	"encoding/json"
	"fmt"
	"github.com/yangsibai/panda/models"
	"log"
	"os"
)

var Config models.Configuration

func readConfig() (error, models.Configuration) {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := models.Configuration{}
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
	err, Config = readConfig()
	if err != nil {
		log.Fatal(err)
	}
}
