package grpc

import (
	"context"

	"github.com/tubagusmf/tbwallet-user-auth/internal/model"
	pb "github.com/tubagusmf/tbwallet-user-auth/pb/user"
)

type UsergRPCHandler struct {
	pb.UnimplementedUserServiceServer
	userUseCase model.IUserUsecase
}

func NewUsergRPCHandler(userUseCase model.IUserUsecase) pb.UserServiceServer {
	return &UsergRPCHandler{userUseCase: userUseCase}
}

func (u *UsergRPCHandler) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	user, err := u.userUseCase.GetByID(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	var kycStatusID int64
	if user.KycStatusID != nil {
		kycStatusID = *user.KycStatusID
	}

	response := &pb.GetUserByIDResponse{
		User: &pb.User{
			Id:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			KycStatusId: kycStatusID,
			Role:        user.Role,
		},
	}

	return response, nil
}

func (u *UsergRPCHandler) ValidateSession(ctx context.Context, req *pb.ValidateSessionRequest) (*pb.ValidateSessionResponse, error) {
	session, err := u.userUseCase.ValidateSession(ctx, model.UserSession{Token: req.GetToken()})
	if err != nil {
		return nil, err
	}

	var kycStatusID int64
	if session.KycStatusID != nil {
		kycStatusID = *session.KycStatusID
	}

	response := &pb.ValidateSessionResponse{
		User: &pb.SessionUser{
			Id:          session.ID,
			KycStatusId: kycStatusID,
		},
	}

	return response, nil
}
