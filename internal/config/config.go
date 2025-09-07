package config
import "os"

type Config struct {
}

const configFileName string = ".gatorconfig.json"

func getConfigFilePath() string, error {
	filePath := ""

	os.UserHomeDir

	return filePath, nil
}

func (c *Config) SetUsername() {

}

func ReadJson() Config {

	return Config{}
}
