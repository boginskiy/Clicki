package repomap

import "context"

type MapStatistician struct {
	Repo *RepoMap
}

func NewMapStatistician(repo *RepoMap) *MapStatistician {
	return &MapStatistician{
		Repo: repo,
	}
}

func (s *MapStatistician) ReadQuantityUsers(_ context.Context) int {
	tmpSet := make(map[int]struct{}, 10)
	for _, v := range s.Repo.Store {
		tmpSet[v.UserID] = struct{}{}
	}
	return len(tmpSet)
}

func (s *MapStatistician) ReadQuantityShortURLs(_ context.Context) int {
	return len(s.Repo.Store)
}
