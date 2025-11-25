package server

import (
	"net/http"

	"github.com/bilyardvmetro/avito-PR-project/internal/handler"
)

type Server struct {
	Mux *http.ServeMux
}

func New(
	teamHandler *handler.TeamHandler,
	userHandler *handler.UserHandler,
	prHandler *handler.PRHandler,
	statsHandler *handler.StatsHandler,
) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/team/add", teamHandler.AddTeam)
	mux.HandleFunc("/team/get", teamHandler.GetTeam)

	mux.HandleFunc("/users/setIsActive", userHandler.SetIsActive)
	mux.HandleFunc("/users/getReview", userHandler.GetReview)

	mux.HandleFunc("/pullRequest/create", prHandler.CreatePR)
	mux.HandleFunc("/pullRequest/merge", prHandler.MergePR)
	mux.HandleFunc("/pullRequest/reassign", prHandler.Reassign)

	// эндпоинт статистики
	mux.HandleFunc("/stats/assignments", statsHandler.GetAssignments)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	return &Server{Mux: mux}
}
