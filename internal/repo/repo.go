package repo

import (
	"github.com/jackc/pgx/v5"
)

type repoImpl struct {
	*pgx.Conn
}
