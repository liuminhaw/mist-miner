package mmerr

const (
	MineCmdType      = "mine"
	CatFileCmdType   = "cat-file"
	LogCmdType       = "log"
	LogReloadCmdType = "log reload"
)

type ArgsError struct {
	CmdType string
	Msg     string
}

func NewArgsError(cmdType, message string) ArgsError {
	return ArgsError{CmdType: cmdType, Msg: message}
}

func (e ArgsError) Error() string {
	return e.Msg
}
