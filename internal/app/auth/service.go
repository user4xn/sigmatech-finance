package auth

import (
	"clean-arch/database"
	"clean-arch/internal/dto"
	"clean-arch/internal/factory"
	"clean-arch/internal/model"
	"clean-arch/internal/repository"
	"clean-arch/pkg/consts"
	"clean-arch/pkg/util"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type service struct {
	UserRepository  repository.User
	RedisRepository repository.Redis
}

type Service interface {
	LoginAttempt(ctx context.Context, reqHandler dto.PayloadLogin) (dto.ResponseJWT, error)
	Logout(ctx context.Context, bearer string) error
}

func NewService(f *factory.Factory) Service {
	return &service{
		UserRepository:  f.UserRepository,
		RedisRepository: f.RedisRepository,
	}
}

func (s *service) Logout(ctx context.Context, bearer string) error {
	tx := database.BeginTx(ctx, factory.NewFactory().InitDB)
	if err := tx.Error; err != nil {
		return err
	}

	err := s.UserRepository.RevokeSession(tx, bearer)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (s *service) GenerateToken(secretKey []byte, userID string, email string) (string, *time.Time, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return "", nil, err
	}

	jwtMode := util.GetEnv("JWT_MODE", "fallback")
	nonUnixTime := time.Now().In(loc).Add(consts.TokenDurationDev)
	expiredTime := nonUnixTime.Unix()

	if jwtMode == "release" {
		nonUnixTime := time.Now().In(loc).Add(consts.TokenDurationRelease)
		expiredTime := nonUnixTime.Unix()

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": userID,
			"email":   email,
			"exp":     expiredTime,
		})

		tokenString, err := token.SignedString(secretKey)
		if err != nil {
			return "", nil, err
		}

		return tokenString, &nonUnixTime, nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     expiredTime,
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", nil, err
	}

	return tokenString, &nonUnixTime, nil
}

func (s *service) LoginAttempt(ctx context.Context, reqHandler dto.PayloadLogin) (dto.ResponseJWT, error) {
	var (
		res dto.ResponseJWT
	)

	user, err := s.UserRepository.FindOne(ctx, "id, email, password", "email = ?", reqHandler.Email)
	if err != nil {
		return res, consts.UserNotFound
	}

	err = util.ComparePasswords(user.Password, reqHandler.Password)
	if err != nil {
		return res, consts.InvalidPassword
	}

	secretKey := []byte(util.GetEnv("APP_SECRET_KEY", "fallback"))
	jwt, exp, err := s.GenerateToken(secretKey, strconv.Itoa(user.ID), user.Email)
	if err != nil {
		return res, consts.ErrorGenerateJwt
	}

	if jwt == "" {
		return res, consts.EmptyGenerateJwt
	}

	dataUser := dto.DataUserLogin{
		ID:    user.ID,
		Email: user.Email,
	}

	sessionModel := model.UserSession{
		UserID:    user.ID,
		JWTToken:  jwt,
		ExpiresAt: *exp,
	}

	tx := database.BeginTx(ctx, factory.NewFactory().InitDB)
	if err := tx.Error; err != nil {
		return res, err
	}

	err = s.UserRepository.CreateSession(tx, sessionModel)
	if err != nil {
		tx.Rollback()
		return res, consts.ErrorGenerateJwt
	}
	tx.Commit()

	cacheKey := fmt.Sprintf("user_session-%d", user.ID)
	_ = s.RedisRepository.Del(ctx, cacheKey)

	res = dto.ResponseJWT{
		TokenJwt:  jwt,
		ExpiredAt: exp.Format(consts.TimeFormatDateTime),
		DataUser:  &dataUser,
	}

	return res, nil
}
