package dto

type (
	PayloadLogin struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	ResponseJWT struct {
		TokenJwt  string         `json:"token_jwt"`
		ExpiredAt string         `json:"expired_at"`
		DataUser  *DataUserLogin `json:"data_user"`
	}

	DataUserLogin struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}
)
