package repo

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"movie-lib/internal/model"
)

const getUserRoleQuery = `
		SELECT "role" FROM "users"
		WHERE "id" = $1;`

func (r *repoImpl) GetUserRole(ctx context.Context, id uint64) (model.Role, error) {
	row := r.QueryRow(ctx, getUserRoleQuery, id)
	var role model.Role
	if err := row.Scan(&role); errors.Is(err, pgx.ErrNoRows) {
		return "", model.ErrUserNotExists
	} else if err != nil {
		return "", errors.Join(model.ErrDatabaseError, err)
	}
	return role, nil
}
