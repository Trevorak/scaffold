package scaffold

import (
	"log"
	"os"
	"strings"
	"testing"
)

func TestScaffold_Make(t *testing.T) {
	scaffold, err := Init("_testdata/template")
	if err != nil {
		log.Fatal("Failed to initialize")
	}

	scaffold.RegisterTokenValue("camelToken", "scaffoldTesting")
	scaffold.RegisterTokenValue("foo", "bar")
	scaffold.RegisterTokenValue("Foo", "Bar")
	scaffold.RegisterTokenValue("slug-token", "SomeCrazy String!")

	err = scaffold.Make("_testdata/destination")
	if err != nil {
		log.Fatal(err)
	}

	sampleFileBytes, err := os.ReadFile("_testdata/destination/sample_file.go")
	if err != nil {
		log.Fatal(err)
	}

	for _, token := range scaffold.GetTokens() {
		if strings.Contains(string(sampleFileBytes), token.Name) && token.Localize == nil {
			log.Fatalf("Token %v was not replaced in sample_file.go", token.Name)
		}

		if strings.Contains(string(sampleFileBytes), token.Value) != true && token.Localize == nil {
			log.Fatalf("sample_file.go content does not contain expected text \"%v\"", token.Value)
		}
	}

	barFileBytes, err := os.ReadFile("_testdata/destination/local/bar.go")
	if err != nil {
		log.Fatal(err)
	}

	if strings.Contains(string(barFileBytes), "type Bar") != true {
		log.Fatal("local/bar.go content does not contain expected text \"type Bar\"")
	}

	if strings.Contains(string(barFileBytes), "getBar()") != true {
		log.Fatal("local/bar.go content does not contain expected text \"getBar()\"")
	}
}
