package database

import (
	"errors"
)

var ErrDuplicateKey = errors.New("duplicate key error")

var DBError = errors.New("mongodb error")

type RecordNotFoundError struct{}

func (e *RecordNotFoundError) Error() string {
	return "record not found"
}
