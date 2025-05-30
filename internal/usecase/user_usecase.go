package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tubagusmf/tbwallet-user-auth/internal/helper"
	"github.com/tubagusmf/tbwallet-user-auth/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo model.IUserRepository
}

func NewUserUsecase(userRepo model.IUserRepository) model.IUserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

func (u *UserUsecase) GetAll(ctx context.Context, user model.User) ([]*model.User, error) {
	log := logrus.WithFields(logrus.Fields{
		"user": user,
	})

	users, err := u.userRepo.GetAll(ctx, user)
	if err != nil {
		log.WithError(err).Error("failed to get all users")
		return nil, err
	}

	return users, nil
}

func (u *UserUsecase) GetByID(ctx context.Context, id int64) (*model.User, error) {
	log := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		log.WithError(err).Error("failed to get user by id")
		return nil, err
	}

	return user, nil
}

func (u *UserUsecase) Create(ctx context.Context, input model.CreateUserInput) (token string, err error) {
	log := logrus.WithFields(logrus.Fields{
		"input": input,
	})

	err = helper.Validator.Struct(input)
	if err != nil {
		log.Error("Validation error: ", err)
		return
	}

	passwordHashed, err := helper.HashRequestPassword(input.PasswordHash)
	if err != nil {
		log.WithError(err).Error("failed to hash password")
		return "", err
	}

	newUser, err := u.userRepo.Create(ctx, model.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: passwordHashed,
		Role:         input.Role,
	})
	if err != nil {
		log.Error(err)
		return
	}

	accessToken, err := helper.GenerateToken(newUser.ID)
	if err != nil {
		log.Error(err)
		return
	}

	return accessToken, nil
}

func (u *UserUsecase) Update(ctx context.Context, id int64, input model.UpdateUserInput) (*model.User, error) {
	log := logrus.WithFields(logrus.Fields{
		"id":    id,
		"input": input,
	})

	err := helper.Validator.Struct(input)
	if err != nil {
		log.Error("Validation error: ", err)
		return &model.User{}, err
	}

	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		log.WithError(err).Error("failed to get user by id")
		return nil, err
	}

	if user == nil || (user.DeletedAt != nil && !user.DeletedAt.IsZero()) {
		log.Error("User is deleted or does not exist")
		return nil, errors.New("user is deleted or does not exist")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to hash password: ", err)
		return nil, err
	}

	updatedUser, err := u.userRepo.Update(ctx, model.User{
		ID:           user.ID,
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		Role:         input.Role,
	})
	if err != nil {
		log.WithError(err).Error("failed to update user")
		return nil, err
	}

	return updatedUser, nil
}

func (u *UserUsecase) Delete(ctx context.Context, id int64) error {
	log := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		log.Error("Failed to find user for deletion: ", err)
		return err
	}

	if user == nil {
		log.Error("User not found")
		return err
	}

	now := time.Now()
	user.DeletedAt = &now

	err = u.userRepo.Delete(ctx, id)
	if err != nil {
		log.Error("Failed to delete user: ", err)
		return err
	}

	log.Info("Successfully deleted user with ID: ", id)
	return nil
}

func (u *UserUsecase) ValidateSession(ctx context.Context, token model.UserSession) (*model.UserSession, error) {
	log := logrus.WithFields(logrus.Fields{
		"token": token,
	})

	session, err := u.userRepo.FindSessionByToken(ctx, token.Token)
	if err != nil {
		log.WithError(err).Error("failed to validate session")
		return nil, err
	}

	return session, nil
}

func (u *UserUsecase) Login(ctx context.Context, input model.LoginInput) (token string, err error) {
	log := logrus.WithFields(logrus.Fields{
		"input": input,
	})

	err = helper.Validator.Struct(input)
	if err != nil {
		log.Error("Validation error: ", err)
		return "", err
	}

	user, err := u.userRepo.GetByEmail(ctx, input.Email)
	if user == nil {
		return "", errors.New("user not found")
	}

	if !helper.CheckPasswordHash(input.Password, user.PasswordHash) {
		return "", errors.New("mismatch password")
	}

	token, err = helper.GenerateToken(user.ID)
	if err != nil {
		log.Error(err)
		return "", err
	}

	err = u.userRepo.DeleteSession(ctx, token)
	if err != nil {
		log.Warn("Failed to delete old session:", err)
	}

	session := model.UserSession{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}
	_, err = u.userRepo.CreateSession(ctx, session)
	if err != nil {
		log.Error("Failed to create session:", err)
		return "", err
	}

	return token, nil
}

func (u *UserUsecase) Logout(ctx context.Context, token model.UserSession) error {
	log := logrus.WithFields(logrus.Fields{
		"token": token,
	})

	err := u.userRepo.DeleteSession(ctx, token.Token)
	if err != nil {
		log.WithError(err).Error("failed to delete session")
		return err
	}

	return nil
}
