package model

import (
	"context"
	"time"
)

type KycDocument struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	DocumentType string    `json:"document_type"`
	DocumentURL  string    `json:"document_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateKycDocInput struct {
	UserID       int64  `json:"user_id" validate:"required"`
	DocumentType string `json:"document_type" validate:"required"`
	DocumentURL  string `json:"document_url" validate:"required"`
}

type UpdateKycDocInput struct {
	UserID       int64  `json:"user_id" validate:"required"`
	DocumentType string `json:"document_type" validate:"required"`
	DocumentURL  string `json:"document_url" validate:"required"`
}

type IKycDocRepository interface {
	GetByID(ctx context.Context, id int64) (*KycDocument, error)
	GetByUserID(ctx context.Context, userID int64) ([]KycDocument, error)
	Create(ctx context.Context, doc KycDocument) (*KycDocument, error)
	Update(ctx context.Context, id int64, doc KycDocument) (*KycDocument, error)
}

type IKycDocUsecase interface {
	GetByID(ctx context.Context, id int64) (*KycDocument, error)
	GetByUserID(ctx context.Context, userID int64) ([]KycDocument, error)
	Create(ctx context.Context, input CreateKycDocInput) (*KycDocument, error)
	Update(ctx context.Context, id int64, input UpdateKycDocInput) (*KycDocument, error)
}
