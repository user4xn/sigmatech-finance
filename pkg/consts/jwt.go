package consts

import "time"

type (
	SessionStatus int
	UserType      string
)

const (
	SessionActive  SessionStatus = 0
	SessionRevoked SessionStatus = 1

	TokenDurationRelease = time.Hour * 24 * 7
	TokenDurationDev     = time.Hour * 24 * 30

	UserTypeAdmin    UserType = "admin"
	UserTypeConsumer UserType = "consumer"
)
