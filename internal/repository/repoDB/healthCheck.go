package repoDB

import "context"

type RepoDBHealthCheck struct {
	Repo *RepoDB
}

func NewRepoDBHealthCheck(repoDB *RepoDB) *RepoDBHealthCheck {
	return &RepoDBHealthCheck{
		Repo: repoDB,
	}
}

// PingStore - implements interface HealthCheckRepo.
func (r *RepoDBHealthCheck) PingStore(ctx context.Context) (bool, error) {
	return r.Repo.DB.CheckOpen()
}
