package httpErrors

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"fmt"
)

const (
	BadRequestError                      = "Bad Request"
	WrongCredentialsError                = "Wrong Credentials"
	EmailAlreadyExistsError              = "User with given email already exists"
	ProfileError                         = "Consultants must have a profile"
	CertificationError                   = "Consultants must have a certification"
	CertificationAndProfileNotExistError = "Consultants must have a profile and a certification"
	NotFoundError                        = "Not Found"
	UnauthorizedError                    = "Unauthorized"
	ForbiddenError                       = "Forbidden"
	BadQueryParamsError                  = "Invalid query params"
)

var (
	ErrBadRequest          = errors.New("Bad Request")
	ErrWrongCredentials    = errors.New("Wrong Credentials")
	ErrNotFound            = errors.New("Not Found")
	ErrUnauthorized        = errors.New("Unauthorized")
	ErrForbidden           = errors.New("Forbidden")
	ErrPermissionDenied    = errors.New("Permission Denied")
	ErrExpiredCSRF         = errors.New("Expired CSRF token")
	ErrWrongCSRFToken      = errors.New("Wrong CSRF token")
	ErrCSRFNotPresented    = errors.New("CSRF not presented")
	ErrNotRequiredFields   = errors.New("No such required fields")
	ErrBadQueryParams      = errors.New("Invalid query params")
	ErrInternalServer      = errors.New("Internal Server Error")
	ErrRequestTimeout      = errors.New("Request Timeout")
	ErrExistsEmail         = errors.New("User with given email already exists")
	ErrInvalidJWTToken     = errors.New("Invalid JWT Token")
	ErrInvalidJWTClaims    = errors.New("Invalid JWT Claims")
	ErrRequestTimeoutError = errors.New("Request Timeout error")
	ErrExistsEmailError    = errors.New("Exist Email Error")
)

// Rest error interface
type RestErr interface {
	Status() int
	Error() string
	Causes() interface{}
}

// REst error struct

type RestError struct {
	ErrStatus int         `json:"status,omitempty"`
	ErrError  string      `json:"error,omitempty"`
	ErrCauses interface{} `json:"-"`
}

// Error
func (e RestError) Error() string {
	return fmt.Sprintf("status: %d - errors: %s - causes : %v", e.ErrStatus, e.ErrError, e.ErrCauses)
}

// Error status
func (e RestError) Status() int {
	return e.ErrStatus
}

// RestError Causes
func (e RestError) Causes() interface{} {
	return e.ErrCauses
}

// New rest error
func NewRestError(status int, err string, causes interface{}) RestErr {
	return RestError{
		ErrStatus: status,
		ErrError:  err,
		ErrCauses: causes,
	}
}

// New rest error with message
func NewRestErrorWithMessage(status int, err string, causes interface{}) RestErr {
	return RestError{
		ErrStatus: status,
		ErrError:  err,
		ErrCauses: causes,
	}
}

// New rest error form bytes
func NewRestErrorFromBytes(bytes []byte) (RestErr, error) {
	var apiErr RestError

	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, errors.New("InvalidJson")
	}
	return &apiErr, nil
}

// New Bad Request Error
func NewBadRequestError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusBadRequest,
		ErrError:  ErrBadRequest.Error(),
		ErrCauses: causes,
	}
}

// New Not Found Error
func NewNotFoundError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusNotFound,
		ErrError:  ErrNotFound.Error(),
		ErrCauses: causes,
	}
}

// New Unauthorized Error
func NewUnauthorizedError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusUnauthorized,
		ErrError:  ErrUnauthorized.Error(),
		ErrCauses: causes,
	}
}

// New Forbidden Error
func NewForbiddenError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusForbidden,
		ErrError:  ErrForbidden.Error(),
		ErrCauses: causes,
	}
}

// New Internal Server Error
func NewInternalServerError(causes interface{}) RestErr {
	result := RestError{
		ErrStatus: http.StatusInternalServerError,
		ErrError:  ErrInternalServer.Error(),
		ErrCauses: causes,
	}
	return result
}

// Parser of error string messages returns RestError
func ParseErrors(err error) RestErr {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return NewRestError(http.StatusNotFound, ErrNotFound.Error(), err)
	case errors.Is(err, context.DeadlineExceeded):
		return NewRestError(http.StatusRequestTimeout, ErrRequestTimeoutError.Error(), err)
	case strings.Contains(err.Error(), "SQLSTATE"):
		return parseSqlErrors(err)
	case strings.Contains(err.Error(), "Field validation"):
		return parseValidatorError(err)
	case strings.Contains(err.Error(), "Unmarshal"):
		return NewRestError(http.StatusBadRequest, ErrBadRequest.Error(), err)
	case strings.Contains(err.Error(), "UUID"):
		return NewRestError(http.StatusBadRequest, err.Error(), err)
	case strings.Contains(strings.ToLower(err.Error()), "cookie"):
		return NewRestError(http.StatusUnauthorized, ErrUnauthorized.Error(), err)
	case strings.Contains(strings.ToLower(err.Error()), "token"):
		return NewRestError(http.StatusUnauthorized, ErrUnauthorized.Error(), err)
	case strings.Contains(strings.ToLower(err.Error()), "bcrypt"):
		return NewRestError(http.StatusBadRequest, ErrBadRequest.Error(), err)
	default:
		if restErr, ok := err.(RestErr); ok {
			return restErr
		}
		return NewInternalServerError(err)
	}
}
func parseSqlErrors(err error) RestErr {
	if strings.Contains(err.Error(), "23505") {
		return NewRestError(http.StatusBadRequest, ErrExistsEmailError.Error(), err)
	}

	return NewRestError(http.StatusBadRequest, ErrBadRequest.Error(), err)
}

func parseValidatorError(err error) RestErr {
	if strings.Contains(err.Error(), "Password") {
		return NewRestError(http.StatusBadRequest, "Invalid password, min length 6", err)
	}

	if strings.Contains(err.Error(), "Email") {
		return NewRestError(http.StatusBadRequest, "Invalid email", err)
	}

	return NewRestError(http.StatusBadRequest, ErrBadRequest.Error(), err)
}

// Error response
func ErrorResponse(err error) (int, interface{}) {
	return ParseErrors(err).Status(), ParseErrors(err)
}
