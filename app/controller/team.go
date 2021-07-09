package controllers

import (
	"../middleware"
	"../models"
	"../repos"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
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
	rAdmin.HandleFunc("", c.handlePostTeam).Methods("POST")
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
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeResponse(w, http.StatusOK, data)
}

// Handles a POST request to a team resource
func (c *TeamController) handlePostTeam(w http.ResponseWriter, req *http.Request) {
	teamData, err := c.getTeamJson(w, req)
	if err != nil {
		return
	}

	team := models.Team{}
	c.fillTeamData(&team, teamData)
	team.Budget = teamData.Budget

	err = c.Repo.Create(&team)
	if err != nil {
		log.Println(err)
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeResponse(w, http.StatusOK, []byte(fmt.Sprintf(`{"error": false, "id": %v}`, team.ID)))
}

// Handles a DELETE request to a team resource
func (c *TeamController) handleDeleteTeam(w http.ResponseWriter, req *http.Request) {
	team, err := c.getTeamFromRequest(w, req)
	if err != nil {
		return
	}

	players := c.Repo.GetPlayers(team.ID)
	if len(players) > 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Cannot delete a team while it still has players")
		return
	}

	err = c.Repo.Delete(team)
	if err != nil {
		log.Println(err)
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeResponse(w, http.StatusOK, noError())
}

// Handles a PATCH request to a team resource
func (c *TeamController) handlePatchTeam(w http.ResponseWriter, req *http.Request) {
	user, err := getAuthenticatedUserFromRequest(w, req)
	team, err := c.getTeamFromRequest(w, req)
	if err != nil || (!user.IsAdmin() && !c.validateTeamOwner(w, user, team)) {
		return
	}

	t, err := c.getTeamJson(w, req)
	if err != nil {
		return
	}

	c.fillTeamData(&team, t)
	if user.IsAdmin() {
		team.Budget = t.Budget
	}

	err = c.Repo.Update(&team)
	if err != nil {
		log.Println(err)
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeResponse(w, http.StatusOK, noError())
}

type teamJson struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	MarketValue   int    `json:"market_value"`
	Budget  int    `json:"budget"`
	Players  []int    `json:"players"`
	Owner  int    `json:"owner"`
}

// Fill a team with data from a json structure
func (c *TeamController) fillTeamData(team *models.Team, t teamJson) {
	team.Country = t.Country
	team.Name = t.Name
	//if admin {
	//	team.OwnerID = t.Id
	//}
}

// Generate a json from a team model
func (c *TeamController) makeTeamJson(team models.Team) ([]byte, error) {
	players := make([]int, 0)
	marketValue := 0
	for _, p := range c.Repo.GetPlayers(team.ID) {
		players = append(players, int(p.ID))
		marketValue += int(p.MarketValue)
	}

	t := teamJson{
		Id:      team.ID,
		Name:    team.Name,
		Country: team.Country,
		MarketValue:   marketValue,
		Budget:  team.Budget,
		Players: players,
		Owner: int(team.OwnerID),
	}
	return json.Marshal(t)
}

func (c *TeamController) getTeamJson(w http.ResponseWriter, req *http.Request) (teamJson, error) {

	decoder := json.NewDecoder(req.Body)
	var t teamJson
	err := decoder.Decode(&t)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Incorrect body parameters")
		return t, err
	}
	return t, nil
}

func (c *TeamController) getTeamFromRequest(w http.ResponseWriter, req *http.Request) (models.Team, error) {
	id, err := parseIdFromRequest(w, req)
	team, err := c.Repo.GetTeam(id)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, "Not found")
		return models.Team{}, err
	}
	return team, nil
}

func (c *TeamController) validateTeamOwner(w http.ResponseWriter, user models.User, team models.Team) bool {
	if user.ID != team.OwnerID {
		writeErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return false
	}
	return true
}