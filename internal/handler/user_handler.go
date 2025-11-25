package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bilyardvmetro/avito-PR-project/internal/domain"
)

type UserHandler struct {
	srv domain.UserService
}

func NewUserHandler(s domain.UserService) *UserHandler {
	return &UserHandler{srv: s}
}

// POST /users/setIsActive
// body: { "user_id": "...", "is_active": true }
func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.UserID == "" {
		writeDomainErr(w, domain.NewError(domain.ErrorNotFound, "user not found"))
		return
	}

	user, err := h.srv.SetIsActive(req.UserID, req.IsActive)
	if err != nil {
		if derr, ok := err.(*domain.DomainError); ok {
			writeDomainErr(w, derr)
		} else {
			writeInternalErr(w)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"user": user,
	})
}

// GET /users/getReview?user_id=...
func (h *UserHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	// даже если user_id пустой — по сути просто вернём пустой список

	prs, err := h.srv.GetAssignedPRs(userID)
	if err != nil {
		if derr, ok := err.(*domain.DomainError); ok {
			writeDomainErr(w, derr)
		} else {
			writeInternalErr(w)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"user_id":       userID,
		"pull_requests": prs,
	})
}
