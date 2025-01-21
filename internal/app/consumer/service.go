package consumer

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
	"strings"
	"time"
)

type service struct {
	ConsumerRepository repository.Consumer
	RedisRepository    repository.Redis
}

type Service interface {
	MyProfile(ctx context.Context) (dto.ResponseMyProfile, error)
	UpdateMyProfile(ctx context.Context, reqHandler dto.PayloadMyProfile) error
}

func NewService(f *factory.Factory) Service {
	return &service{
		ConsumerRepository: f.ConsumerRepository,
		RedisRepository:    f.RedisRepository,
	}
}

func (s *service) MyProfile(ctx context.Context) (dto.ResponseMyProfile, error) {
	var (
		res dto.ResponseMyProfile
	)

	userSession := ctx.Value("user").(dto.JwtSession)

	if userSession.ConsumerID == 0 {
		return res, fmt.Errorf("failed to get consumer data")
	}

	consumer, err := s.ConsumerRepository.FindOne(ctx, "*", "id = ?", userSession.ConsumerID)
	if err != nil {
		return res, err
	}

	res = dto.ResponseMyProfile{
		NIK:        consumer.NIK,
		FullName:   consumer.FullName,
		LegalName:  consumer.LegalName,
		BirthDate:  consumer.BirthDate.Format(consts.TimeFormatDate),
		BirthPlace: consumer.BirthPlace,
		Salary:     int(consumer.Salary),
		KTPURL:     consumer.KtpURL,
		SelfieURL:  consumer.SelfieURL,
	}

	return res, nil
}

func (s *service) UpdateMyProfile(ctx context.Context, reqHandler dto.PayloadMyProfile) error {
	userSession := ctx.Value("user").(dto.JwtSession)

	if userSession.ConsumerID == 0 {
		return fmt.Errorf("failed to get consumer data")
	}

	fetch, err := s.ConsumerRepository.FindOne(ctx, "*", "id = ?", userSession.ConsumerID)
	if err != nil {
		return err
	}

	appUrl := util.GetEnv("APP_URL", "http://localhost")
	appPort := util.GetEnv("APP_PORT", "8080")
	baseURL := fmt.Sprintf("%s:%s/", appUrl, appPort)
	updateModel := model.Consumer{}
	tx := database.BeginTx(ctx, factory.NewFactory().InitDB)

	if reqHandler.KTPURL != "" {
		updateModel.KtpURL = reqHandler.KTPURL

		oldLink := fetch.KtpURL
		sanitizedLink := strings.Replace(fetch.KtpURL, baseURL, "", 1)

		if sanitizedLink != oldLink {
			err := util.DeleteFile(sanitizedLink)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}

	if reqHandler.SelfieURL != "" {
		updateModel.SelfieURL = reqHandler.SelfieURL

		oldLink := fetch.SelfieURL
		sanitizedLink := strings.Replace(fetch.SelfieURL, baseURL, "", 1)

		if sanitizedLink != oldLink {
			err := util.DeleteFile(sanitizedLink)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}

	if reqHandler.BirthDate != "" {
		birthDate, err := time.Parse(consts.TimeFormatDate, reqHandler.BirthDate)
		if err != nil {
			return err
		}

		updateModel.BirthDate = birthDate
	}

	if reqHandler.BirthPlace != "" {
		updateModel.BirthPlace = reqHandler.BirthPlace
	}

	if reqHandler.FullName != "" {
		updateModel.FullName = reqHandler.FullName
	}

	if reqHandler.LegalName != "" {
		updateModel.LegalName = reqHandler.LegalName
	}

	if reqHandler.NIK != "" {
		updateModel.NIK = reqHandler.NIK
	}

	if reqHandler.Salary != 0 {
		updateModel.Salary = float64(reqHandler.Salary)
	}

	err = s.ConsumerRepository.UpdateOne(tx, userSession.ConsumerID, updateModel)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
