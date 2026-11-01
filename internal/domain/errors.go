package domain

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user with that id already exists")
var ErrEmptyName = errors.New("username is empty or has not enough length")
var ErrInvalidData = errors.New("invalid data")
