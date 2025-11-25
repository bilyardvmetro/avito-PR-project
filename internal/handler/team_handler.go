package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bilyardvmetro/avito-PR-project/internal/domain"
)

type TeamHandler struct {
	srv domain.TeamService
}

func NewTeamHandler(s domain.TeamService) *TeamHandler {
	return &TeamHandler{srv: s}
}

func (h *TeamHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req domain.Team
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	team, err := h.srv.AddTeam(req)
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
		"team": team,
	})
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		writeDomainErr(w, domain.NewError(domain.ErrorNotFound, "team not found"))
		return
	}

	team, err := h.srv.GetTeam(teamName)
	if err != nil {
		if derr, ok := err.(*domain.DomainError); ok {
			writeDomainErr(w, derr)
		} else {
			writeInternalErr(w)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(team)
}
