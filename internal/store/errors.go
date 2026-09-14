package store

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var ErrLoginTaken = errors.New("[store] login already taken")

var ErrNotFound = errors.New("[store] not found")

var ErrBadFilterKey = errors.New("[store] invalid filter key")

const pgUniqueViolation = "23505"

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation
}
