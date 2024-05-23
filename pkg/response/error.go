package response

type ErrorResponse struct {
	Code    int    `json:"code"`
	Err     string `json:"error,omitempty"`
	Trace   error  `json:"trace,omitempty"`
	Message string `json:"message"`
}

func (e *ErrorResponse) Error() string {
	return e.Message
}
