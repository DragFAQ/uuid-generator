package handler_test

import (
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/DragFAQ/uuid-generator/generator"
	"github.com/DragFAQ/uuid-generator/handler"
	mocklog "github.com/DragFAQ/uuid-generator/mocks"
)

func TestNewHttpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock_logger := mocklog.NewMockLogger(ctrl)
	hash := "test-hash"
	timestamp := time.Now()
	currentHash := &generator.Hash{
		Value:          hash,
		GenerationTime: timestamp,
	}
	hashLock := sync.RWMutex{}

	handl := handler.NewHttpHandler(currentHash, &hashLock, mock_logger)

	assert.NotNil(t, handl)
}

func TestHttpGetCurrentHash(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock_logger := mocklog.NewMockLogger(ctrl)
	hash := "test-hash"
	timestamp := time.Now()
	currentHash := &generator.Hash{
		Value:          hash,
		GenerationTime: timestamp,
	}
	hashLock := sync.RWMutex{}
	handl := handler.NewHttpHandler(currentHash, &hashLock, mock_logger)

	recorder := httptest.NewRecorder()
	handl.GetCurrentHash(recorder, &http.Request{})

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, `{"generation_time":"`+timestamp.Format(time.RFC3339)+`","hash":"test-hash"}`, recorder.Body.String())
}
