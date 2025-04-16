package lib

import (
	"errors"
	"github.com/mattn/go-sqlite3"
)

func IsUniqueViolation(err error) bool {
	var sqliteErr sqlite3.Error
	if ok := asSqliteError(err, &sqliteErr); ok {
		return errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) &&
			errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique)
	}

	return false
}

func asSqliteError(err error, target *sqlite3.Error) bool {
	var e sqlite3.Error
	if errors.As(err, &e) {
		*target = e
		return true
	}

	return false
}
