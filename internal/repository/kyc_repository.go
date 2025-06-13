package repository

import (
	"context"
	"errors"
	"time"

	"github.com/tubagusmf/tbwallet-user-auth/internal/model"
	"gorm.io/gorm"
)

type KycRepo struct {
	db *gorm.DB
}

func NewKycRepo(db *gorm.DB) model.IKycDocRepository {
	return &KycRepo{db: db}
}

func (k *KycRepo) GetByID(ctx context.Context, id int64) (*model.KycDocument, error) {
	var doc model.KycDocument
	err := k.db.WithContext(ctx).First(&doc, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("document not found")
	}

	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func (k *KycRepo) GetByUserID(ctx context.Context, userID int64) ([]model.KycDocument, error) {
	var docs []model.KycDocument
	err := k.db.WithContext(ctx).Where("user_id = ?", userID).Find(&docs).Error

	if err != nil {
		return nil, err
	}

	return docs, nil
}

func (k *KycRepo) Create(ctx context.Context, doc model.KycDocument) (*model.KycDocument, error) {
	doc.CreatedAt = time.Now()
	doc.UpdatedAt = time.Now()

	err := k.db.WithContext(ctx).Create(&doc).Error
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func (k *KycRepo) Update(ctx context.Context, id int64, doc model.KycDocument) (*model.KycDocument, error) {
	doc.UpdatedAt = time.Now()

	err := k.db.WithContext(ctx).Where("id = ?", id).Updates(&doc).Error
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func (k *KycRepo) ValidateStatus(ctx context.Context, id int64, status string) (*model.KycDocument, error) {
	var doc model.KycDocument
	if err := k.db.WithContext(ctx).First(&doc, id).Error; err != nil {
		return nil, err
	}

	if doc.Status == status {
		return &doc, nil
	}

	doc.Status = status
	doc.UpdatedAt = time.Now()

	if err := k.db.WithContext(ctx).Save(&doc).Error; err != nil {
		return nil, err
	}

	return &doc, nil
}

func (k *KycRepo) GetKycStatus(ctx context.Context, userID int64) (*model.KycDocument, error) {
	var doc model.KycDocument
	if err := k.db.WithContext(ctx).Where("user_id = ?", userID).First(&doc).Error; err != nil {
		return nil, err
	}

	return &doc, nil
}
