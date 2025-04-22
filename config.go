package scaffold

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

type Token struct {
	Name       string `toml:"name"`
	Value      string
	ValueToken string   `toml:"value"`
	Modifiers  []string `toml:"modifiers"`
	Localize   []string `toml:"localize"`
	Priority   int      `toml:"priority"`
}

type Config struct {
	Tokens []Token `toml:"token"`
}

func getConfig(configPath string) (Config, error) {
	config := Config{}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	err = toml.Unmarshal(content, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
