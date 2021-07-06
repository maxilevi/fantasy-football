package controllers

import (
	"../middleware"
	"../models"
	"../repos"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type TeamController struct {
	Repo repos.Repository
}

func (c *TeamController) AddRoutes(r *mux.Router) {
	rAuth := r.PathPrefix("/team").Subrouter()
	rAuth.Use(middleware.Auth(c.Repo))
	rAuth.HandleFunc("/{id}", c.handleGetTeam).Methods("GET")
	rAuth.HandleFunc("/{id}", c.handlePatchTeam).Methods("PATCH")

	rAdmin := r.PathPrefix("/team").Subrouter()
	rAdmin.Use(middleware.Auth(c.Repo))
	rAdmin.Use(middleware.Admin)
	rAdmin.HandleFunc("/{id}", c.handlePostTeam).Methods("POST")
	rAdmin.HandleFunc("/{id}", c.handleDeleteTeam).Methods("DELETE")
}

// Handles a GET request to a team resource
func (c *TeamController) handleGetTeam(w http.ResponseWriter, req *http.Request) {
	team, err := c.getTeamFromRequest(w, req)
	if err != nil {
		return
	}

	data, err := c.makeTeamJson(team)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeResponse(w, http.StatusOK, data)
}

// Handles a PATCH request to a team resource
func (c *TeamController) handlePatchTeam(w http.ResponseWriter, req *http.Request) {
	user, err := getAuthenticatedUserFromRequest(w, req)
	team, err := c.getTeamFromRequest(w, req)
	if err != nil || !c.validateTeamOwner(w, user, team) {
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
	err = c.Repo.Update(team)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeResponse(w, http.StatusOK, noError())
}


func (c *TeamController) makeTeamJson(team models.Team) ([]byte, error) {
	type TeamJson struct {
		Id      uint   `json:"id"`
		Name    string `json:"name"`
		Country string `json:"country"`
		Value   int    `json:"value"`
		Budget  int    `json:"budget"`
		Players  []int    `json:"players"`
	}
	players := make([]int, 0)
	marketValue := 0
	for _, p := range c.Repo.GetPlayers(team.ID) {
		players = append(players, int(p.ID))
		marketValue += int(p.MarketValue)
	}

	t := TeamJson{
		Id:      team.ID,
		Name:    team.Name,
		Country: team.Country,
		Value:   marketValue,
		Budget:  team.Budget,
		Players: players,
	}
	return json.Marshal(t)
}

func (c *TeamController) getTeamFromRequest(w http.ResponseWriter, req *http.Request) (models.Team, error) {
	id, err := parseIdFromRequest(w, req)
	team, err := c.Repo.GetTeam(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Not found")
		return models.Team{}, err
	}
	return team, nil
}

func (c *TeamController) validateTeamOwner(w http.ResponseWriter, user models.User, team models.Team) bool {
	if user.ID != team.OwnerID {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return false
	}
	return true
}