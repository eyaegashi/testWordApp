package config

import (
	"encoding/json"

	"github.com/eyaegashi/wordTestApp/util"
)

const configFileName = "./config.json"

var conf Config

// Config is the struct for config file
type Config struct {
	ready   bool
	WordAPI struct {
		APIID  string `json:"APIID"`
		APIKey string `json:"APIKey"`
		URL    string `json:"Url"`
	} `json:"WordAPI"`
}

// setConfig is to load config file and set it
func setConfig() (err error) {
	jsonData, err := util.LoadjsonFile(configFileName)
	if err != nil {
		return err
	}

	// convert json to struct
	err = json.Unmarshal(jsonData, &conf)
	if err != nil {
		return err
	}
	conf.ready = true
	return nil
}

//GetConfig is to get config
func GetConfig() *Config {
	if !conf.ready {
		err := setConfig()
		if err != nil {
			// todo: config error: cannot get onfig informaiton
		}
	}
	return &conf
}
