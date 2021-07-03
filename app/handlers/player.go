package handlers

import (
	"../middleware"
	"../models"
	"../repos"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// Add player resource routes
func AddPlayerRoutes(r *mux.Router, repo repos.Repository) {
	rAdminPlayer := r.PathPrefix("/player").Subrouter()
	rAdminPlayer.Use(middleware.Auth(repo))
	rAdminPlayer.Use(middleware.Admin)
	rAdminPlayer.HandleFunc("", wrap(handlePostPlayer, repo)).Methods("POST")
	rAdminPlayer.HandleFunc("/{id}", wrap(handleDeletePlayer, repo)).Methods("DELETE")

	rAuthPlayer := r.PathPrefix("/player").Subrouter()
	rAuthPlayer.Use(middleware.Auth(repo))
	rAuthPlayer.HandleFunc("/{id}", wrap(handleGetPlayer, repo)).Methods("GET")
	rAuthPlayer.HandleFunc("/{id}", wrap(handlePatchPlayer, repo)).Methods("PATCH")
}

// Handles a GET request to the player resource
func handleGetPlayer(w http.ResponseWriter, req *http.Request, repo repos.Repository) {
	player, err := getPlayerFromRequest(w, req, repo)
	if err != nil {
		return
	}
	/// Add check if it's the team owner
	data, err := makePlayerJson(player)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeResponse(w, http.StatusOK, data)
}

// Handles a POST request to the player resource
func handlePostPlayer(w http.ResponseWriter, req *http.Request, repo repos.Repository) {
	payload, err := getPlayerJsonFromRequest(w, req)
	if err != nil {
		return
	}

	player := models.Player{}
	fillPlayerData(&player, payload)

	err = repo.UpdatePlayer(player)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
	}

	writeResponse(w, http.StatusOK, []byte(`{"error": false}`))
}

// Handles a PATCH request to the player resource
func handlePatchPlayer(w http.ResponseWriter, req *http.Request, repo repos.Repository) {
	payload, err1 := getPlayerJsonFromRequest(w, req)
	player, err2 := getPlayerFromRequest(w, req, repo)
	if err1 != nil || err2 != nil {
		return
	}

	fillPlayerData(&player, payload)

	err := repo.UpdatePlayer(player)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
	}

	writeResponse(w, http.StatusOK, []byte(`{"error": false}`))
}

// Handles a DELETE request to the player resource
func handleDeletePlayer(w http.ResponseWriter, req *http.Request, repo repos.Repository) {
	player, err := getPlayerFromRequest(w, req, repo)
	if err != nil {
		return
	}

	err = repo.DeletePlayer(player)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
	}

	writeResponse(w, http.StatusOK, []byte(`{"error": false}`))
}

// Gets a Player model from the id in the request
func getPlayerFromRequest(w http.ResponseWriter, req *http.Request, repo repos.Repository) (models.Player, error) {
	id, err := parseIdFromRequest(w, req)
	player, err := repo.GetPlayer(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Not found")
		return models.Player{}, err
	}
	return player, nil
}

type playerPayload struct {
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Country string `json:"country"`
	Age int `json:"age"`
	MarketValue int `json:"market_value"`
	Team int `json:"team"`
}

func makePlayerJson(p models.Player) ([]byte, error) {
	return json.Marshal(playerPayload{
		FirstName: p.FirstName,
		LastName: p.LastName,
		Country: p.Country,
		Age: p.Age,
		MarketValue: int(p.MarketValue),
		Team: int(p.TeamID),
	})
}

func getPlayerJsonFromRequest(w http.ResponseWriter, req *http.Request) (playerPayload, error) {
	decoder := json.NewDecoder(req.Body)
	var t playerPayload
	err := decoder.Decode(&t)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Incorrect body parameters")
		return t, err
	}
	return t, nil
}

func fillPlayerData(player *models.Player, payload playerPayload) {
	player.FirstName = payload.FirstName
	player.LastName = payload.LastName
	player.Country = payload.Country
	player.TeamID = uint(payload.Team)
	player.MarketValue = int32(payload.MarketValue)
	player.Age = payload.Age
}