package eam

import "errors"

var (
	ErrAPIServiceLogin      = errors.New("cannot login to api service, check user's name and pass")
	ErrUnexpectedDataFormat = errors.New("wrong expected data format, json or text required")
)
