# Scaffold

A powerful Go library for templating and scaffolding projects with customizable token replacement and path-specific transformations.

## Features

- Replace tokens in file names and content
- Apply multiple modifiers to token values (lower, upper, slug, snake, pascal, camel)
- Path-specific token application
- Token value binding
- Simple and intuitive configuration via TOML

## Installation

```bash
go get github.com/trevorak/scaffold/v2
```

### VS Code Extension

For a better experience when creating templates, you can use the [Scaffold Token Highlighter](https://github.com/trevorak/scaffold-token-highlighter-vs) VS Code extension. This extension provides:

- Syntax highlighting for scaffold tokens in your template files
- Easier template creation and maintenance

To install the extension, search for "Scaffold Token Highlighter" in the VS Code marketplace or visit the [extension repository](https://github.com/trevorak/scaffold-token-highlighter-vs).

## Usage

### Template Structure

Create a template directory containing all files and directories that will serve as your template. The structure should look like this:

```
example-template/
├── foo.txt
├── subfolder/
│   └── foo.txt
└── scaffold.toml
```

### Configuration

The `scaffold.toml` file defines tokens, their modifiers, and path-specific rules. Here's an example:

```toml
[[token]]
# The token to replace in files and directories
name = "camelToken"

# Available modifiers: lower, upper, slug, snake, pascal, camel
modifiers = ["camel"]

# Optional: Restrict token application to specific paths
#localize = ["path/to/dir", "path/to/another", "path/to/a/file.go"]

# Token names can also use regular expressions for more flexible matching
# For example: "user[0-9]+" will match "user1", "user42", etc.

[[token]]
name = "PascalToken"
# Bind to another token's value
token = "camelToken"
modifiers = ["pascal"]

[[token]]
name = "slug-token"
modifiers = ["slug"]

[[token]]
name = "foo"
modifiers = ["camel"]
# Apply only to items within subfolder
localize = ["subfolder"]
```

### Available Modifiers

- `lower`: Convert to lowercase
- `upper`: Convert to uppercase
- `slug`: Convert to URL-friendly slug
- `snake`: Convert to snake_case
- `pascal`: Convert to PascalCase
- `camel`: Convert to camelCase
- `singular`: Convert to singular form
- `plural`: Convert to plural form

### Code Example

```go
package main

import (
    "github.com/trevorak/scaffold/v2"
    "log"
)

func main() {
    // Initialize scaffold with template directory
    scaf, err := scaffold.Init("example-template")
    if err != nil {
        log.Fatal("Failed to initialize:", err)
    }

    // Register token values
    scaf.RegisterTokenValue("camelToken", "scaffoldTesting")
    scaf.RegisterTokenValue("foo", "bar")
    scaf.RegisterTokenValue("slug-token", "SomeCrazy String!")

    // Generate scaffolded content
    err = scaf.Make("destination/path")
    if err != nil {
        log.Fatal("Failed to generate scaffold:", err)
    }
}
```
