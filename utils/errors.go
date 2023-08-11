package utils

import "fmt"

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
	if msg == "" {
		return &ErrorHandler{Message: "Bad request", Code: 400}
	}

	return &ErrorHandler{Message: msg, Code: 400}
}

/*
Creates 500 Internal Server Error
*/
func InternalServerError(err error) *ErrorHandler {
	if err != nil {
		return &ErrorHandler{Message: fmt.Sprintf("Something went wrong:%s", err.Error()), Code: 500}
	}

	return &ErrorHandler{Message: fmt.Sprintf("Something went wrong"), Code: 500}
}

/*
Creates 401 Unautorized Error
*/
func UnauthorizedError(msg string) *ErrorHandler {
	if msg == "" {
		return &ErrorHandler{Message: "Unauthorized user", Code: 401}
	}
	return &ErrorHandler{Message: msg, Code: 401}
}

/*
Creates 404 Not Found Error
*/
func NotFoundError(msg string) *ErrorHandler {
	if msg == "" {
		return &ErrorHandler{Message: "Resource not found", Code: 404}
	}
	return &ErrorHandler{Message: msg, Code: 404}
}
