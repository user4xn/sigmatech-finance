package user

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

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type service struct {
	UserRepository repository.User
}

type Service interface {
	Store(ctx context.Context, reqHandler dto.PayloadUser) error
	FindAll(ctx context.Context, reqHandler dto.PayloadBasicTable) (*dto.ResponseUser, error)
	FindOne(ctx context.Context, id int) (dto.User, error)
	Update(ctx context.Context, id int, reqHandler dto.PayloadUpdateUser) error
	Delete(ctx context.Context, id int) error
}

func NewService(f *factory.Factory) Service {
	return &service{
		UserRepository: f.UserRepository,
	}
}

func (s *service) Store(ctx context.Context, reqHandler dto.PayloadUser) error {
	tx := database.BeginTx(ctx, factory.NewFactory().InitDB)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqHandler.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	insertModel := model.User{
		Email:    reqHandler.Email,
		Password: string(hashedPassword),
	}

	existingEmail, err := s.UserRepository.FindOne(ctx, "email", "email = ?", reqHandler.Email)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return err
		}
	}

	if existingEmail.Email != "" {
		return fmt.Errorf("email already exists")
	}

	if err := s.UserRepository.Store(tx, insertModel); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (s *service) FindAll(ctx context.Context, reqHandler dto.PayloadBasicTable) (*dto.ResponseUser, error) {
	var (
		total dto.ResponseTotalRow
		res   *dto.ResponseUser
		users []dto.User
		query string
		args  []interface{}
	)

	if reqHandler.Search != "" {
		query = "email LIKE ?"
		args = append(args, "%"+reqHandler.Search+"%")
	}

	count, err := s.UserRepository.Count(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	fetch, err := s.UserRepository.FindAll(ctx, "id, email, created_at, updated_at", reqHandler.Limit, reqHandler.Offset, query, args...)
	if err != nil {
		return nil, err
	}

	for _, user := range fetch {
		users = append(users, dto.User{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(consts.TimeFormatDateTime),
			UpdatedAt: user.UpdatedAt.Format(consts.TimeFormatDateTime),
		})
	}

	total = dto.ResponseTotalRow{
		TotalRow: count,
	}

	res = &dto.ResponseUser{
		ResponseTotalRow: total,
		Data:             users,
	}

	return res, nil
}

func (s *service) Update(ctx context.Context, id int, reqHandler dto.PayloadUpdateUser) error {
	tx := database.BeginTx(ctx, factory.NewFactory().InitDB)
	updatedModel := model.User{
		Email: reqHandler.Email,
	}

	if reqHandler.NewPassword != "" {
		if reqHandler.LastPassword == "" {
			return fmt.Errorf("last password is required")
		}

		user, err := s.UserRepository.FindOne(ctx, "password", "id = ?", id)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("user not found")
			}

			return err
		}

		err = util.ComparePasswords(user.Password, reqHandler.LastPassword)
		if err != nil {
			return fmt.Errorf("last password does not match")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqHandler.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		updatedModel.Password = string(hashedPassword)
	}

	if err := s.UserRepository.UpdateOne(tx, id, updatedModel); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (s *service) FindOne(ctx context.Context, id int) (dto.User, error) {
	var (
		res dto.User
	)

	fetch, err := s.UserRepository.FindOne(ctx, "*", "id = ?", id)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return res, err
		}
	}

	res = dto.User{
		ID:        fetch.ID,
		Email:     fetch.Email,
		CreatedAt: fetch.CreatedAt.Format(consts.TimeFormatDateTime),
	}

	return res, nil
}

func (s *service) Delete(ctx context.Context, id int) error {
	tx := database.BeginTx(ctx, factory.NewFactory().InitDB)

	if err := s.UserRepository.DeleteOne(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
