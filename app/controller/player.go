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

type PlayerController struct {
	Repo repos.Repository
}

// Add player resource routes
func (c *PlayerController) AddRoutes(r *mux.Router) {
	rAdminPlayer := r.PathPrefix("/player").Subrouter()
	rAdminPlayer.Use(middleware.Auth(c.Repo))
	rAdminPlayer.Use(middleware.Admin)
	rAdminPlayer.HandleFunc("", c.handlePostPlayer).Methods("POST")
	rAdminPlayer.HandleFunc("/{id}", c.handleDeletePlayer).Methods("DELETE")

	rAuthPlayer := r.PathPrefix("/player").Subrouter()
	rAuthPlayer.Use(middleware.Auth(c.Repo))
	rAuthPlayer.HandleFunc("/{id}", c.handleGetPlayer).Methods("GET")
	rAuthPlayer.HandleFunc("/{id}", c.handlePatchPlayer).Methods("PATCH")
}

// Handles GET requests to the user's resource
// @Summary Get a player
// @Description get a player by ID
// @Tags players
// @Accept  json
// @Produce  json
// @Param id path int true "Player ID"
// @Success 200 {object} model.User
// @Failure 401 {object} httputil.HTTPError
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /player/{id} [get]
func (c *PlayerController) handleGETPlayer(w http.ResponseWriter, req *http.Request) {
	player, err := c.getPlayerFromRequest(w, req)
	if err != nil {
		return
	}
	/// Add check if it's the team owner
	data, err := c.makePlayerJson(player)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeResponse(w, http.StatusOK, data)
}

// Handles a POST request to the player resource
// @Summary Create a player
// @Description create a player
// @Tags players
// @Accept  json
// @Produce  json
// @Param player body model.CreatePlayer true "Create player"
// @Success 200 {object} model.User
// @Failure 401 {object} httputil.HTTPError
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /player [post]
func (c *PlayerController) handlePOSTPlayer(w http.ResponseWriter, req *http.Request) {
	payload, err := c.getPlayerJsonFromRequest(w, req)
	if err != nil {
		return
	}

	player := models.Player{}
	c.fillPlayerData(&player, payload, true)

	err = c.Repo.Update(&player)
	if err != nil {
		log.Println(err)
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeResponse(w, http.StatusOK, []byte(fmt.Sprintf(`{"error": false, "id": %v}`, player.ID)))
}

// Handles a PUT request to the player resource
// @Summary Update a player
// @Description update a player
// @Tags players
// @Accept  json
// @Produce  json
// @Param player body model.UpdatePlayer true "Update player"
// @Param id path int true "Player ID"
// @Success 200 {object} models.Player
// @Failure 401 {object} httputil.HTTPError
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /player/{id} [put]
func (c *PlayerController) handlePUTPlayer(w http.ResponseWriter, req *http.Request) {
	payload, err1 := c.getPlayerJsonFromRequest(w, req)
	player, err2 := c.getPlayerFromRequest(w, req)
	user, err3 := getAuthenticatedUserFromRequest(w, req)
	if err1 != nil || err2 != nil || err3 != nil {
		return
	}
	isTeamOwner, isAdmin := player.Team.Owner.ID == user.ID, user.IsAdmin()
	if !isAdmin && !isTeamOwner {
		writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	c.fillPlayerData(&player, payload, isAdmin)

	err := c.Repo.Update(&player)
	if err != nil {
		log.Println(err)
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeResponse(w, http.StatusOK, noError())
}

// Handles a DELETE request to the player resource
// @Summary Delete a player
// @Description delete a player
// @Tags players
// @Accept  json
// @Produce  json
// @Param id path int true "Player ID"
// @Success 200
// @Failure 401 {object} httputil.HTTPError
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /player/{id} [delete]
func (c *PlayerController) handleDELETEPlayer(w http.ResponseWriter, req *http.Request) {
	player, err := c.getPlayerFromRequest(w, req)
	if err != nil {
		return
	}

	err = c.Repo.Delete(player)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
	}

	writeResponse(w, http.StatusOK, noError())
}

// Gets a Player model from the id in the request
func (c *PlayerController) getPlayerFromRequest(w http.ResponseWriter, req *http.Request) (models.Player, error) {
	id, err := parseIdFromRequest(w, req)
	player, err := c.Repo.GetPlayer(id)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, "Not found")
		return models.Player{}, err
	}
	return player, nil
}

type playerPayload struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Country     string `json:"country"`
	Age         int    `json:"age"`
	MarketValue int    `json:"market_value"`
	Position    int    `json:"position"`
	Team        int    `json:"team"`
}

// Returns a json from a player object
func (c *PlayerController) makePlayerJson(p models.Player) ([]byte, error) {
	return json.Marshal(playerPayload{
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		Country:     p.Country,
		Age:         p.Age,
		MarketValue: int(p.MarketValue),
		Team:        int(p.TeamID),
		Position:    p.Position,
	})
}

// Returns the player json from a request
func (c *PlayerController) getPlayerJsonFromRequest(w http.ResponseWriter, req *http.Request) (playerPayload, error) {
	decoder := json.NewDecoder(req.Body)
	var t playerPayload
	t.Team = -1
	err := decoder.Decode(&t)
	if err != nil {
		log.Println(err)
		writeErrorResponse(w, http.StatusBadRequest, "Incorrect body parameters")
		return t, err
	}
	return t, nil
}

// Fill player data from a json payload into a model
func (c *PlayerController) fillPlayerData(player *models.Player, payload playerPayload, isAdmin bool) {
	player.FirstName = payload.FirstName
	player.LastName = payload.LastName
	player.Country = payload.Country
	if isAdmin {
		if payload.Team >= 0 {
			player.TeamID = uint(payload.Team)
		}
		player.MarketValue = int32(payload.MarketValue)
		player.Age = payload.Age
		player.Position = payload.Position
	}
}
