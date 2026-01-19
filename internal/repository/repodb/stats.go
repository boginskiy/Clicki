package repodb

import "context"

type DBStatistician struct {
	Repo *RepoDB
}

func NewDBStatistician(repo *RepoDB) *DBStatistician {
	return &DBStatistician{
		Repo: repo,
	}
}

func (s *DBStatistician) ReadQuantityUsers(ctx context.Context) int {
	row := s.Repo.Store.QueryRowContext(ctx,
		`SELECT COUNT(DISTINCT user_id)
		 FROM urls;`)

	var tmpQuantity int

	err := row.Scan(&tmpQuantity)
	if err != nil {
		s.Repo.Logg.RaiseError(err, "RepoDB>ReadQuantityUsers", nil)
	}
	return tmpQuantity
}

func (s *DBStatistician) ReadQuantityShortURLs(ctx context.Context) int {
	row := s.Repo.Store.QueryRowContext(ctx,
		`SELECT COUNT(DISTINCT short_url)
	 	 FROM urls;`)

	var tmpQuantity int

	err := row.Scan(&tmpQuantity)
	if err != nil {
		s.Repo.Logg.RaiseError(err, "RepoDB>ReadQuantityShortURLs", nil)
	}
	return tmpQuantity
}
