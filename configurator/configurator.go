package configurator

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	IsUsed       bool     `json:"setup"`
	InboxPath    string   `json:"inbox"`
	FilteredPath string   `json:"filtered"`
	SpamPath     string   `json:"spam"`
	TargetEmail  string   `json:"targetEmail"`
	ListOfMails  []string `json:"listOfMails"`
}

type Configurator struct{}

func (c *Configurator) GetConfig() (Config, error) {
	file, err := os.Open("conf.json")
	if err != nil {
		fmt.Println(err)
	}

	decoder := json.NewDecoder(file)
	configuration := Config{}
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println(err)
	}

	err = file.Close()
	if err != nil {
		fmt.Println(err)
	}

	return configuration, nil
}

func (c *Configurator) SetConfig(configuration *Config) error {
	file, err := os.OpenFile("conf.json", os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println(err)
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(configuration)
	if err != nil {
		fmt.Println(err)
	}

	err = file.Close()
	if err != nil {
		fmt.Println(err)
	}

	return err
}
