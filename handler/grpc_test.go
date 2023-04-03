package handler_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/DragFAQ/uuid-generator/generator"
	"github.com/DragFAQ/uuid-generator/handler"
	mocklog "github.com/DragFAQ/uuid-generator/mocks"
	pb "github.com/DragFAQ/uuid-generator/proto"
)

func TestNewGrpcHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock_logger := mocklog.NewMockLogger(ctrl)
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
	mock_logger := mocklog.NewMockLogger(ctrl)
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
