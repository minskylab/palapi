package palapi

import "errors"

// import "github.com/pkg/errors"

var ErrFeatureNotExist = errors.New("feature not exists")
var ErrWordNotFound = errors.New("word not found at any provider")
