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
	rAuth := r.PathPrefix("/team").Subrouter()
	rAuth.Use(middleware.Auth(repo))
	rAuth.HandleFunc("/{id}", wrap(handleGetTeam, repo)).Methods("GET")
	rAuth.HandleFunc("/{id}", wrap(handlePatchTeam, repo)).Methods("PATCH")

	rAdmin := r.PathPrefix("/team").Subrouter()
	rAdmin.Use(middleware.Auth(repo))
	rAdmin.Use(middleware.Admin)
	rAdmin.HandleFunc("/{id}", wrap(handlePostTeam, repo)).Methods( "POST")
	rAdmin.HandleFunc("/{id}", wrap(handleDeleteTeam, repo)).Methods( "DELETE")
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
		Id      uint   `json:"id"`
		Name    string `json:"name"`
		Country string `json:"country"`
		Value   int    `json:"value"`
		Budget  int    `json:"budget"`
	}
	t := TeamJson{
		Id:      team.ID,
		Name:    team.Name,
		Country: team.Country,
		Value:   int(teamMarketValue(team, repo)),
		Budget:  team.Budget,
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

func handlePatchTeam(w http.ResponseWriter, req *http.Request, repo repos.Repository) {
	user, err := getUserFromRequest(w, req)
	team, err := getTeamFromRequest(w, req, repo)
	if err != nil || !validateTeamOwner(w, user, team) {
		return
	}

	type patchTeamData struct {
		Country string `json:"country"`
		Name    string `json:"name"`
		Budget  int    `json:"budget"`
	}

	decoder := json.NewDecoder(req.Body)
	var t patchTeamData
	err = decoder.Decode(&t)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Incorrect body parameters")
		return
	}

	team.Country = t.Country
	team.Name = t.Name
	if user.IsAdmin() {
		team.Budget = t.Budget
	}
	err = repo.UpdateTeam(team)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
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
	id, err := parseIdFromRequest(w, req)
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

func parseIdFromRequest(w http.ResponseWriter, req *http.Request) (uint, error) {
	vars := mux.Vars(req)
	id, ok := strconv.ParseInt(vars["id"], 10, 32)
	if ok != nil {
		writeError(w, http.StatusBadRequest, "Invalid id")
		return 0, fmt.Errorf("invalid id")
	}
	return uint(id), nil
}
