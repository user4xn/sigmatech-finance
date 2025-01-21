package factory

import (
	"clean-arch/database"
	"clean-arch/internal/repository"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Factory struct {
	UserRepository  repository.User
	RedisRepository repository.Redis
	RedisClient     *redis.Client
	InitDB          *gorm.DB
}

func NewFactory() *Factory {
	// Check db connection
	db := database.GetConnection()
	rdb := database.GetRedisClient()

	return &Factory{
		// Pass the db connection to repository package for database query calling
		UserRepository:  repository.NewUserRepository(db),
		RedisRepository: repository.NewRedisRepository(rdb),
		RedisClient:     rdb,
		InitDB:          db,
	}
}
