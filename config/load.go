package config

import (
	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
	"io/ioutil"
)

var (
	conf Config
)

// load all the config for the application
func LoadConfig() Config {
	conf, err := readEnv()
	if err != nil {
		log.Fatal("Error While Loading config from Env - " + err.Error())
	}
	return conf
}

// Read Config on App start
func ReadConfig() Config {
	return conf
}

// Read info file to read build release file
func readInfoFile() string {
	data, err := ioutil.ReadFile("data/info.json") // just pass the file name
	if err != nil {
		log.Fatal("Error While opening data/info.json file - " + err.Error())
		return ""
	}
	return string(data)
}

// Read Env on App start
func readEnv() (Config, error) {
	err := envconfig.Process("app", &conf)
	if err != nil {
		return conf, nil
	}
	infoData := readInfoFile()
	conf.SystemInfo = infoData
	return conf, nil
}
