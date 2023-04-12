package generator_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/DragFAQ/uuid-generator/generator"
	mocklog "github.com/DragFAQ/uuid-generator/mocks"
)

func TestGenerate(t *testing.T) {
	ctrl := gomock.NewController(t)

	mock_logger := mocklog.NewMockLogger(ctrl)
	mock_logger.EXPECT().Debugf("%s: New UUID was generated '%s'", gomock.Any(), gomock.Any()).AnyTimes()
	mock_logger.EXPECT().Infof("GenerateHash worker stopped.").AnyTimes()

	ctx, cancel := context.WithCancel(context.Background())
	go generator.GenerateHash(ctx, mock_logger, 1)

	time.Sleep(1100 * time.Millisecond)

	first := generator.GetHash()

	assert.NotEmpty(t, first.Value)

	_, err := uuid.Parse(first.Value)
	assert.NoError(t, err)

	assert.Equal(t, first.GenerationTime.IsZero(), false)

	cancel()
	time.Sleep(100 * time.Millisecond)
}

func TestChangedValAfterTTL(t *testing.T) {
	ctrl := gomock.NewController(t)

	mock_logger := mocklog.NewMockLogger(ctrl)
	mock_logger.EXPECT().Debugf("%s: New UUID was generated '%s'", gomock.Any(), gomock.Any()).AnyTimes()
	mock_logger.EXPECT().Infof("GenerateHash worker stopped.").AnyTimes()

	ctx, cancel := context.WithCancel(context.Background())
	go generator.GenerateHash(ctx, mock_logger, 1)

	time.Sleep(1100 * time.Millisecond)

	first := generator.GetHash()

	assert.NotEmpty(t, first.Value)

	time.Sleep(1100 * time.Millisecond)

	second := generator.GetHash()
	assert.NotEqual(t, first.Value, second.Value)

	cancel()
	time.Sleep(100 * time.Millisecond)
}

func TestNotChangedBeforeTTL(t *testing.T) {
	ctrl := gomock.NewController(t)

	mock_logger := mocklog.NewMockLogger(ctrl)
	mock_logger.EXPECT().Debugf("%s: New UUID was generated '%s'", gomock.Any(), gomock.Any()).AnyTimes()
	mock_logger.EXPECT().Infof("GenerateHash worker stopped.").AnyTimes()

	ctx, cancel := context.WithCancel(context.Background())
	go generator.GenerateHash(ctx, mock_logger, 1)

	time.Sleep(1100 * time.Millisecond)

	first := generator.GetHash()

	time.Sleep(500 * time.Millisecond)

	second := generator.GetHash()
	assert.Equal(t, first.Value, second.Value)

	cancel()
	time.Sleep(100 * time.Millisecond)
}
