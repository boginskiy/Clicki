package model

import "time"

// Struct - struct for one of URL.
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

// Struct - struct for set of URL.
type (
	ReqURLSet struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	ResURLSet struct {
		CorrelationID string    `json:"correlation_id"` // CorrelationID is - Unic id.
		OriginalURL   string    `json:"-"`              // OriginalURL is - URL for shorting.
		ShortURL      string    `json:"short_url"`      // ShortURL is - short link.
		CreatedAt     time.Time `json:"-"`              // CreatedAt is - time of create record.
		UserID        int       `json:"-"`              // UserID is - ID user.
	}

	ResUserURLSet struct {
		OriginalURL string `json:"original_url"`
		ShortURL    string `json:"short_url"`
	}
)

func NewResURLSet(correlationID, origin, short string, userID int) ResURLSet {
	return ResURLSet{
		CorrelationID: correlationID,
		OriginalURL:   origin,
		ShortURL:      short,
		CreatedAt:     time.Now(),
		UserID:        userID,
	}
}

func NewResUserURLSet(origin, short string) ResUserURLSet {
	return ResUserURLSet{
		OriginalURL: origin,
		ShortURL:    short,
	}
}
