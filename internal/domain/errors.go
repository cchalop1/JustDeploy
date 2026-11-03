package domain

import "errors"

var ErrNotFound = errors.New("service not found")
var ErrConflict = errors.New("another operation is in progress; retry later")
var ErrAlreadyConfigured = errors.New("admin account already exists")
var ErrUnsupported = errors.New("this legacy operation is no longer supported")
