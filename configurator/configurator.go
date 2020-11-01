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
		fmt.Println("Error::Configurator::SetConfig::WithFileOpen::", err)
	}

	decoder := json.NewDecoder(file)
	configuration := Config{}
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("Error::Configurator::SetConfig::WithDecode::", err)
	}

	err = file.Close()
	if err != nil {
		fmt.Println("Error::Configurator::SetConfig::WithFileClose::", err)
	}

	return configuration, nil
}

func (c *Configurator) SetConfig(configuration *Config) error {
	file, err := os.OpenFile("conf.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Println("Error::Configurator::SetConfig::WithFileOpen::", err)
	}

	err = file.Truncate(0)
	if err != nil {
		fmt.Println("Error::Configurator::SetConfig::WithFileTruncate::", err)
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(configuration)
	if err != nil {
		fmt.Println("Error::Configurator::SetConfig::WithEncode::", err)
	}

	err = file.Close()
	if err != nil {
		fmt.Println("Error::Configurator::SetConfig::WithFileClose::", err)
	}

	return err
}
