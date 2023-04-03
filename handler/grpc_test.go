package handler_test

import (
	"context"
	mock_log "github.com/DragFAQ/uuid-generator/mocks"
	"github.com/golang/mock/gomock"
	"sync"
	"testing"
	"time"

	"github.com/DragFAQ/uuid-generator/generator"
	pb "github.com/DragFAQ/uuid-generator/proto"
	"github.com/stretchr/testify/assert"

	"github.com/DragFAQ/uuid-generator/handler"
)

func TestNewGrpcHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock_logger := mock_log.NewMockLogger(ctrl)
	hash := "test-hash"
	timestamp := time.Now()
	currentHash := &generator.Hash{
		Value:          hash,
		GenerationTime: timestamp,
	}
	hashLock := sync.RWMutex{}

	handl := handler.NewGrpcHandler(currentHash, &hashLock, mock_logger)

	assert.NotNil(t, handl)
}

func TestGrpcGetCurrentHash(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock_logger := mock_log.NewMockLogger(ctrl)
	hash := "test-hash"
	timestamp := time.Now()
	currentHash := &generator.Hash{
		Value:          hash,
		GenerationTime: timestamp,
	}
	hashLock := sync.RWMutex{}
	handl := handler.NewGrpcHandler(currentHash, &hashLock, mock_logger)

	response, err := handl.GetCurrentHash(context.Background(), &pb.HashRequest{})

	assert.NoError(t, err)
	assert.Equal(t, response.Hash, hash)
	assert.Equal(t, response.GenerationTime, timestamp.Format(time.RFC3339))
}
