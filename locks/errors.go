package locks

import "errors"

var ErrIsLocked = errors.New(
	"File is locked, another process is writing to it, wait and try again later.",
)
