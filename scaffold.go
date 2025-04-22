package scaffold

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	configFileName = "scaffold.toml"
)

type modifierMap map[string][]func(string) string

func (modMap *modifierMap) Add(token string, modifier func(string) string) *modifierMap {
	(*modMap)[token] = append((*modMap)[token], modifier)

	return modMap
}

type Scaffold struct {
	Path          string
	Config        Config
	Modifiers     modifierMap
	TokenValueMap map[string]string
	onMakeFunc    func(string)
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
}

func (scaf *Scaffold) GetTokens() []Token {
	return scaf.Config.Tokens
}

func (scaf *Scaffold) OnMake(onMakeFunc func(string)) {
	scaf.onMakeFunc = onMakeFunc
}

func (scaf *Scaffold) Make(destination string) error {
	var localTokens []Token
	var globalTokens []Token
	for i, token := range scaf.Config.Tokens {
		if token.ValueToken == "" {
			token.Value = scaf.TokenValueMap[token.Name]
		} else {
			token.Value = scaf.TokenValueMap[token.ValueToken]
		}

		for _, modifier := range token.Modifiers {
			for j := range scaf.Modifiers[modifier] {
				token.Value = scaf.Modifiers[modifier][j](token.Value)
			}
		}

		scaf.Config.Tokens[i].Value = token.Value

		if token.Localize != nil {
			localTokens = append(localTokens, token)
		} else {
			globalTokens = append(globalTokens, token)
		}
	}

	_ = filepath.WalkDir(scaf.Path, func(path string, info os.DirEntry, err error) error {
		if info.Name() == configFileName {
			return nil
		}

		if path == scaf.Path {
			return nil
		}

		relativePath := strings.TrimPrefix(path, scaf.Path)

		relativePath = scaf.replaceTokens(relativePath, path, globalTokens)
		relativePath = scaf.replaceTokens(relativePath, path, localTokens)

		makeDestination := destination + relativePath

		if info.IsDir() {
			err = os.MkdirAll(makeDestination, os.ModePerm)
		} else {
			contents, err := os.ReadFile(path)
			if err != nil {
				panic(err)
			}

			stringcontents := string(contents)
			stringcontents = scaf.replaceTokens(stringcontents, path, globalTokens)
			stringcontents = scaf.replaceTokens(stringcontents, path, localTokens)

			err = os.WriteFile(makeDestination, []byte(stringcontents), 0644)
			if err != nil {
				panic(err)
			}
		}

		scaf.onMakeFunc(makeDestination)

		return nil
	})

	return nil
}

func (scaf *Scaffold) replaceTokens(subject string, path string, tokens []Token) string {
	for _, token := range tokens {
		if token.Localize != nil {
			for _, tokenPath := range token.Localize {
				if strings.HasPrefix(path, scaf.Path+"/"+tokenPath) {
					subject = strings.ReplaceAll(subject, token.Name, token.Value)
				}
			}
		} else {
			subject = strings.ReplaceAll(subject, token.Name, token.Value)
		}
	}

	return subject
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
