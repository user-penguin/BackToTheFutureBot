package reader

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Token string `json:"token"`
}

func GetTokenFromFile() (string, error) {
	jsonFile, err := os.Open("config/.env.json")
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.Println(err)
		}
	}(jsonFile)
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Println(err)
		return "", err
	}
	config := &Config{}
	if err := json.Unmarshal(byteValue, config); err != nil {
		log.Println(err)
		return "", err
	}
	return config.Token, nil
}
