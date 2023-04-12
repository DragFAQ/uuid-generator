package handler_test

import (
	"context"
	"testing"

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
	handl := handler.NewGrpcHandler(mock_logger)

	assert.NotNil(t, handl)
}

func TestGrpcGetCurrentHash(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock_logger := mocklog.NewMockLogger(ctrl)
	currentHash := generator.GetHash()
	handl := handler.NewGrpcHandler(mock_logger)

	response, err := handl.GetCurrentHash(context.Background(), &pb.HashRequest{})

	assert.NoError(t, err)
	assert.Equal(t, response.Hash, currentHash.Value)
}
