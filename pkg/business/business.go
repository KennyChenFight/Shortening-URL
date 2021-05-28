package business

func NewError(businessCode int, httpStatusCode int, message string, reason error) *Error {
	return &Error{BusinessCode: businessCode, HTTPStatusCode: httpStatusCode, Message: message, Reason: reason}
}

type Error struct {
	BusinessCode     int               `json:"code"`
	HTTPStatusCode   int               `json:"-"`
	Message          string            `json:"message"`
	ValidationErrors map[string]string `json:"validationErrors,omitempty"`
	Reason           error             `json:"-"`
}

func (b *Error) Error() string {
	return b.Reason.Error()
}
