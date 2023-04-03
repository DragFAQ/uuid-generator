package generator_test

import (
	mock_log "github.com/DragFAQ/uuid-generator/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
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

	if val == "" {
		t.Errorf("currentHash.Value is empty")
	}

	_, err := uuid.Parse(val)
	if err != nil {
		t.Errorf("currentHash.Value is not UUID")
	}

	if tim.IsZero() {
		t.Errorf("currentHash.GenerationTime is zero")
	}

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

	if firstVal == "" {
		t.Errorf("currentHash.Value is empty")
	}

	time.Sleep(1100 * time.Millisecond)

	hashLock.RLock()
	if firstVal == currentHash.Value {
		t.Errorf("currentHash.Value not changed")
	}
	hashLock.RUnlock()

	shutDownCh <- syscall.SIGTERM
	time.Sleep(100 * time.Millisecond)
}
