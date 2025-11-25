package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bilyardvmetro/avito-PR-project/internal/domain"
)

type PRHandler struct {
	srv domain.PRService
}

func NewPRHandler(s domain.PRService) *PRHandler {
	return &PRHandler{srv: s}
}

// POST /pullRequest/create
// body: { pull_request_id, pull_request_name, author_id }
func (h *PRHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PullRequestID   string `json:"pull_request_id"`
		PullRequestName string `json:"pull_request_name"`
		AuthorID        string `json:"author_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.PullRequestID == "" || req.PullRequestName == "" || req.AuthorID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pr, err := h.srv.CreatePR(req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		if derr, ok := err.(*domain.DomainError); ok {
			writeDomainErr(w, derr)
		} else {
			writeInternalErr(w)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"pr": pr,
	})
}

// POST /pullRequest/merge
// body: { pull_request_id }
func (h *PRHandler) MergePR(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PullRequestID string `json:"pull_request_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.PullRequestID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pr, err := h.srv.MergePR(req.PullRequestID)
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
		"pr": pr,
	})
}

// POST /pullRequest/reassign
// body: { pull_request_id, old_user_id }
func (h *PRHandler) Reassign(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PullRequestID string `json:"pull_request_id"`
		OldUserID     string `json:"old_user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.PullRequestID == "" || req.OldUserID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pr, newReviewer, err := h.srv.ReassignReviewer(req.PullRequestID, req.OldUserID)
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
		"pr":          pr,
		"replaced_by": newReviewer,
	})
}
