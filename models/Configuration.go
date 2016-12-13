package models

type Configuration struct {
	Addr                  string   `json: "addr"`
	SaveDir               string   `json: "saveDir"`
	BaseURL               string   `json: "baseURL"`
	ResourceServerBaseURL string   `json: "resourceServerBaseURL"`
	MongoURL              string   `json: "mongo"`
	CorHosts              []string `json: "corHosts"`
	MaxSize               int64    `json: "maxSize"`
	Oss2Endpoint          string   `json: "oss2Endpoint"`
	Oss2AccessKeyId       string   `json: "oss2AccessKeyId"`
	Oss2AccessKeySecret   string   `json: "oss2AccesskeySecret"`
	Oss2BucketName        string   `json: "oss2BucketName"`
}
