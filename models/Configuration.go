package models

type Configuration struct {
	Addr     string   `json: "addr"`
	SaveDir  string   `json: "saveDir"`
	BaseURL  string   `json: "baseURL"`
	MongoURL string   `json: "mongo"`
	CorHosts []string `json: "corHosts"`
	MaxSize  int64    `json:"maxSize"`
}
