package lib

import (
	"errors"
	"github.com/lib/pq"
)

const PGUniqueViolation = "23505"

func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == PGUniqueViolation
	}

	return false
}
