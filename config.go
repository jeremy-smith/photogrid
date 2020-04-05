package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const ConfigFile string = "config.json"

type DB struct {
	DbHost string `json:dbhost,omitempty`
	DbUser string `json:dbuser,omitempty`
	DbPass string `json:dbpass,omitempty`
	DbName string `json:dbname,omitempty`
	DbPort string `json:dbport,omitempty`
}

type Config struct {
	AdminIps []string
	Db       DB
}

var Conf Config

func ReadConf(fileName string) Config {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln("Could not read config file: ", err)
	}
	return parseConfig(b)
}

func parseConfig(b []byte) Config {
	c := Config{}
	err := json.Unmarshal(b, &c)
	if err != nil {
		log.Fatalln("Could not parse config: ", err)
	}
	return c
}
