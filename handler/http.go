package handler

import (
	"encoding/json"
	hash "github.com/DragFAQ/uuid-generator/generator"
	log "github.com/DragFAQ/uuid-generator/logger"
	"net/http"
	"sync"
	"time"
)

type HttpHandler struct {
	logger      log.Logger
	currentHash *hash.Hash
	hashLock    *sync.RWMutex
}

func NewHttpHandler(currentHash *hash.Hash, hashLock *sync.RWMutex, logger log.Logger) *HttpHandler {
	return &HttpHandler{
		logger:      logger,
		currentHash: currentHash,
		hashLock:    hashLock,
	}
}

func (h *HttpHandler) GetCurrentHash(w http.ResponseWriter, r *http.Request) {
	h.hashLock.RLock()
	defer h.hashLock.RUnlock()

	resp := map[string]string{
		"hash":            h.currentHash.Value,
		"generation_time": h.currentHash.GenerationTime.Format(time.RFC3339),
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}
