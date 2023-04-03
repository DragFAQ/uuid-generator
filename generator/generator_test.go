package generator_test

import (
	mock_log "github.com/DragFAQ/uuid-generator/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/DragFAQ/uuid-generator/generator"
)

func TestGenerate(t *testing.T) {
	hashLock := &sync.RWMutex{}
	ctrl := gomock.NewController(t)

	mock_logger := mock_log.NewMockLogger(ctrl)
	mock_logger.EXPECT().Debugf("%s: New UUID was generated '%s'", gomock.Any(), gomock.Any()).AnyTimes()
	mock_logger.EXPECT().Infof("GenerateHash worker stopped.").AnyTimes()

	currentHash := &generator.Hash{}
	shutDownCh := make(chan os.Signal, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go generator.GenerateHash(currentHash, hashLock, mock_logger, 1, shutDownCh, &wg)

	time.Sleep(1100 * time.Millisecond)

	hashLock.RLock()
	val := currentHash.Value
	tim := currentHash.GenerationTime
	hashLock.RUnlock()

	assert.NotEmpty(t, val)

	_, err := uuid.Parse(val)
	assert.NoError(t, err)

	assert.Equal(t, tim.IsZero(), false)

	shutDownCh <- syscall.SIGTERM
	time.Sleep(100 * time.Millisecond)
}

func TestChangedValAfterTTL(t *testing.T) {
	hashLock := &sync.RWMutex{}
	ctrl := gomock.NewController(t)

	mock_logger := mock_log.NewMockLogger(ctrl)
	mock_logger.EXPECT().Debugf("%s: New UUID was generated '%s'", gomock.Any(), gomock.Any()).AnyTimes()
	mock_logger.EXPECT().Infof("GenerateHash worker stopped.").AnyTimes()

	currentHash := &generator.Hash{}
	shutDownCh := make(chan os.Signal, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go generator.GenerateHash(currentHash, hashLock, mock_logger, 1, shutDownCh, &wg)

	time.Sleep(1100 * time.Millisecond)

	hashLock.RLock()
	firstVal := currentHash.Value
	hashLock.RUnlock()

	assert.NotEmpty(t, firstVal)

	time.Sleep(1100 * time.Millisecond)

	hashLock.RLock()
	assert.NotEqual(t, firstVal, currentHash.Value)
	hashLock.RUnlock()

	shutDownCh <- syscall.SIGTERM
	time.Sleep(100 * time.Millisecond)
}

func TestNotChangedBeforeTTL(t *testing.T) {
	hashLock := &sync.RWMutex{}
	ctrl := gomock.NewController(t)

	mock_logger := mock_log.NewMockLogger(ctrl)
	mock_logger.EXPECT().Debugf("%s: New UUID was generated '%s'", gomock.Any(), gomock.Any()).AnyTimes()
	mock_logger.EXPECT().Infof("GenerateHash worker stopped.").AnyTimes()

	currentHash := &generator.Hash{}
	shutDownCh := make(chan os.Signal, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go generator.GenerateHash(currentHash, hashLock, mock_logger, 1, shutDownCh, &wg)

	time.Sleep(1100 * time.Millisecond)

	hashLock.RLock()
	firstVal := currentHash.Value
	hashLock.RUnlock()

	time.Sleep(500 * time.Millisecond)

	hashLock.RLock()
	assert.Equal(t, firstVal, currentHash.Value)
	hashLock.RUnlock()

	shutDownCh <- syscall.SIGTERM
	time.Sleep(100 * time.Millisecond)
}

func TestGenerateWhileLock(t *testing.T) {
	hashLock := &sync.RWMutex{}
	ctrl := gomock.NewController(t)

	mock_logger := mock_log.NewMockLogger(ctrl)
	mock_logger.EXPECT().Debugf("%s: New UUID was generated '%s'", gomock.Any(), gomock.Any()).AnyTimes()
	mock_logger.EXPECT().Infof("GenerateHash worker stopped.").AnyTimes()

	currentHash := &generator.Hash{}
	shutDownCh := make(chan os.Signal, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go generator.GenerateHash(currentHash, hashLock, mock_logger, 1, shutDownCh, &wg)

	time.Sleep(1100 * time.Millisecond)

	hashLock.RLock()
	firstVal := currentHash.Value
	time.Sleep(1100 * time.Millisecond)
	hashLock.RUnlock()

	hashLock.RLock()
	assert.NotEqual(t, firstVal, currentHash.Value)
	hashLock.RUnlock()

	shutDownCh <- syscall.SIGTERM
	time.Sleep(100 * time.Millisecond)
}

func TestNotChangedWhileLock(t *testing.T) {
	hashLock := &sync.RWMutex{}
	ctrl := gomock.NewController(t)

	mock_logger := mock_log.NewMockLogger(ctrl)
	mock_logger.EXPECT().Debugf("%s: New UUID was generated '%s'", gomock.Any(), gomock.Any()).AnyTimes()
	mock_logger.EXPECT().Infof("GenerateHash worker stopped.").AnyTimes()

	currentHash := &generator.Hash{}
	shutDownCh := make(chan os.Signal, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go generator.GenerateHash(currentHash, hashLock, mock_logger, 1, shutDownCh, &wg)

	time.Sleep(1100 * time.Millisecond)

	hashLock.RLock()
	firstVal := currentHash.Value
	time.Sleep(1100 * time.Millisecond)
	assert.Equal(t, firstVal, currentHash.Value)
	hashLock.RUnlock()

	shutDownCh <- syscall.SIGTERM
	time.Sleep(100 * time.Millisecond)
}
