package camelToken

type PascalToken struct {
	Url string
}

func (camelToken *PascalToken) NewPascalToken() *PascalToken {
	camelToken.Url = "slug-token"

	return &PascalToken{}
}
