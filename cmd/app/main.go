package main

import (
	"log"
	"net/http"

	"github.com/bilyardvmetro/avito-PR-project/internal/handler"
	"github.com/bilyardvmetro/avito-PR-project/internal/repository/memory"
	"github.com/bilyardvmetro/avito-PR-project/internal/server"
	"github.com/bilyardvmetro/avito-PR-project/internal/service"
)

func main() {
	teamRepo := memory.NewTeamRepoMemory()
	userRepo := memory.NewUserRepoMemory()
	prRepo := memory.NewPRRepoMemory()

	teamSvc := service.NewTeamService(teamRepo, userRepo)
	userSvc := service.NewUserService(userRepo, prRepo)
	prSvc := service.NewPRService(prRepo, userRepo, teamRepo)

	teamHandler := handler.NewTeamHandler(teamSvc)
	userHandler := handler.NewUserHandler(userSvc)
	prHandler := handler.NewPRHandler(prSvc)

	srv := server.New(teamHandler, userHandler, prHandler)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", srv.Mux)
}
