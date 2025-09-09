package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DB_URL           string `json:"db_url"`
	Current_Username string `json:"current_username"`
}

const configFileName string = ".gatorconfig.json"

func getConfigFilePath() (string, error) {

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filePath := home + "/" + configFileName

	return filePath, nil
}

func (c *Config) SetUsername(username string) {
	c.Current_Username = username
	c.write()
}

func (c *Config) write() error {

	filepath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, data, os.FileMode(0644))
	if err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {

	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	var c = Config{}
	if err := json.Unmarshal(data, &c); err != nil {
		return Config{}, err
	}

	return c, nil
}
