package repository

import (
	"clean-arch/internal/model"
	"clean-arch/pkg/util"
	"context"

	"gorm.io/gorm"
)

type Transaction interface {
	Store(db *gorm.DB, transaction model.Transaction) error
	FindAll(ctx context.Context, selectedFields string, limit, offset int, query string, args ...interface{}) ([]*model.Transaction, error)
	FindOne(ctx context.Context, selectedFields string, query string, args ...any) (model.Transaction, error)
	UpdateOne(db *gorm.DB, id int, data model.Transaction) error
	DeleteOne(db *gorm.DB, id int) error
	Count(ctx context.Context, query string, args ...any) (int, error)
}

type transaction struct {
	Db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) Transaction {
	return &transaction{
		Db: db,
	}
}

func (r *transaction) Store(db *gorm.DB, transaction model.Transaction) error {
	if err := db.Create(&transaction).Error; err != nil {
		return err
	}
	return nil
}

func (r *transaction) Count(ctx context.Context, query string, args ...any) (int, error) {
	var (
		res int64
	)

	err := r.Db.WithContext(ctx).Model(model.Transaction{}).Select("id").Where(query, args...).Count(&res).Error
	if err != nil {
		return 0, err
	}

	return int(res), nil
}

func (r *transaction) FindAll(ctx context.Context, selectedFields string, limit, offset int, query string, args ...interface{}) ([]*model.Transaction, error) {
	var res []*model.Transaction
	db := r.Db.Model(&model.Transaction{})
	db = util.SetSelectFields(db, selectedFields)
	db = db.Where(query, args...).Limit(limit).Offset(offset)

	err := db.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *transaction) FindOne(ctx context.Context, selectedFields string, query string, args ...any) (model.Transaction, error) {
	var res model.Transaction

	db := r.Db.WithContext(ctx).Model(model.Transaction{})
	db = util.SetSelectFields(db, selectedFields)

	if err := db.Where(query, args...).Take(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (r *transaction) UpdateOne(db *gorm.DB, id int, data model.Transaction) error {
	if err := db.Model(&model.Transaction{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func (r *transaction) DeleteOne(db *gorm.DB, id int) error {
	if err := db.Delete(&model.Transaction{}, id).Error; err != nil {
		return err
	}
	return nil
}
