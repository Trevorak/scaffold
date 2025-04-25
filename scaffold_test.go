package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMake(t *testing.T) {
	// Create a temporary directory for test templates
	tmpDir, err := os.MkdirTemp("", "scaffold-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test template structure
	templateDir := filepath.Join(tmpDir, "template")
	err = os.MkdirAll(templateDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create template dir: %v", err)
	}

	// Create scaffold.toml
	configContent := `
		[[token]]
		name = "{{name}}"
		value = "test"
		priority = 1

		[[token]]
		name = "{{name_upper}}"
		token = "{{name}}"
		modifiers = ["upper"]
		priority = 2

		[[token]]
		name = "{{name_lower}}"
		token = "{{name_upper}}"
		modifiers = ["lower"]
		priority = 3

		[[token]]
		name = "{{dependent}}"
		token = "{{name_upper}}"
		priority = 4

		[[token]]
		name = "{{chain1}}"
		token = "{{chain2}}"
		priority = 5

		[[token]]
		name = "{{chain2}}"
		token = "{{chain3}}"
		priority = 6

		[[token]]
		name = "{{chain3}}"
		value = "final"
		priority = 7
	`
	err = os.WriteFile(filepath.Join(templateDir, "scaffold.toml"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Create a test file with tokens
	fileContent := `package {{name}}

func {{name_upper}}Func() {
    return "{{dependent}}"
}

var {{name_lower}}Var = "{{chain1}}"
`
	err = os.WriteFile(filepath.Join(templateDir, "test.go"), []byte(fileContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	scaf, err := Init(templateDir)
	if err != nil {
		t.Fatalf("Failed to init scaffold: %v", err)
	}

	// Set up destination directory
	destDir := filepath.Join(tmpDir, "output")
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}

	// Run Make
	err = scaf.Make(destDir)
	if err != nil {
		t.Fatalf("Failed to make scaffold: %v", err)
	}

	// Read and verify the generated file
	generatedContent, err := os.ReadFile(filepath.Join(destDir, "test.go"))
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	expectedContent := `package test

func TESTFunc() {
    return "TEST"
}

var testVar = "final"
`

	if string(generatedContent) != expectedContent {
		t.Errorf("Generated content does not match expected.\nExpected:\n%s\nGot:\n%s", expectedContent, string(generatedContent))
	}
}

func TestMakeWithLocalizedTokens(t *testing.T) {
	// Create a temporary directory for test templates
	tmpDir, err := os.MkdirTemp("", "scaffold-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test template structure
	templateDir := filepath.Join(tmpDir, "template")
	err = os.MkdirAll(templateDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create template dir: %v", err)
	}

	// Create scaffold.toml with localized tokens
	configContent := `
		[[token]]
		name = "{{name}}"
		value = "test"
		priority = 1
		
		[[token]]
		name = "{{local1}}"
		value = "localVal1"
		priority = 2
		localize = ["dir1"]
		
		[[token]]
		name = "{{local2}}"
		value = "localVal2"
		priority = 2
		localize = ["dir2"]
	`
	err = os.WriteFile(filepath.Join(templateDir, "scaffold.toml"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Create test files in different directories
	dir1Content := `{{name}} {{local1}}`
	dir2Content := `{{name}} {{local2}}`

	err = os.MkdirAll(filepath.Join(templateDir, "dir1"), 0755)
	if err != nil {
		t.Fatalf("Failed to create dir1: %v", err)
	}
	err = os.MkdirAll(filepath.Join(templateDir, "dir2"), 0755)
	if err != nil {
		t.Fatalf("Failed to create dir2: %v", err)
	}

	err = os.WriteFile(filepath.Join(templateDir, "dir1", "test.txt"), []byte(dir1Content), 0644)
	if err != nil {
		t.Fatalf("Failed to write dir1 file: %v", err)
	}
	err = os.WriteFile(filepath.Join(templateDir, "dir2", "test.txt"), []byte(dir2Content), 0644)
	if err != nil {
		t.Fatalf("Failed to write dir2 file: %v", err)
	}

	// Initialize scaffold
	scaf, err := Init(templateDir)
	if err != nil {
		t.Fatalf("Failed to init scaffold: %v", err)
	}

	// Set up destination directory
	destDir := filepath.Join(tmpDir, "output")
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}

	// Run Make
	err = scaf.Make(destDir)
	if err != nil {
		t.Fatalf("Failed to make scaffold: %v", err)
	}

	// Verify dir1 content
	dir1Generated, err := os.ReadFile(filepath.Join(destDir, "dir1", "test.txt"))
	if err != nil {
		t.Fatalf("Failed to read dir1 generated file: %v", err)
	}
	if string(dir1Generated) != "test localVal1" {
		t.Errorf("Dir1 content incorrect. Expected 'test local1', got '%s'", string(dir1Generated))
	}

	// Verify dir2 content
	dir2Generated, err := os.ReadFile(filepath.Join(destDir, "dir2", "test.txt"))
	if err != nil {
		t.Fatalf("Failed to read dir2 generated file: %v", err)
	}
	if string(dir2Generated) != "test localVal2" {
		t.Errorf("Dir2 content incorrect. Expected 'test local2', got '%s'", string(dir2Generated))
	}
}

func TestMakeWithRegisteredValues(t *testing.T) {
	// Create a temporary directory for test templates
	tmpDir, err := os.MkdirTemp("", "scaffold-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test template structure
	templateDir := filepath.Join(tmpDir, "template")
	err = os.MkdirAll(templateDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create template dir: %v", err)
	}

	// Create scaffold.toml
	configContent := `
		[[token]]
		name = "{{name}}"
		priority = 1

		[[token]]
		name = "{{name_upper}}"
		token = "{{name}}"
		modifiers = ["upper"]
		priority = 2

		[[token]]
		name = "{{name_lower}}"
		token = "{{name}}"
		modifiers = ["lower"]
		priority = 3

		[[token]]
		name = "{{custom}}"
		priority = 4
	`
	err = os.WriteFile(filepath.Join(templateDir, "scaffold.toml"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Create a test file with tokens
	fileContent := `package {{name}}

func {{name_upper}}Func() {
    return "{{custom}}"
}

var {{name_lower}}Var = "{{name}}"
`
	err = os.WriteFile(filepath.Join(templateDir, "test.go"), []byte(fileContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Initialize scaffold
	scaf, err := Init(templateDir)
	if err != nil {
		t.Fatalf("Failed to init scaffold: %v", err)
	}

	// Register token values programmatically
	scaf.RegisterTokenValue("{{name}}", "myapp")
	scaf.RegisterTokenValue("{{custom}}", "custom_value")

	// Set up destination directory
	destDir := filepath.Join(tmpDir, "output")
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}

	// Run Make
	err = scaf.Make(destDir)
	if err != nil {
		t.Fatalf("Failed to make scaffold: %v", err)
	}

	// Read and verify the generated file
	generatedContent, err := os.ReadFile(filepath.Join(destDir, "test.go"))
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	expectedContent := `package myapp

func MYAPPFunc() {
    return "custom_value"
}

var myappVar = "myapp"
`

	if string(generatedContent) != expectedContent {
		t.Errorf("Generated content does not match expected content.\nExpected:\n%s\nGot:\n%s", expectedContent, string(generatedContent))
	}
}

func TestMakeWithCustomModifier(t *testing.T) {
	// Create a temporary directory for test templates
	tmpDir, err := os.MkdirTemp("", "scaffold-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test template structure
	templateDir := filepath.Join(tmpDir, "template")
	err = os.MkdirAll(templateDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create template dir: %v", err)
	}

	// Create scaffold.toml
	configContent := `
		[[token]]
		name = "{{name}}"
		priority = 1

		[[token]]
		name = "{{name_reversed}}"
		token = "{{name}}"
		modifiers = ["reverse"]
		priority = 2
	`
	err = os.WriteFile(filepath.Join(templateDir, "scaffold.toml"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Create a test file with tokens
	fileContent := `package {{name}}

func {{name_reversed}}Func() {
    return "{{name_reversed}}"
}
`
	err = os.WriteFile(filepath.Join(templateDir, "test.go"), []byte(fileContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Initialize scaffold
	scaf, err := Init(templateDir)
	if err != nil {
		t.Fatalf("Failed to init scaffold: %v", err)
	}

	// Register a custom modifier that reverses strings
	scaf.RegisterModifier("reverse", func(s string) string {
		runes := []rune(s)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	})

	// Register token value
	scaf.RegisterTokenValue("{{name}}", "myapp")

	// Set up destination directory
	destDir := filepath.Join(tmpDir, "output")
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}

	// Run Make
	err = scaf.Make(destDir)
	if err != nil {
		t.Fatalf("Failed to make scaffold: %v", err)
	}

	// Read and verify the generated file
	generatedContent, err := os.ReadFile(filepath.Join(destDir, "test.go"))
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	expectedContent := `package myapp

func ppaymFunc() {
    return "ppaym"
}
`

	if string(generatedContent) != expectedContent {
		t.Errorf("Generated content does not match expected content.\nExpected:\n%s\nGot:\n%s", expectedContent, string(generatedContent))
	}
}
