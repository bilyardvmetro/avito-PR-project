package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

func main() {
	baseURL := "http://localhost:8080"
	client := &http.Client{Timeout: 5 * time.Second}

	log.Printf("Seeding data to %s ...", baseURL)

	// 20 команд по 10 пользователей = 200 пользователей
	totalTeams := 20
	usersPerTeam := 10

	userCounter := 1

	for i := 1; i <= totalTeams; i++ {
		teamName := fmt.Sprintf("team-%02d", i)
		members := make([]TeamMember, 0, usersPerTeam)

		for j := 1; j <= usersPerTeam; j++ {
			userID := fmt.Sprintf("u%d", userCounter)
			username := fmt.Sprintf("User-%d", userCounter)
			members = append(members, TeamMember{
				UserID:   userID,
				Username: username,
				IsActive: true,
			})
			userCounter++
		}

		team := Team{
			TeamName: teamName,
			Members:  members,
		}

		body, _ := json.Marshal(team)
		req, err := http.NewRequest(http.MethodPost, baseURL+"/team/add", bytes.NewReader(body))
		if err != nil {
			log.Fatalf("build request failed: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("seed team %s failed: %v", teamName, err)
		}
		_ = resp.Body.Close()

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusConflict {
			log.Fatalf("team %s: unexpected status %d", teamName, resp.StatusCode)
		}

		log.Printf("Team %s seeded (status %d)", teamName, resp.StatusCode)
	}

	// по 5 PR на каждую команду от первых трёх пользователей
	prID := 1
	for i := 1; i <= totalTeams; i++ {
		teamName := fmt.Sprintf("team-%02d", i)
		for k := 1; k <= 5; k++ {
			authorUserID := fmt.Sprintf("u%d", (i-1)*usersPerTeam+1)

			reqBody := map[string]string{
				"pull_request_id":   fmt.Sprintf("pr-%d", prID),
				"pull_request_name": fmt.Sprintf("Seed PR %d in %s", prID, teamName),
				"author_id":         authorUserID,
			}
			body, _ := json.Marshal(reqBody)

			req, err := http.NewRequest(http.MethodPost, baseURL+"/pullRequest/create", bytes.NewReader(body))
			if err != nil {
				log.Fatalf("build PR request failed: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				log.Fatalf("seed PR %d failed: %v", prID, err)
			}
			_ = resp.Body.Close()

			if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusConflict {
				log.Fatalf("PR %d: unexpected status %d", prID, resp.StatusCode)
			}
			log.Printf("PR %d seeded (status %d)", prID, resp.StatusCode)
			prID++
		}
	}

	log.Println("Seeding finished.")
}
