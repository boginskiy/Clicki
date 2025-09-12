package model

type (
	URLJson struct {
		URL string `json:"url"`
	}

	ResultJSON struct {
		*URLJson `json:"-"`
		Result   string `json:"result"`
	}
)

func NewURLJson() *URLJson {
	return &URLJson{}
}

func NewResultJSON(url *URLJson, result string) *ResultJSON {
	return &ResultJSON{
		URLJson: url,
		Result:  result,
	}
}
