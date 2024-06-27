package email

import "errors"

var ErrInvalidCode = errors.New("invalid code")
var ErrSendingVerificationEmailFaild = errors.New("sending verification email faild")
var ErrCheckingCodeFaild = errors.New("checking code faild") 
