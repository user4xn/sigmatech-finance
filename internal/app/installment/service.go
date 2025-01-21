package installment

import (
	"clean-arch/internal/factory"
	"clean-arch/internal/repository"
)

type service struct {
	InstallmentRepository repository.Installment
	RedisRepository       repository.Redis
}

type Service interface {
}

func NewService(f *factory.Factory) Service {
	return &service{
		InstallmentRepository: f.InstallmentRepository,
		RedisRepository:       f.RedisRepository,
	}
}
