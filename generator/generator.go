package generator

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	log "github.com/DragFAQ/uuid-generator/logger"
)

type Hash struct {
	Value          string
	GenerationTime time.Time
}

var (
	currentHash Hash
	hashLock    sync.RWMutex
)

func GetHash() Hash {
	hashLock.RLock()
	defer hashLock.RUnlock()

	return currentHash
}

func GenerateHash(quit context.Context, logger log.Logger, ttlSeconds int) {
	ticker := time.NewTicker(time.Duration(ttlSeconds) * time.Second)
	currentHash = Hash{
		Value:          uuid.New().String(),
		GenerationTime: time.Now(),
	}

	for {
		select {
		case <-quit.Done():
			logger.Infof("GenerateHash worker stopped.")
			return
		case <-ticker.C:
			newHash := Hash{
				Value:          uuid.New().String(),
				GenerationTime: time.Now(),
			}
			hashLock.Lock()
			currentHash = newHash
			logger.Debugf("%s: New UUID was generated '%s'", newHash.GenerationTime, newHash.Value)
			hashLock.Unlock()
		}
	}
}
