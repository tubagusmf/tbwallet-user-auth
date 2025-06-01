package usecase

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/tubagusmf/tbwallet-user-auth/internal/helper"
	"github.com/tubagusmf/tbwallet-user-auth/internal/model"
)

type KycUsecase struct {
	kycRepo model.IKycDocRepository
}

func NewKycUsecase(kycRepo model.IKycDocRepository) model.IKycDocUsecase {
	return &KycUsecase{kycRepo: kycRepo}
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

	doc := model.KycDocument{
		UserID:       input.UserID,
		DocumentType: input.DocumentType,
		DocumentURL:  input.DocumentURL,
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
