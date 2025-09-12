package model

type (
	URLJson struct {
		URL string `json:"url"`
	}

	ResultJson struct {
		*URLJson `json:"-"`
		Result   string `json:"result"`
	}
)

func NewURLJson() *URLJson {
	return &URLJson{}
}

func NewResultJson(url *URLJson, result string) *ResultJson {
	return &ResultJson{
		URLJson: url,
		Result:  result,
	}
}
