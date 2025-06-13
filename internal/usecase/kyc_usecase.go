package usecase

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/tubagusmf/tbwallet-user-auth/internal/helper"
	"github.com/tubagusmf/tbwallet-user-auth/internal/model"
	"github.com/tubagusmf/tbwallet-user-auth/pb/kycdoc"
)

type KycUsecase struct {
	kycRepo   model.IKycDocRepository
	kycClient kycdoc.KycdocServiceClient
}

func NewKycUsecase(kycRepo model.IKycDocRepository, kycClient kycdoc.KycdocServiceClient) model.IKycDocUsecase {
	return &KycUsecase{
		kycRepo:   kycRepo,
		kycClient: kycClient,
	}
}

func (k *KycUsecase) GetByID(ctx context.Context, id int64) (*model.KycDocument, error) {
	log := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	doc, err := k.kycRepo.GetByID(ctx, id)
	if err != nil {
		log.WithError(err).Error("failed to get kyc by id")
		return nil, err
	}

	return doc, nil
}

func (k *KycUsecase) GetByUserID(ctx context.Context, userID int64) ([]model.KycDocument, error) {
	log := logrus.WithFields(logrus.Fields{
		"user_id": userID,
	})

	docs, err := k.kycRepo.GetByUserID(ctx, userID)
	if err != nil {
		log.WithError(err).Error("failed to get kyc by user id")
		return nil, err
	}

	return docs, nil
}

func (k *KycUsecase) Create(ctx context.Context, input model.CreateKycDocInput) (*model.KycDocument, error) {
	log := logrus.WithFields(logrus.Fields{
		"input": input,
	})

	err := helper.Validator.Struct(input)
	if err != nil {
		log.Error("Validation error: ", err)
		return nil, err
	}

	defaultStatus := "pending"
	doc := model.KycDocument{
		UserID:       input.UserID,
		DocumentType: input.DocumentType,
		DocumentURL:  input.DocumentURL,
		Status:       defaultStatus,
	}

	createDoc, err := k.kycRepo.Create(ctx, doc)
	if err != nil {
		log.WithError(err).Error("failed to create kyc")
		return nil, err
	}

	return createDoc, nil
}

func (k *KycUsecase) Update(ctx context.Context, id int64, input model.UpdateKycDocInput) (*model.KycDocument, error) {
	log := logrus.WithFields(logrus.Fields{
		"id":    id,
		"input": input,
	})

	err := helper.Validator.Struct(input)
	if err != nil {
		log.Error("Validation error: ", err)
		return nil, err
	}

	doc := model.KycDocument{
		UserID:       input.UserID,
		DocumentType: input.DocumentType,
		DocumentURL:  input.DocumentURL,
	}

	updateDoc, err := k.kycRepo.Update(ctx, id, doc)
	if err != nil {
		log.WithError(err).Error("failed to update kyc")
		return nil, err
	}

	return updateDoc, nil
}

func (k *KycUsecase) ValidateStatus(ctx context.Context, id int64, input model.ValidateStatusInput) (*model.KycDocument, error) {
	log := logrus.WithFields(logrus.Fields{
		"id":    id,
		"input": input,
	})

	err := helper.Validator.Struct(input)
	if err != nil {
		log.Error("Validation error: ", err)
		return nil, err
	}

	doc, err := k.kycRepo.ValidateStatus(ctx, id, input.Status)
	if err != nil {
		log.WithError(err).Error("failed to validate kyc status")
		return nil, err
	}

	return doc, nil
}

func (k *KycUsecase) GetKycStatus(ctx context.Context, userID int64) (*model.KycDocument, error) {
	log := logrus.WithFields(logrus.Fields{
		"user_id": userID,
	})

	status, err := k.kycRepo.GetKycStatus(ctx, userID)
	if err != nil {
		log.WithError(err).Error("failed to get kyc status")
		return nil, err
	}

	return status, nil
}
