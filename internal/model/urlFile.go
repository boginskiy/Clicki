package model

type URLFile struct {
	UUID          int    `json:"uuid"`
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
	OriginalURL   string `json:"original_url"`
}

func NewURLFile(origin, short, correID string) *URLFile {
	return &URLFile{
		ShortURL:      short,
		CorrelationID: correID,
		OriginalURL:   origin,
	}
}
