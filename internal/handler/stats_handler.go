package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bilyardvmetro/avito-PR-project/internal/domain"
)

type StatsHandler struct {
	srv domain.StatsService
}

func NewStatsHandler(s domain.StatsService) *StatsHandler {
	return &StatsHandler{srv: s}
}

// GET /stats/assignments
func (h *StatsHandler) GetAssignments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.srv.GetAssignmentStats()
	if err != nil {
		// считаем ошибку здесь технической
		writeInternalErr(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(stats)
}
