package model

type (
	BaseLink struct {
		URL string `json:"url"`
	}

	ExtraLink struct {
		*BaseLink `json:"-"`
		Result    string `json:"result"`
	}
)

func NewBaseLink() *BaseLink {
	return &BaseLink{}
}

func NewExtraLink(base *BaseLink, result string) *ExtraLink {
	return &ExtraLink{
		BaseLink: base,
		Result:   result,
	}
}
