package push

import "errors"

var (
	deviceTokenNotFound = errors.New("no device token found for account")
)
