package scaffold

import (
	"cmp"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

const (
	configFileName = "scaffold.toml"
)

type Scaffold struct {
	Path          string
	Config        Config
	Modifiers     modifierMap
	TokenValueMap map[string]string
	onMakeFunc    func(string)
}

func Init(templatesPath string) (*Scaffold, error) {
	config, err := getConfig(templatesPath + "/" + configFileName)
	if err != nil {
		return nil, err
	}

	scaffold := &Scaffold{
		Path:          templatesPath,
		Config:        config,
		Modifiers:     make(modifierMap),
		TokenValueMap: make(map[string]string),
	}

	scaffold.onMakeFunc = func(_ string) {}

	scaffold.registerDefaultModifiers()

	return scaffold, nil
}

func (scaf *Scaffold) RegisterModifier(tokenName string, modifier func(string) string) {
	scaf.Modifiers.Add(tokenName, modifier)
}

func (scaf *Scaffold) RegisterTokenValue(tokenName string, value string) {
	scaf.TokenValueMap[tokenName] = value
}

func (scaf *Scaffold) registerDefaultModifiers() {
	scaf.RegisterModifier("lower", ModifierLower)
	scaf.RegisterModifier("upper", ModifierUpper)
	scaf.RegisterModifier("slug", ModifierSlug)
	scaf.RegisterModifier("title", ModifierTitle)
	scaf.RegisterModifier("snake", ModifierSnake)
	scaf.RegisterModifier("camel", ModifierCamel)
	scaf.RegisterModifier("pascal", ModifierPascal)
	scaf.RegisterModifier("plural", ModifierPlural)
	scaf.RegisterModifier("singular", ModifierSingular)
}

func (scaf *Scaffold) GetTokens() []Token {
	return scaf.Config.Tokens
}

func (scaf *Scaffold) OnMake(onMakeFunc func(string)) {
	scaf.onMakeFunc = onMakeFunc
}

func (scaf *Scaffold) Make(destination string) error {
	// First pass: Set all token values
	for i := range scaf.Config.Tokens {
		token := &scaf.Config.Tokens[i]

		// If token depends on another token, get its value
		if token.Token != "" {
			parentToken, err := scaf.GetTokenByName(token.Token)
			if err == nil {
				token.Value = parentToken.Value
			}
		}

		// If no value is set yet, try to get it from TokenValueMap (user-supplied values)
		if token.Value == "" {
			token.Value = scaf.TokenValueMap[token.Name]
		}

		// Apply modifiers
		for _, modifier := range token.Modifiers {
			for _, modFunc := range scaf.Modifiers[modifier] {
				token.Value = modFunc(token.Value)
			}
		}
	}

	// Second pass: Process any remaining token dependencies
	for i := range scaf.Config.Tokens {
		token := &scaf.Config.Tokens[i]
		if token.Token != "" && token.Value == "" {
			parentToken, err := scaf.GetTokenByName(token.Token)
			if err == nil && parentToken.Value != "" {
				token.Value = parentToken.Value
				// Apply modifiers again for newly set values
				for _, modifier := range token.Modifiers {
					for _, modFunc := range scaf.Modifiers[modifier] {
						token.Value = modFunc(token.Value)
					}
				}
			}
		}
	}

	// Sort tokens by priority for replacement order
	slices.SortFunc(scaf.Config.Tokens, func(a, b Token) int {
		return cmp.Compare(b.Priority, a.Priority)
	})

	// if scaf.Path does not exist, create it
	if _, err := os.Stat(scaf.Path); os.IsNotExist(err) {
		err := os.MkdirAll(scaf.Path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	if err := filepath.WalkDir(scaf.Path, func(path string, info os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if info.Name() == configFileName {
			return nil
		}

		if path == scaf.Path {
			return nil
		}

		relativePath := strings.TrimPrefix(path, scaf.Path)

		relativePath = scaf.replaceTokens(relativePath, path)

		makeDestination := destination + relativePath

		if info.IsDir() {
			return os.MkdirAll(makeDestination, os.ModePerm)
		}

		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		stringcontents := string(contents)
		stringcontents = scaf.replaceTokens(stringcontents, path)

		if err := os.WriteFile(makeDestination, []byte(stringcontents), 0644); err != nil {
			return err
		}

		scaf.onMakeFunc(makeDestination)

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (scaf *Scaffold) replaceTokens(subject string, path string) string {
	for _, token := range scaf.Config.Tokens {
		if token.Localize != nil {
			for _, tokenPath := range token.Localize {
				if strings.HasPrefix(path, scaf.Path+"/"+tokenPath) {
					re := regexp.MustCompile(token.Name)
					subject = re.ReplaceAllString(subject, token.Value)
				}
			}
		} else {
			re := regexp.MustCompile(token.Name)
			subject = re.ReplaceAllString(subject, token.Value)
		}
	}

	return subject
}

func (scaf *Scaffold) GetTokenByName(name string) (Token, error) {
	for _, token := range scaf.Config.Tokens {
		if token.Name == name {
			return token, nil
		}
	}

	return Token{}, errors.New("token not found")
}
