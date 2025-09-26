package model

import "time"

// Struct for one of URL
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

// Struct for set of URL/UserURL
type (
	ReqURLSet struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	ResURLSet struct {
		CorrelationID string    `json:"correlation_id"` // CorrelationID is - Уникальный строковый идентификатор
		OriginalURL   string    `json:"-"`              // OriginalURL is - URL для сокращения
		ShortURL      string    `json:"short_url"`      // ShortURL is - Сокращённая ссылка
		CreatedAt     time.Time `json:"-"`              // CreatedAt is - Время создания записи
		UserID        int       `json:"-"`              // UserID is - Идентификатор пользователя
	}

	ResUserURLSet struct {
		OriginalURL string `json:"original_url"`
		ShortURL    string `json:"short_url"`
	}
)

func NewResURLSet(correlationID, origin, short string, id int) ResURLSet {
	return ResURLSet{
		CorrelationID: correlationID,
		OriginalURL:   origin,
		ShortURL:      short,
		CreatedAt:     time.Now(),
		UserID:        id,
	}
}

func NewResUserURLSet(origin, short string) ResUserURLSet {
	return ResUserURLSet{
		OriginalURL: origin,
		ShortURL:    short,
	}
}
