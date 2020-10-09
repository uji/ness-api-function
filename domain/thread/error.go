package thread

import "errors"

var (
	ErrorLastEvaluatedIDCanNotBeBlank error = errors.New("lastEvaluatedID can not be blank")
	ErrorTitleIsRequired              error = errors.New("title is required")
	ErrorTimeFormatInValid            error = errors.New("time format invalid")
)
