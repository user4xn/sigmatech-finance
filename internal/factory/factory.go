package factory

import (
	"clean-arch/database"
	"clean-arch/internal/repository"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Factory struct {
	TransactionPenaltyRepository repository.TransactionPenalty
	InstallmentRepository        repository.Installment
	TransactionRepository        repository.Transaction
	LimitRepository              repository.Limit
	ConsumerRepository           repository.Consumer
	UserRepository               repository.User
	RedisRepository              repository.Redis
	RedisClient                  *redis.Client
	InitDB                       *gorm.DB
}

func NewFactory() *Factory {
	// Check db connection
	db := database.GetConnection()
	rdb := database.GetRedisClient()

	return &Factory{
		// Pass the db connection to repository package for database query calling
		TransactionPenaltyRepository: repository.NewTransactionPenaltyRepository(db),
		InstallmentRepository:        repository.NewInstallmentRepository(db),
		TransactionRepository:        repository.NewTransactionRepository(db),
		LimitRepository:              repository.NewLimitRepository(db),
		ConsumerRepository:           repository.NewConsumerRepository(db),
		UserRepository:               repository.NewUserRepository(db),
		RedisRepository:              repository.NewRedisRepository(rdb),
		RedisClient:                  rdb,
		InitDB:                       db,
	}
}
