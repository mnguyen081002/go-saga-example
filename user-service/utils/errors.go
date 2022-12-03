package utils

import (
	"errors"
	"net/http"
)

func GetStatusCode(err error) (int, bool) {
	if v, ok := MapErrorStatusCode[err.Error()]; !ok {
		return http.StatusInternalServerError, false
	} else {
		return v, true
	}
}

var (
	ErrInternalServerError = errors.New("Internal Server Error")
	ErrorBadParamInput     = "Bad param input %v"
)

const (
	ItemNotFound = "Item not found"
)

var MapErrorStatusCode = map[string]int{
	ItemNotFound: http.StatusNotFound,
}
