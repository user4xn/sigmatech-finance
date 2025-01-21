package repository

import (
	"clean-arch/internal/model"
	"clean-arch/pkg/util"
	"context"

	"gorm.io/gorm"
)

type Limit interface {
	Store(db *gorm.DB, consumerLimit model.Limit) error
	FindAll(ctx context.Context, selectedFields string, limit, offset int, query string, args ...interface{}) ([]*model.Limit, error)
	FindOne(ctx context.Context, selectedFields string, query string, args ...any) (model.Limit, error)
	UpdateOne(db *gorm.DB, id int, data model.Limit) error
	DeleteOne(db *gorm.DB, id int) error
	Count(ctx context.Context, query string, args ...any) (int, error)
}

type consumerLimit struct {
	Db *gorm.DB
}

func NewLimitRepository(db *gorm.DB) Limit {
	return &consumerLimit{
		Db: db,
	}
}

func (r *consumerLimit) Store(db *gorm.DB, consumerLimit model.Limit) error {
	if err := db.Create(&consumerLimit).Error; err != nil {
		return err
	}
	return nil
}

func (r *consumerLimit) Count(ctx context.Context, query string, args ...any) (int, error) {
	var (
		res int64
	)

	err := r.Db.WithContext(ctx).Model(model.Limit{}).Select("id").Where(query, args...).Count(&res).Error
	if err != nil {
		return 0, err
	}

	return int(res), nil
}

func (r *consumerLimit) FindAll(ctx context.Context, selectedFields string, limit, offset int, query string, args ...interface{}) ([]*model.Limit, error) {
	var res []*model.Limit
	db := r.Db.Model(&model.Limit{})
	db = util.SetSelectFields(db, selectedFields)
	db = db.Where(query, args...).Limit(limit).Offset(offset)

	err := db.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *consumerLimit) FindOne(ctx context.Context, selectedFields string, query string, args ...any) (model.Limit, error) {
	var res model.Limit

	db := r.Db.WithContext(ctx).Model(model.Limit{})
	db = util.SetSelectFields(db, selectedFields)

	if err := db.Where(query, args...).Take(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (r *consumerLimit) UpdateOne(db *gorm.DB, id int, data model.Limit) error {
	if err := db.Model(&model.Limit{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func (r *consumerLimit) DeleteOne(db *gorm.DB, id int) error {
	if err := db.Delete(&model.Limit{}, id).Error; err != nil {
		return err
	}
	return nil
}
