package users

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrInvalidID      = errors.New("invalid user ID")
	ErrPhoneNumberTaken  = errors.New("phone number already taken")
	ErrUsernameTaken  = errors.New("username already taken")
	ErrInvalidRequest = errors.New("invalid request")
	ErrInactiveUser   = errors.New("inactive user")
)