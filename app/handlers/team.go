package handlers

import (
	"../middleware"
	"../models"
	"../repos"
	"encoding/json"
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func AddTeamRoutes(r *mux.Router, repo repos.Repository) {
	rTeam := r.PathPrefix("/team").Subrouter()
	rTeam.Use(middleware.Auth(repo))
	rTeam.HandleFunc("/{id}", wrap(handleGetTeam, repo)).Methods( "GET")
	rTeam.HandleFunc("/{id}", wrap(handlePatchTeam, repo)).Methods( "PATCH")
}

func handleGetTeam(w http.ResponseWriter, req *http.Request, repo repos.Repository) {
	team, err := getTeamFromRequest(w, req, repo)
	if err != nil {
		return
	}

	data, err := makeTeamJson(team, repo)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeResponse(w, http.StatusOK, data)
}


func makeTeamJson(team models.Team, repo repos.Repository) ([]byte, error) {
	type TeamJson struct {
		Id uint `json:"id"`
		Name string `json:"name"`
		Country string `json:"country"`
		Value int `json:"value"`
		Budget int `json:"budget"`
	}
	t := TeamJson{
		Id: team.ID,
		Name: team.Name,
		Country: team.Country,
		Value: int(teamMarketValue(team, repo)),
		Budget: team.Budget,
	}
	return json.Marshal(t)
}

func teamMarketValue(team models.Team, repo repos.Repository) int32 {
	var sum int32
	for _, player := range repo.GetPlayers(team.ID) {
		sum += player.MarketValue
	}
	return sum
}

type patchTeamData struct {
	country string
	name string
	budget int
}

func handlePatchTeam(w http.ResponseWriter, req *http.Request, repo repos.Repository) {
	user, err := getUserFromRequest(w, req)
	team, err := getTeamFromRequest(w, req, repo)
	if err != nil || !validateTeamOwner(w, user, team) {
		return
	}

	decoder := json.NewDecoder(req.Body)
	var t patchTeamData
	err = decoder.Decode(&t)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Incorrect body parameters")
		return
	}

	team.Country = t.country
	team.Name = t.name
	if user.IsAdmin() {
		team.Budget = t.budget
	}
	repo.UpdateTeam(team)
	writeResponse(w, http.StatusOK, []byte(`{"error": false}`))
}

func getUserFromRequest(w http.ResponseWriter, req *http.Request) (models.User, error) {
	val, ok := context.GetOk(req, "user")
	if !ok {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return models.User{}, fmt.Errorf("internal server error")
	}
	user, ok := val.(models.User)
	if !ok {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return models.User{}, fmt.Errorf("internal server error")
	}
	return user, nil
}

func getTeamFromRequest(w http.ResponseWriter, req *http.Request, repo repos.Repository) (models.Team, error) {
	id, err := parseTeamId(w, req)
	team, err := repo.GetTeam(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Not found")
		return models.Team{}, err
	}
	return team, nil
}

func validateTeamOwner(w http.ResponseWriter, user models.User, team models.Team) bool {
	if user.ID != team.OwnerID {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return false
	}
	return true
}

func parseTeamId(w http.ResponseWriter, req *http.Request) (uint, error) {
	vars := mux.Vars(req)
	id, ok := strconv.ParseInt(vars["id"], 10, 32)
	if ok != nil {
		writeError(w, http.StatusBadRequest, "Invalid team id")
		return 0, fmt.Errorf("invalid team id")
	}
	return uint(id), nil
}