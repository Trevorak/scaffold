# Scaffold

## Usage

### Layout
- Create a folder, which contains all files/directories that will be used as a template.
- File content and names will be modified according to the token definitions contained in scaffold.toml, which should exist in the root of the template folder.
- Example
```
~/example-template/foo.txt
~/example-template/subfolder/foo.txt
~/example-template/scaffold.toml
```

### Configuration
- scaffold.toml defines token names and modifiers, can bind a token to the value of another token, and can restrict application to specific paths.
```
[[token]]
# The token to replace in files and directories
name = "camelToken"

# Modifies the value. One of lower, upper, slug, snake, pascal, camel
modifiers = ["camel"]

# specifies the path for a token. This token will only be applied to items within the specified file/directory
#localize = ["path/to/dir", "path/to/another", "path/to/a/file.go"]


[[token]]
name = "PascalToken"
value = "camelToken" # Use the same supplied value as the token camelToken. Order is important here. This must be defined after camelToken.
modifiers = ["pascal"]

[[token]]
name = "slug-token"
modifiers = ["slug"]

[[token]]
name = "foo"
modifiers = ["camel"]
localize = ["subfolder"] # This would apply only to items within subfolder. Renaming foo.txt to the value set for the token.
```

### Execution

```go
    package main

    import (
        "github.com/trevorak/scaffold"
        "log"
    )

    func main() {
		scaf, err := scaffold.Init("example-template")
		if err != nil {
			log.Fatal("Failed to initialize")
		}

		scaf.RegisterTokenValue("camelToken", "scaffoldTesting")
		scaf.RegisterTokenValue("foo", "bar")
		scaf.RegisterTokenValue("slug-token", "SomeCrazy String!")

		err = scaf.Make("destination/path")
		if err != nil {
			log.Fatal(err)
		}
    }
```