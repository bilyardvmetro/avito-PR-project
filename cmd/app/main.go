package main

import (
	"context"
	"log"
	"net/http"

	"github.com/bilyardvmetro/avito-PR-project/internal/handler"
	pg "github.com/bilyardvmetro/avito-PR-project/internal/repository/postgres"
	"github.com/bilyardvmetro/avito-PR-project/internal/server"
	"github.com/bilyardvmetro/avito-PR-project/internal/service"
)

func main() {
	ctx := context.Background()
	db, err := pg.Connect(ctx)
	if err != nil {
		log.Fatal("DB connect error:", err)
	}

	teamRepo := pg.NewTeamRepoPostgres(db)
	userRepo := pg.NewUserRepoPostgres(db)
	prRepo := pg.NewPRRepoPostgres(db)

	teamSvc := service.NewTeamService(teamRepo, userRepo)
	userSvc := service.NewUserService(userRepo, prRepo)
	prSvc := service.NewPRService(prRepo, userRepo, teamRepo)

	teamH := handler.NewTeamHandler(teamSvc)
	userH := handler.NewUserHandler(userSvc)
	prH := handler.NewPRHandler(prSvc)

	srv := server.New(teamH, userH, prH)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", srv.Mux); err != nil {
		log.Fatal(err)
	}
}
