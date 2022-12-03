package errorh

import "fmt"

type RaiseError struct {
	ErrorMessage string
	ErrorCode    int
}

func (r *RaiseError) Error() string {
	msg := fmt.Sprintf("status-code:%d-%s", r.ErrorCode, r.ErrorMessage)
	return msg
}
