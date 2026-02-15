package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	// Result of operation (OK, Created, Error)
	Status  string `json:"status" example:"created"`
	Message string `json:"message,omitempty" example:"success"`
}

const (
	StatusOK      = "OK"
	StatusCreated = "Created"
	StatusError   = "Error"
)

func OK(msg string) Response {
	return Response{
		Status:  StatusOK,
		Message: msg,
	}
}

func Created(msg string) Response {
	return Response{
		Status:  StatusCreated,
		Message: msg,
	}
}

func Error(msg string) Response {
	return Response{
		Status:  StatusError,
		Message: msg,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status:  StatusError,
		Message: strings.Join(errMsgs, ", "),
	}
}
