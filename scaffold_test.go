package scaffold

import (
	"fmt"
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

func TestMakeWithRegexToken(t *testing.T) {
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

	// Create scaffold.toml with various regex tokens
	configContent := `
		[[token]]
		name = "{{\\w+_id}}"
		value = "12345"
		priority = 1

		[[token]]
		name = "{{prefix_\\w+}}"
		value = "prefixed"
		priority = 2

		[[token]]
		name = "{{\\w+_suffix}}"
		value = "suffixed"
		priority = 3

		[[token]]
		name = "{{test_\\d+}}"
		value = "numbered"
		priority = 4

		[[token]]
		name = "{{[a-z]+_[A-Z]+}}"
		value = "mixed_case"
		priority = 5

		[[token]]
		name = "{{name}}"
		value = "test"
		priority = 6
	`
	err = os.WriteFile(filepath.Join(templateDir, "scaffold.toml"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Create a test file with various regex token patterns
	fileContent := `package {{name}}

const (
    UserID     = "{{user_id}}"
    OrderID    = "{{order_id}}"
    ProductID  = "{{product_id}}"
    PrefixVal  = "{{prefix_value}}"
    ValSuffix  = "{{value_suffix}}"
    TestNum    = "{{test_123}}"
    CaseMix    = "{{lower_UPPER}}"
)

type Config struct {
    SessionID  string = "{{session_id}}"
    PrefixKey string = "{{prefix_key}}"
    DataSuffix string = "{{data_suffix}}"
    TestCase   string = "{{test_999}}"
    MixedCase  string = "{{snake_CASE}}"
}
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

const (
    UserID     = "12345"
    OrderID    = "12345"
    ProductID  = "12345"
    PrefixVal  = "prefixed"
    ValSuffix  = "suffixed"
    TestNum    = "numbered"
    CaseMix    = "mixed_case"
)

type Config struct {
    SessionID  string = "12345"
    PrefixKey string = "prefixed"
    DataSuffix string = "suffixed"
    TestCase   string = "numbered"
    MixedCase  string = "mixed_case"
}
`

	if string(generatedContent) != expectedContent {
		t.Errorf("Generated content does not match expected.\nExpected:\n%s\nGot:\n%s", expectedContent, string(generatedContent))
	}
}

func BenchmarkMakeBasic(b *testing.B) {
	// Create a temporary directory for test templates
	tmpDir, err := os.MkdirTemp("", "scaffold-bench-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Set up template directory
	templateDir := filepath.Join(tmpDir, "template")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		b.Fatalf("Failed to create template dir: %v", err)
	}

	// Create basic scaffold.toml
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
    `
	if err := os.WriteFile(filepath.Join(templateDir, "scaffold.toml"), []byte(configContent), 0644); err != nil {
		b.Fatalf("Failed to write config: %v", err)
	}

	// Create a test file
	fileContent := `package {{name}}

func {{name_upper}}Func() string {
    return "{{name}}"
}
`
	if err := os.WriteFile(filepath.Join(templateDir, "test.go"), []byte(fileContent), 0644); err != nil {
		b.Fatalf("Failed to write test file: %v", err)
	}

	scaf, err := Init(templateDir)
	if err != nil {
		b.Fatalf("Failed to init scaffold: %v", err)
	}

	// Set up destination directory
	destDir := filepath.Join(tmpDir, "output")
	if err := os.MkdirAll(destDir, 0755); err != nil {
		b.Fatalf("Failed to create output dir: %v", err)
	}

	// Prepare all directories before measurement starts
	iterDirs := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		iterDirs[i] = filepath.Join(destDir, fmt.Sprintf("iter_%d", i))
		if err := os.MkdirAll(iterDirs[i], 0755); err != nil {
			b.Fatalf("Failed to create iteration dir: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := scaf.Make(filepath.Join(destDir, fmt.Sprintf("iter_%d", i))); err != nil {
			b.Fatalf("Failed to make scaffold: %v", err)
		}
	}
}

func BenchmarkMakeWithComplexTokens(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "scaffold-bench-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	templateDir := filepath.Join(tmpDir, "template")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		b.Fatalf("Failed to create template dir: %v", err)
	}

	// Create complex scaffold.toml with chained and regex tokens
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
        name = "{{chain1}}"
        token = "{{chain2}}"
        priority = 3

        [[token]]
        name = "{{chain2}}"
        token = "{{chain3}}"
        priority = 4

        [[token]]
        name = "{{chain3}}"
        value = "final"
        priority = 5

        [[token]]
        name = "{{\\w+_id}}"
        value = "12345"
        priority = 6
        is_regex = true
    `
	if err := os.WriteFile(filepath.Join(templateDir, "scaffold.toml"), []byte(configContent), 0644); err != nil {
		b.Fatalf("Failed to write config: %v", err)
	}

	// Create a complex test file
	fileContent := `package {{name}}

type Data struct {
    UserID    string = "{{user_id}}"
    OrderID   string = "{{order_id}}"
    Name      string = "{{name_upper}}"
    ChainVal  string = "{{chain1}}"
}
`
	if err := os.WriteFile(filepath.Join(templateDir, "test.go"), []byte(fileContent), 0644); err != nil {
		b.Fatalf("Failed to write test file: %v", err)
	}

	scaf, err := Init(templateDir)
	if err != nil {
		b.Fatalf("Failed to init scaffold: %v", err)
	}

	destDir := filepath.Join(tmpDir, "output")
	if err := os.MkdirAll(destDir, 0755); err != nil {
		b.Fatalf("Failed to create output dir: %v", err)
	}

	// Prepare all directories before measurement starts
	iterDirs := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		iterDirs[i] = filepath.Join(destDir, fmt.Sprintf("iter_%d", i))
		if err := os.MkdirAll(iterDirs[i], 0755); err != nil {
			b.Fatalf("Failed to create iteration dir: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := scaf.Make(filepath.Join(destDir, fmt.Sprintf("iter_%d", i))); err != nil {
			b.Fatalf("Failed to make scaffold: %v", err)
		}
	}
}

func BenchmarkMakeWithCustomModifiers(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "scaffold-bench-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	templateDir := filepath.Join(tmpDir, "template")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		b.Fatalf("Failed to create template dir: %v", err)
	}

	configContent := `
        [[token]]
        name = "{{name}}"
        value = "test"
        priority = 1

        [[token]]
        name = "{{name_reversed}}"
        token = "{{name}}"
        modifiers = ["reverse"]
        priority = 2
    `
	if err := os.WriteFile(filepath.Join(templateDir, "scaffold.toml"), []byte(configContent), 0644); err != nil {
		b.Fatalf("Failed to write config: %v", err)
	}

	fileContent := `package {{name}}

func {{name_reversed}}Func() string {
    return "{{name_reversed}}"
}
`
	if err := os.WriteFile(filepath.Join(templateDir, "test.go"), []byte(fileContent), 0644); err != nil {
		b.Fatalf("Failed to write test file: %v", err)
	}

	scaf, err := Init(templateDir)
	if err != nil {
		b.Fatalf("Failed to init scaffold: %v", err)
	}

	scaf.RegisterModifier("reverse", func(s string) string {
		runes := []rune(s)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	})

	destDir := filepath.Join(tmpDir, "output")
	if err := os.MkdirAll(destDir, 0755); err != nil {
		b.Fatalf("Failed to create output dir: %v", err)
	}

	// Prepare all directories before measurement starts
	iterDirs := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		iterDirs[i] = filepath.Join(destDir, fmt.Sprintf("iter_%d", i))
		if err := os.MkdirAll(iterDirs[i], 0755); err != nil {
			b.Fatalf("Failed to create iteration dir: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := scaf.Make(filepath.Join(destDir, fmt.Sprintf("iter_%d", i))); err != nil {
			b.Fatalf("Failed to make scaffold: %v", err)
		}
	}
}
