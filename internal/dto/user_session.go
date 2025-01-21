package dto

import (
	"clean-arch/pkg/consts"
	"time"
)

type (
	UserSession struct {
		UserID    int                  `json:"user_id"`
		JWTToken  string               `json:"jwt_token"`
		ExpiresAt string               `json:"expires_at"`
		Revoked   consts.SessionStatus `json:"revoked"`
	}

	JwtSession struct {
		ID         int       `json:"id"`
		ConsumerID int       `json:"consumer_id"`
		Email      string    `json:"email"`
		UserType   string    `json:"user_type"`
		CreatedAt  time.Time `json:"created_at"`
	}
)
