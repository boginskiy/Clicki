package repodb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/boginskiy/Clicki/internal/repository/utils"
)

type RepoDBMarker struct {
	Repo *RepoDB
}

func NewRepoDBMarker(repoDB *RepoDB) *RepoDBMarker {
	return &RepoDBMarker{
		Repo: repoDB,
	}
}

// MarkRecords - implements interface MarkerRepo.
func (r *RepoDBMarker) MarkRecords(ctx context.Context, msg any) error {
	messages, ok := msg.([]utils.DelMessage)
	if !ok {
		return errors.New(`{"error": "message for delete is bad"}`)
	}

	values := make([]string, 0, 10)
	args := make([]any, 0, 10)
	c := 1

	for _, mess := range messages {

		for _, correlID := range mess.ListCorrelID {
			values = append(values, fmt.Sprintf("($%d,$%d)", c, c+1))
			args = append(args, mess.UserID, correlID)
			c += 2
		}
	}

	query := fmt.Sprintf(`UPDATE urls
                          SET deleted_flag = TRUE
                          WHERE (user_id, correlation_id) IN (%s)`, strings.Join(values, ","))

	_, err := r.Repo.Store.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}
