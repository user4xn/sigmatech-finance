package repository

import (
	"clean-arch/internal/model"
	"clean-arch/pkg/util"
	"context"

	"gorm.io/gorm"
)

type TransactionPenalty interface {
	Store(db *gorm.DB, transactionPenalty model.TransactionPenalty) error
	FindAll(ctx context.Context, selectedFields string, limit, offset int, query string, args ...interface{}) ([]*model.TransactionPenalty, error)
	FindOne(ctx context.Context, selectedFields string, query string, args ...any) (model.TransactionPenalty, error)
	UpdateOne(db *gorm.DB, id int, data model.TransactionPenalty) error
	DeleteOne(db *gorm.DB, id int) error
	Count(ctx context.Context, query string, args ...any) (int, error)
}

type transactionPenalty struct {
	Db *gorm.DB
}

func NewTransactionPenaltyRepository(db *gorm.DB) TransactionPenalty {
	return &transactionPenalty{
		Db: db,
	}
}

func (r *transactionPenalty) Store(db *gorm.DB, transactionPenalty model.TransactionPenalty) error {
	if err := db.Create(&transactionPenalty).Error; err != nil {
		return err
	}
	return nil
}

func (r *transactionPenalty) Count(ctx context.Context, query string, args ...any) (int, error) {
	var (
		res int64
	)

	err := r.Db.WithContext(ctx).Model(model.TransactionPenalty{}).Select("id").Where(query, args...).Count(&res).Error
	if err != nil {
		return 0, err
	}

	return int(res), nil
}

func (r *transactionPenalty) FindAll(ctx context.Context, selectedFields string, limit, offset int, query string, args ...interface{}) ([]*model.TransactionPenalty, error) {
	var res []*model.TransactionPenalty
	db := r.Db.Model(&model.TransactionPenalty{})
	db = util.SetSelectFields(db, selectedFields)
	db = db.Where(query, args...).Limit(limit).Offset(offset)

	err := db.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *transactionPenalty) FindOne(ctx context.Context, selectedFields string, query string, args ...any) (model.TransactionPenalty, error) {
	var res model.TransactionPenalty

	db := r.Db.WithContext(ctx).Model(model.TransactionPenalty{})
	db = util.SetSelectFields(db, selectedFields)

	if err := db.Where(query, args...).Take(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (r *transactionPenalty) UpdateOne(db *gorm.DB, id int, data model.TransactionPenalty) error {
	if err := db.Model(&model.TransactionPenalty{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func (r *transactionPenalty) DeleteOne(db *gorm.DB, id int) error {
	if err := db.Delete(&model.TransactionPenalty{}, id).Error; err != nil {
		return err
	}
	return nil
}
