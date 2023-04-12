package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/DragFAQ/uuid-generator/generator"
	"github.com/DragFAQ/uuid-generator/handler"
	mocklog "github.com/DragFAQ/uuid-generator/mocks"
)

func TestNewHttpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock_logger := mocklog.NewMockLogger(ctrl)
	handl := handler.NewHttpHandler(mock_logger)

	assert.NotNil(t, handl)
}

func TestHttpGetCurrentHash(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock_logger := mocklog.NewMockLogger(ctrl)
	currentHash := generator.GetHash()
	handl := handler.NewHttpHandler(mock_logger)

	recorder := httptest.NewRecorder()
	handl.GetCurrentHash(recorder, &http.Request{})

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, `{"generation_time":"`+currentHash.GenerationTime.Format(time.RFC3339)+`","hash":"`+currentHash.Value+`"}`, recorder.Body.String())
}
