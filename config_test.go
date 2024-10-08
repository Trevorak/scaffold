package scaffold

import (
	"log"
	"os"
	"testing"
)

func TestGetConfig(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	config, err := getConfig(wd + "/_testdata/template/scaffold.toml")
	if err != nil {
		log.Fatal(err)
	}

	if len(config.Tokens) == 0 {
		log.Fatal("No tokens found in config")
	}
}
