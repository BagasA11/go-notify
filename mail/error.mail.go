package mail

import "errors"

var ErrRecipient = errors.New("unreachable receiver")
var ErrSender = errors.New("unreachable Sender")
var ErrAuth = errors.New("Auth error")
