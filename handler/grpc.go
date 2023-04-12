package handler

import (
	"context"
	"time"

	hash "github.com/DragFAQ/uuid-generator/generator"
	log "github.com/DragFAQ/uuid-generator/logger"
	pb "github.com/DragFAQ/uuid-generator/proto"
)

type GrpcHandler struct {
	logger log.Logger
	pb.UnimplementedHashServiceServer
}

func NewGrpcHandler(logger log.Logger) *GrpcHandler {
	return &GrpcHandler{
		logger: logger,
	}
}

func (h *GrpcHandler) GetCurrentHash(_ context.Context, _ *pb.HashRequest) (*pb.HashResponse, error) {
	currentHash := hash.GetHash()

	resp := &pb.HashResponse{
		Hash:           currentHash.Value,
		GenerationTime: currentHash.GenerationTime.Format(time.RFC3339),
	}

	return resp, nil
}
