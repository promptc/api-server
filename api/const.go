package api

import "errors"

var ErrNotFound = errors.New("Not Found")
var ErrBindingError = errors.New("BINDING_ERROR")

type Request struct {
	Input map[string]string `json:"input"`
}
