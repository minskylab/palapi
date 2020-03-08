package static

import "github.com/pkg/errors"

var ErrInvalidStaticFileType = errors.New("type file not supported. this provider only supports: json, yaml, xml, csv ")
