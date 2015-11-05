package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

func main() {
	confFilePath := "./config.toml"

	c := &Config{
		User:   "reiki4040",
		Github: "https://github.com/reiki4040",
	}

	if err := StoreConfig(c, confFilePath); err != nil {
		log.Fatal("save toml error: %s", err.Error())
	}

	c2, err := LoadConfig(confFilePath)
	if err != nil {
		log.Fatal("load toml error: %s", err.Error())
	}

	log.Printf("user: %s, github: %s", c2.User, c2.Github)
}

type Config struct {
	User   string
	Github string
}

func LoadConfig(configFilePath string) (*Config, error) {
	var conf Config
	if _, err := toml.DecodeFile(configFilePath, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

// TODO file swap? backup conf file.
func StoreConfig(config *Config, confFilePath string) error {
	confFile, err := os.Create(confFilePath)
	if err != nil {
		return err
	}
	defer confFile.Close()

	//w := bufio.NewWriter(confFile)
	// wrap with bufio.NewWriter in toml
	enc := toml.NewEncoder(confFile)
	if err := enc.Encode(config); err != nil {
		return err
	}

	return nil
}
