package repofile

import "context"

type FileStatistician struct {
	Repo *RepoFile
}

func NewFileStatistician(repo *RepoFile) *FileStatistician {
	return &FileStatistician{
		Repo: repo,
	}
}

func (s *FileStatistician) ReadQuantityUsers(_ context.Context) int {
	tmpSet := make(map[int]struct{}, 10)
	for _, v := range s.Repo.tmpStore {
		tmpSet[v.UserID] = struct{}{}
	}
	return len(tmpSet)
}

func (s *FileStatistician) ReadQuantityShortURLs(_ context.Context) int {
	return len(s.Repo.tmpStore)
}
