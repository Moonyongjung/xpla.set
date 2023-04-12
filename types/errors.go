package types

type XGoError struct {
	code uint64
	desc string
}

// Return error code and message generating on the anchor.
var (
	Success           = new(0, "success")
	ErrParse          = new(101, "error parsing data")
	ErrInvalidRequest = new(102, "invalid request")
	ErrInvalidMsgType = new(103, "invalid message type")
	ErrAlreadyExist   = new(104, "already exist")
	ErrNotFound       = new(105, "not found")
)

func new(code uint64, desc string) XGoError {
	var xErr XGoError
	xErr.code = code
	xErr.desc = desc

	return xErr
}

func (x XGoError) Code() uint64 {
	return x.code
}

func (x XGoError) Desc() string {
	return x.desc
}
