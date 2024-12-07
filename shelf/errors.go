package shelf

import "errors"

var (
	ErrRefHeadNotFound = errors.New("reference head not found")
	ErrDiaryNotFound   = errors.New("diary not found")
)
