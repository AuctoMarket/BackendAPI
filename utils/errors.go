package utils

/*
Error handler for response errors
*/
type ErrorHandler struct {
	Message string
	Code    int
}

func (m *ErrorHandler) Error() string {
	return m.Message
}

func (m *ErrorHandler) ErrorCode() int {
	return m.Code
}

/*
Creates 400 Bad request error
*/
func BadRequestError(msg string) *ErrorHandler {
	return &ErrorHandler{Message: msg, Code: 400}
}

/*
Creates 500 Internal Server Error
*/
func InternalServerError() *ErrorHandler {
	return &ErrorHandler{Message: "Something went wrong", Code: 500}
}

/*
Creates 401 Unautorized Error
*/
func UnauthorizedError(msg string) *ErrorHandler {
	return &ErrorHandler{Message: msg, Code: 401}
}
