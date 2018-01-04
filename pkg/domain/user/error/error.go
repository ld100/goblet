package error

import "errors"

var (
	INTERNAL_SERVER_ERROR = errors.New("Internal Server Error")
	NOT_FOUND_ERROR       = errors.New("Your requested Item is not found")
	CONFLICT_ERROR        = errors.New("Your Item already exist")
	WRONG_PASSWORD        = errors.New("Provided password is wrong")
)
