package repository

import (
	"clean-arch/internal/model"
	"clean-arch/pkg/util"
	"context"

	"gorm.io/gorm"
)

type Installment interface {
	Store(db *gorm.DB, installment model.Installment) error
	FindAll(ctx context.Context, selectedFields string, limit, offset int, query string, args ...interface{}) ([]*model.Installment, error)
	FindOne(ctx context.Context, selectedFields string, query string, args ...any) (model.Installment, error)
	UpdateOne(db *gorm.DB, id int, data model.Installment) error
	DeleteOne(db *gorm.DB, id int) error
	Count(ctx context.Context, query string, args ...any) (int, error)
}

type installment struct {
	Db *gorm.DB
}

func NewInstallmentRepository(db *gorm.DB) Installment {
	return &installment{
		Db: db,
	}
}

func (r *installment) Store(db *gorm.DB, installment model.Installment) error {
	if err := db.Create(&installment).Error; err != nil {
		return err
	}
	return nil
}

func (r *installment) Count(ctx context.Context, query string, args ...any) (int, error) {
	var (
		res int64
	)

	err := r.Db.WithContext(ctx).Model(model.Installment{}).Select("id").Where(query, args...).Count(&res).Error
	if err != nil {
		return 0, err
	}

	return int(res), nil
}

func (r *installment) FindAll(ctx context.Context, selectedFields string, limit, offset int, query string, args ...interface{}) ([]*model.Installment, error) {
	var res []*model.Installment
	db := r.Db.Model(&model.Installment{})
	db = util.SetSelectFields(db, selectedFields)
	db = db.Where(query, args...).Limit(limit).Offset(offset)

	err := db.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *installment) FindOne(ctx context.Context, selectedFields string, query string, args ...any) (model.Installment, error) {
	var res model.Installment

	db := r.Db.WithContext(ctx).Model(model.Installment{})
	db = util.SetSelectFields(db, selectedFields)

	if err := db.Where(query, args...).Take(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (r *installment) UpdateOne(db *gorm.DB, id int, data model.Installment) error {
	if err := db.Model(&model.Installment{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func (r *installment) DeleteOne(db *gorm.DB, id int) error {
	if err := db.Delete(&model.Installment{}, id).Error; err != nil {
		return err
	}
	return nil
}
