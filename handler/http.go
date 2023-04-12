package handler

import (
	"encoding/json"
	"net/http"
	"time"

	hash "github.com/DragFAQ/uuid-generator/generator"
	log "github.com/DragFAQ/uuid-generator/logger"
)

type HttpHandler struct {
	logger log.Logger
}

func NewHttpHandler(logger log.Logger) *HttpHandler {
	return &HttpHandler{
		logger: logger,
	}
}

func (h *HttpHandler) GetCurrentHash(w http.ResponseWriter, _ *http.Request) {
	currentHash := hash.GetHash()
	resp := map[string]string{
		"hash":            currentHash.Value,
		"generation_time": currentHash.GenerationTime.Format(time.RFC3339),
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
