package generator_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/DragFAQ/uuid-generator/generator"
	mocklog "github.com/DragFAQ/uuid-generator/mocks"
)

func TestGenerate(t *testing.T) {
	hashLock := &sync.RWMutex{}
	ctrl := gomock.NewController(t)

	mock_logger := mocklog.NewMockLogger(ctrl)
	mock_logger.EXPECT().Debugf("%s: New UUID was generated '%s'", gomock.Any(), gomock.Any()).AnyTimes()
	mock_logger.EXPECT().Infof("GenerateHash worker stopped.").AnyTimes()

	currentHash := &generator.Hash{}
	ctx, cancel := context.WithCancel(context.Background())
	go generator.GenerateHash(currentHash, hashLock, mock_logger, 1, ctx)

	time.Sleep(1100 * time.Millisecond)

	hashLock.RLock()
	val := currentHash.Value
	tim := currentHash.GenerationTime
	hashLock.RUnlock()

	assert.NotEmpty(t, val)

	_, err := uuid.Parse(val)
	assert.NoError(t, err)

	assert.Equal(t, tim.IsZero(), false)

	cancel()
	time.Sleep(100 * time.Millisecond)
}

func TestChangedValAfterTTL(t *testing.T) {
	hashLock := &sync.RWMutex{}
	ctrl := gomock.NewController(t)

	mock_logger := mocklog.NewMockLogger(ctrl)
	mock_logger.EXPECT().Debugf("%s: New UUID was generated '%s'", gomock.Any(), gomock.Any()).AnyTimes()
	mock_logger.EXPECT().Infof("GenerateHash worker stopped.").AnyTimes()

	currentHash := &generator.Hash{}
	ctx, cancel := context.WithCancel(context.Background())
	go generator.GenerateHash(currentHash, hashLock, mock_logger, 1, ctx)

	time.Sleep(1100 * time.Millisecond)

	hashLock.RLock()
	firstVal := currentHash.Value
	hashLock.RUnlock()

	assert.NotEmpty(t, firstVal)

	time.Sleep(1100 * time.Millisecond)

	hashLock.RLock()
	assert.NotEqual(t, firstVal, currentHash.Value)
	hashLock.RUnlock()

	cancel()
	time.Sleep(100 * time.Millisecond)
}

func TestNotChangedBeforeTTL(t *testing.T) {
	hashLock := &sync.RWMutex{}
	ctrl := gomock.NewController(t)

	mock_logger := mocklog.NewMockLogger(ctrl)
	mock_logger.EXPECT().Debugf("%s: New UUID was generated '%s'", gomock.Any(), gomock.Any()).AnyTimes()
	mock_logger.EXPECT().Infof("GenerateHash worker stopped.").AnyTimes()

	currentHash := &generator.Hash{}
	ctx, cancel := context.WithCancel(context.Background())
	go generator.GenerateHash(currentHash, hashLock, mock_logger, 1, ctx)

	time.Sleep(1100 * time.Millisecond)

	hashLock.RLock()
	firstVal := currentHash.Value
	hashLock.RUnlock()

	time.Sleep(500 * time.Millisecond)

	hashLock.RLock()
	assert.Equal(t, firstVal, currentHash.Value)
	hashLock.RUnlock()

	cancel()
	time.Sleep(100 * time.Millisecond)
}

func TestGenerateWhileLock(t *testing.T) {
	hashLock := &sync.RWMutex{}
	ctrl := gomock.NewController(t)

	mock_logger := mocklog.NewMockLogger(ctrl)
	mock_logger.EXPECT().Debugf("%s: New UUID was generated '%s'", gomock.Any(), gomock.Any()).AnyTimes()
	mock_logger.EXPECT().Infof("GenerateHash worker stopped.").AnyTimes()

	currentHash := &generator.Hash{}
	ctx, cancel := context.WithCancel(context.Background())
	go generator.GenerateHash(currentHash, hashLock, mock_logger, 1, ctx)

	time.Sleep(1100 * time.Millisecond)

	hashLock.RLock()
	firstVal := currentHash.Value
	time.Sleep(1100 * time.Millisecond)
	hashLock.RUnlock()

	hashLock.RLock()
	assert.NotEqual(t, firstVal, currentHash.Value)
	hashLock.RUnlock()

	cancel()
	time.Sleep(100 * time.Millisecond)
}

func TestNotChangedWhileLock(t *testing.T) {
	hashLock := &sync.RWMutex{}
	ctrl := gomock.NewController(t)

	mock_logger := mocklog.NewMockLogger(ctrl)
	mock_logger.EXPECT().Debugf("%s: New UUID was generated '%s'", gomock.Any(), gomock.Any()).AnyTimes()
	mock_logger.EXPECT().Infof("GenerateHash worker stopped.").AnyTimes()

	currentHash := &generator.Hash{}
	ctx, cancel := context.WithCancel(context.Background())
	go generator.GenerateHash(currentHash, hashLock, mock_logger, 1, ctx)

	time.Sleep(1100 * time.Millisecond)

	hashLock.RLock()
	firstVal := currentHash.Value
	time.Sleep(1100 * time.Millisecond)
	assert.Equal(t, firstVal, currentHash.Value)
	hashLock.RUnlock()

	cancel()
	time.Sleep(100 * time.Millisecond)
}
