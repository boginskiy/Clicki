package model

type URLFile struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewURLFile(origin, short string) *URLFile {
	return &URLFile{
		ShortURL:    short,
		OriginalURL: origin,
	}
}
