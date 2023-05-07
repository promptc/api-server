package api

import "errors"

var ErrNotFound = errors.New("Not Found")
var ErrBindingError = errors.New("BINDING_ERROR")

type Request struct {
	Input map[string]string `json:"input"`
}

type VarInfo struct {
	Type       string `json:"type"`
	Constraint any    `json:"constraint"`
}

type VarInfoResponse map[string]VarInfo
