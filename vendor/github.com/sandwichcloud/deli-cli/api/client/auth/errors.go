package auth

import "errors"

var ErrOTPRequired = errors.New("OTP Token is required.")
var ErrOTPInvalid = errors.New("OTP Token is invalid.")
