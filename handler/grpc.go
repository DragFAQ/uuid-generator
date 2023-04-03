package handler

import (
	"context"
	"sync"
	"time"

	generator "github.com/DragFAQ/uuid-generator/generator"
	log "github.com/DragFAQ/uuid-generator/logger"
	pb "github.com/DragFAQ/uuid-generator/proto"
)

type GrpcHandler struct {
	logger      log.Logger
	currentHash *generator.Hash
	hashLock    *sync.RWMutex
	pb.UnimplementedHashServiceServer
}

func NewGrpcHandler(currentHash *generator.Hash, hashLock *sync.RWMutex, logger log.Logger) *GrpcHandler {
	return &GrpcHandler{
		logger:      logger,
		currentHash: currentHash,
		hashLock:    hashLock,
	}
}

func (h *GrpcHandler) GetCurrentHash(ctx context.Context, r *pb.HashRequest) (*pb.HashResponse, error) {
	h.hashLock.RLock()
	defer h.hashLock.RUnlock()

	resp := &pb.HashResponse{
		Hash:           h.currentHash.Value,
		GenerationTime: h.currentHash.GenerationTime.Format(time.RFC3339),
	}

	return resp, nil
}
