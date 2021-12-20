package reader

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	BotToken   string `json:"bot-token" validate:"gt=0"`
	YahooToken string `json:"yahoo-finance-token" validate:"gt=0"`
}

func GetConfig() (*Config, error) {
	jsonFile, err := os.Open("config/.env.json")
	if err != nil {
		log.Println(err)
		return nil, err
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
		return nil, err
	}
	config := &Config{}
	if err := json.Unmarshal(byteValue, config); err != nil {
		log.Println(err)
		return nil, err
	}
	validate := validator.New()
	err = validate.Var(config, "required")
	if err != nil {
		panic(err)
	}
	return config, nil
}
