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

func GenerateHash(currentHash *Hash, hashLock *sync.RWMutex, logger log.Logger, ttlSeconds int, quit context.Context) {
	ticker := time.NewTicker(time.Duration(ttlSeconds) * time.Second)

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
			*currentHash = newHash
			logger.Debugf("%s: New UUID was generated '%s'", newHash.GenerationTime, newHash.Value)
			hashLock.Unlock()
		}
	}
}
