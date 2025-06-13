package grpc

import (
	"context"

	"github.com/tubagusmf/tbwallet-user-auth/internal/model"
	pb "github.com/tubagusmf/tbwallet-user-auth/pb/kycdoc"
)

type KycdocGRPCHandler struct {
	pb.UnimplementedKycdocServiceServer
	kycdocService model.IKycDocUsecase
}

func NewKycdocGRPCHandler(kycdocService model.IKycDocUsecase) pb.KycdocServiceServer {
	return &KycdocGRPCHandler{kycdocService: kycdocService}
}

func (h *KycdocGRPCHandler) GetKycdocByUserID(ctx context.Context, req *pb.GetKycdocByUserIDRequest) (*pb.GetKycdocByUserIDResponse, error) {
	kycdocs, err := h.kycdocService.GetByUserID(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	var pbKycdocs []*pb.Kycdoc
	for _, k := range kycdocs {
		pbKycdocs = append(pbKycdocs, &pb.Kycdoc{
			Id:           k.ID,
			UserId:       k.UserID,
			DocumentType: k.DocumentType,
			DocumentUrl:  k.DocumentURL,
			Status:       k.Status,
		})
	}

	return &pb.GetKycdocByUserIDResponse{Kycdocs: pbKycdocs}, nil
}

func (h *KycdocGRPCHandler) GetKycStatus(ctx context.Context, req *pb.GetKycStatusRequest) (*pb.GetKycStatusResponse, error) {
	doc, err := h.kycdocService.GetKycStatus(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	return &pb.GetKycStatusResponse{
		Status: doc.Status,
	}, nil
}
