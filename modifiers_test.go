package scaffold

import "testing"

func TestModifierLower(t *testing.T) {
	expected := "lower stringtest"
	input := "Lower StringTest"

	actual := ModifierLower(input)

	if actual != expected {
		t.Errorf("Unexpected result from modifier. Expected: %v, Got %v", expected, actual)
	}
}

func TestModifierUpper(t *testing.T) {
	expected := "UPPER STRING"
	input := "upper string"

	actual := ModifierUpper(input)

	if actual != expected {
		t.Errorf("Unexpected result from modifier. Expected: %v, Got %v", expected, actual)
	}
}

func TestModifierSlug(t *testing.T) {
	expected := "some-crazy-string"
	input := "SomeCrazy String!"

	actual := ModifierSlug(input)

	if actual != expected {
		t.Errorf("Unexpected result from modifier. Expected: %v, Got %v", expected, actual)
	}
}

func TestModifierSnake(t *testing.T) {
	expected := "snake_case"
	input := "SnakeCase"

	actual := ModifierSnake(input)

	if actual != expected {
		t.Errorf("Unexpected result from modifier. Expected: %v, Got %v", expected, actual)
	}

	input = "Snake-Case"

	actual = ModifierSnake(input)

	if actual != expected {
		t.Errorf("Unexpected result from modifier. Expected: %v, Got %v", expected, actual)
	}
}

func TestModifierPascal(t *testing.T) {
	expected := "PascalCase"
	input := "pascal-case"

	actual := ModifierPascal(input)

	if actual != expected {
		t.Errorf("Unexpected result from modifier. Expected: %v, Got %v", expected, actual)
	}

	input = "pascal_case"

	actual = ModifierPascal(input)

	if actual != expected {
		t.Errorf("Unexpected result from modifier. Expected: %v, Got %v", expected, actual)
	}
}

func TestModifierCamel(t *testing.T) {
	expected := "camelCase"
	input := "camel-case"

	actual := ModifierCamel(input)

	if actual != expected {
		t.Errorf("Unexpected result from modifier. Expected: %v, Got %v", expected, actual)
	}

	input = "Camel_case"

	actual = ModifierCamel(input)

	if actual != expected {
		t.Errorf("Unexpected result from modifier. Expected: %v, Got %v", expected, actual)
	}
}

func TestModifierTitle(t *testing.T) {
	expected := "A Title Test"
	input := "ATitleTest"

	actual := ModifierTitle(input)

	if actual != expected {
		t.Errorf("Unexpected result from modifier. Expected: %v, Got %v", expected, actual)
	}

	input = "A_Title-Test"

	actual = ModifierTitle(input)

	if actual != expected {
		t.Errorf("Unexpected result from modifier. Expected: %v, Got %v", expected, actual)
	}
}

func TestModifierPlural(t *testing.T) {
	expected := "Countries"
	input := "Country"

	actual := ModifierPlural(input)

	if actual != expected {
		t.Errorf("Unexpected result from modifier. Expected: %v, Got %v", expected, actual)
	}

	expected = "States"
	input = "State"

	actual = ModifierPlural(input)

	if actual != expected {
		t.Errorf("Unexpected result from modifier. Expected: %v, Got %v", expected, actual)
	}
}
