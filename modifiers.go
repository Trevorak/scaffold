package scaffold

import (
	"github.com/gertd/go-pluralize"
	"strings"
)

type modifierMap map[string][]func(string) string

func (modMap *modifierMap) Add(token string, modifier func(string) string) *modifierMap {
	(*modMap)[token] = append((*modMap)[token], modifier)

	return modMap
}

func ModifierLower(subject string) string {
	return strings.ToLower(subject)
}

func ModifierUpper(subject string) string {
	return strings.ToUpper(subject)
}

func ModifierSlug(subject string) string {
	var slug []uint8
	for i := range subject {
		// if it's a number, add it to slug
		if subject[i] >= 48 && subject[i] <= 57 {
			slug = append(slug, subject[i])
		} else if subject[i] >= 97 && subject[i] <= 122 {
			slug = append(slug, subject[i])
		} else
		// if the character is an upper-case letter, make it lower-case.
		if subject[i] >= 65 && subject[i] <= 90 {

			// If this isn't the first letter, and there's not already a dash preceding, put a dash before the letter
			slugLength := len(slug)
			if i != 0 && slug[slugLength-1] != 45 {
				slug = append(slug, 45)
			}

			lower := subject[i] + 32
			slug = append(slug, lower)
		} else {
			// any other character, add a dash, if there's not already one preceding
			slugLength := len(slug)
			if i != 0 && slug[slugLength-1] != 45 {
				slug = append(slug, 45)
			}
		}
	}

	// if the last char in the slug is a dash, remove it.
	slugLength := len(slug)
	if slug[slugLength-1] == 45 {
		slug = slug[:slugLength-1]
	}

	return string(slug)
}

func ModifierSnake(subject string) string {
	var modified []uint8
	for i := range subject {
		// if it's a number, add it to modified
		if subject[i] >= 48 && subject[i] <= 57 {
			modified = append(modified, subject[i])
		} else if subject[i] >= 97 && subject[i] <= 122 {
			modified = append(modified, subject[i])
		} else
		// if the character is an upper-case letter, make it lower-case.
		if subject[i] >= 65 && subject[i] <= 90 {

			// If this isn't the first letter, and there's not already an underscore preceding, put it before the letter
			modifiedLength := len(modified)
			if i != 0 && modified[modifiedLength-1] != 95 {
				modified = append(modified, 95)
			}

			lower := subject[i] + 32
			modified = append(modified, lower)
		} else {
			// any other character, add an underscore, if there's not already one preceding
			modifiedLength := len(modified)
			if i != 0 && modified[modifiedLength-1] != 95 {
				modified = append(modified, 95)
			}
		}
	}

	// if the last char in the modified is an underscore, remove it.
	modifiedLength := len(modified)
	if modified[modifiedLength-1] == 95 {
		modified = modified[:modifiedLength-1]
	}

	return string(modified)
}

func ModifierPascal(subject string) string {
	var modified []uint8
	for i := range subject {
		// if it's a number, add it to modified
		if subject[i] >= 48 && subject[i] <= 57 {
			modified = append(modified, subject[i])
		} else if subject[i] >= 97 && subject[i] <= 122 {
			// if it's lower case
			// if it's the first letter, capitalize
			// if it's not the first, and the preceding character is
			// a space, hyphen, underscore, capitalize
			if i == 0 {
				upper := subject[i] - 32
				modified = append(modified, upper)
			} else {
				if subject[i-1] == 95 || subject[i-1] == 45 || subject[i-1] == 32 {
					upper := subject[i] - 32
					modified = append(modified, upper)
				} else {
					modified = append(modified, subject[i])
				}
			}
		} else
		// if the character is an upper-case letter
		if subject[i] >= 65 && subject[i] <= 90 {
			modified = append(modified, subject[i])
		}
	}

	return string(modified)
}

func ModifierCamel(subject string) string {
	var modified []uint8
	for i := range subject {
		// if it's a number, add it to modified
		if subject[i] >= 48 && subject[i] <= 57 {
			modified = append(modified, subject[i])
		} else if subject[i] >= 97 && subject[i] <= 122 {
			// if it's lower case
			// if it's not the first, and the preceding character is
			// a space, hyphen, underscore, capitalize
			if i == 0 {
				modified = append(modified, subject[i])
			} else {
				if subject[i-1] == 95 || subject[i-1] == 45 || subject[i-1] == 32 {
					upper := subject[i] - 32
					modified = append(modified, upper)
				} else {
					modified = append(modified, subject[i])
				}
			}
		} else
		// if the character is an upper-case letter
		if subject[i] >= 65 && subject[i] <= 90 {
			// if it's the first character, lower case
			if i == 0 {
				lower := subject[i] + 32
				modified = append(modified, lower)
			} else {
				modified = append(modified, subject[i])
			}
		}
	}

	return string(modified)
}

func ModifierTitle(subject string) string {
	var modified []uint8
	for i := range subject {
		// if it's a number, add it to modified
		if subject[i] >= 48 && subject[i] <= 57 {
			modified = append(modified, subject[i])
		} else if subject[i] >= 97 && subject[i] <= 122 {
			// if first char, make capital.
			if i == 0 {
				upper := subject[i] - 32

				modified = append(modified, upper)
			} else {
				// if preceding char was not a char, make capital, if it was not a space, add one

				modified = append(modified, subject[i])
			}
		} else
		// if char is capital
		if subject[i] >= 65 && subject[i] <= 90 {

			// if not the first char and preceding char was not a space, add space
			modifiedLength := len(modified)
			if i != 0 && modified[modifiedLength-1] != 32 {
				modified = append(modified, 32)
			}

			modified = append(modified, subject[i])
		}
	}

	return string(modified)
}

func ModifierPlural(subject string) string {
	client := pluralize.NewClient()

	return client.Plural(subject)
}

func ModifierSingular(subject string) string {
	client := pluralize.NewClient()

	return client.Singular(subject)
}
