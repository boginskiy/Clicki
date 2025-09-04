package db

type StoreModel struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewStoreModel(id int, shortURL, originURL string) *StoreModel {
	return &StoreModel{
		UUID:        id,
		ShortURL:    shortURL,
		OriginalURL: originURL,
	}
}
