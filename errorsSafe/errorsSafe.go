package errorsSafe

import (
	"log"
	"net/http"
)

var (
	ErrNotFound   = &UserSafeError{status: http.StatusNotFound, msg: "not found"}
	ErrServer     = &UserSafeError{status: http.StatusInternalServerError, msg: "internal server error"}
	ErrNotAllowed = &UserSafeError{status: http.StatusMethodNotAllowed, msg: "method not allowed"}
)

type HttpError interface {
	HttpError() (int, string)
}

type UserSafeError struct {
	status int
	msg    string
}

func (e UserSafeError) Error() string {
	return e.msg
}

func (e UserSafeError) HttpError() (int, string) {
	return e.status, e.msg
}

// matches internal server errors with user safe errors
type wrappedOriginalError struct {
	error
	UserSafeError *UserSafeError
}

func (e wrappedOriginalError) Is(err error) bool {
	return e.UserSafeError == err
}

func (e wrappedOriginalError) HttpError() (int, string) {
	return e.UserSafeError.HttpError()
}

// logs original error and returns a user safe error
func WrapError(err error, safeErr *UserSafeError) error {
	log.Print(err)
	return wrappedOriginalError{error: err, UserSafeError: safeErr}
}
