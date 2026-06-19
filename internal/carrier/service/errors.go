package service

import "errors"

// ErrCarrierNameExists indicates that a carrier name is already in use.
var ErrCarrierNameExists = errors.New("carrier name already exists")
