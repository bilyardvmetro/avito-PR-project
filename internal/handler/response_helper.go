package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bilyardvmetro/avito-PR-project/internal/domain"
)

func writeDomainErr(w http.ResponseWriter, err *domain.DomainError) {
	w.Header().Set("Content-Type", "application/json")
	status := domainErrToCode(err.Code)
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]any{
			"code":    err.Code,
			"message": err.Message,
		},
	})
}

func writeInternalErr(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]any{
			"code":    "INTERNAL",
			"message": "internal server error",
		},
	})
}

func domainErrToCode(code domain.ErrorCode) int {
	switch code {
	case domain.ErrorNotFound:
		return http.StatusNotFound

	case domain.ErrorTeamExists,
		domain.ErrorPRExists,
		domain.ErrorPRMerged,
		domain.ErrorNotAssigned,
		domain.ErrorNoCandidate:
		return http.StatusConflict

	default:
		return http.StatusBadRequest
	}
}
