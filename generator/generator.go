package generator

import (
	log "github.com/DragFAQ/uuid-generator/logger"
	"github.com/google/uuid"
	"sync"
	"time"
)

type Hash struct {
	Value          string
	GenerationTime time.Time
}

func GenerateHash(currentHash *Hash, hashLock *sync.RWMutex, logger log.Logger, ttlSeconds int) {
	for {
		time.Sleep(time.Duration(ttlSeconds) * time.Second)
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
