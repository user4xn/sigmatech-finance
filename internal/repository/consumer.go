package repository

import (
	"clean-arch/internal/model"
	"clean-arch/pkg/util"
	"context"

	"gorm.io/gorm"
)

type Consumer interface {
	Store(db *gorm.DB, consumer model.Consumer) error
	FindAll(ctx context.Context, selectedFields string, limit, offset int, query string, args ...interface{}) ([]*model.Consumer, error)
	FindOne(ctx context.Context, selectedFields string, query string, args ...any) (model.Consumer, error)
	UpdateOne(db *gorm.DB, id int, data model.Consumer) error
	DeleteOne(db *gorm.DB, id int) error
	Count(ctx context.Context, query string, args ...any) (int, error)
}

type consumer struct {
	Db *gorm.DB
}

func NewConsumerRepository(db *gorm.DB) Consumer {
	return &consumer{
		Db: db,
	}
}

func (r *consumer) Store(db *gorm.DB, consumer model.Consumer) error {
	if err := db.Create(&consumer).Error; err != nil {
		return err
	}
	return nil
}

func (r *consumer) Count(ctx context.Context, query string, args ...any) (int, error) {
	var (
		res int64
	)

	err := r.Db.WithContext(ctx).Model(model.Consumer{}).Select("id").Where(query, args...).Count(&res).Error
	if err != nil {
		return 0, err
	}

	return int(res), nil
}

func (r *consumer) FindAll(ctx context.Context, selectedFields string, limit, offset int, query string, args ...interface{}) ([]*model.Consumer, error) {
	var res []*model.Consumer
	db := r.Db.Model(&model.Consumer{})
	db = util.SetSelectFields(db, selectedFields)
	db = db.Where(query, args...).Limit(limit).Offset(offset)

	err := db.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *consumer) FindOne(ctx context.Context, selectedFields string, query string, args ...any) (model.Consumer, error) {
	var res model.Consumer

	db := r.Db.WithContext(ctx).Model(model.Consumer{})
	db = util.SetSelectFields(db, selectedFields)

	if err := db.Where(query, args...).Take(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (r *consumer) UpdateOne(db *gorm.DB, id int, data model.Consumer) error {
	if err := db.Model(&model.Consumer{}).Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r *consumer) DeleteOne(db *gorm.DB, id int) error {
	if err := db.Delete(&model.Consumer{}, id).Error; err != nil {
		return err
	}
	return nil
}
