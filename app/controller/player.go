package controller

import (
	"../httputil"
	"../models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Handles GET requests to the user's resource
// @Summary Show a player
// @Description Get a player by ID
// @Tags Players
// @Accept  json
// @Produce  json
// @Param id path int true "Player ID"
// @Success 200 {object} models.ShowPlayer
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /player/{id} [get]
func (c *Controller) ShowPlayer(ctx *gin.Context) {
	player, err := c.getPlayerFromRequest(ctx)
	if err != nil {
		return
	}

	payload := models.ShowPlayer{
		BasePlayer: models.BasePlayer{
			FirstName:   player.FirstName,
			LastName:    player.LastName,
			MarketValue: player.MarketValue,
			Age:         player.Age,
			Country:     player.Country,
			Position:    player.Position,
		},
	}

	httputil.NoError(ctx, payload)
}

// Handles a POST request to the player resource
// @Summary Create a player
// @Description Create a player
// @Tags Players
// @Accept  json
// @Produce  json
// @Param player body models.CreatePlayer true "Create player"
// @Success 200
// @Failure 401 {object} httputil.HTTPError
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /player [post]
// @Security BearerAuth
func (c *Controller) CreatePlayer(ctx *gin.Context) {
	var payload models.CreatePlayer
	err := ctx.BindJSON(&payload)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, "Incorrect body parameters")
		return
	}

	player := models.Player{
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		Country:     payload.Country,
		Age:         payload.Age,
		MarketValue: payload.MarketValue,
		Position:    payload.Position,
		TeamID:      payload.Team,
	}

	if _, err = c.Repo.GetTeam(payload.Team); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, "An invalid team id was provided")
		return
	}

	err = c.Repo.Update(&player)
	if err != nil {
		log.Println(err)
		httputil.NewError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	httputil.NoError(ctx, map[string]interface{}{
		"id": player.ID,
	})
}

// Handles a PUT request to the player resource
// @Summary Update a player
// @Description Update a player
// @Tags Players
// @Accept  json
// @Produce  json
// @Param player body models.UpdatePlayer true "Update player"
// @Param id path int true "Player ID"
// @Success 200
// @Failure 401 {object} httputil.HTTPError
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /player/{id} [patch]
// @Security BearerAuth
func (c *Controller) UpdatePlayer(ctx *gin.Context) {
	var payload models.UpdatePlayer
	payload.Team = -1
	err1 := ctx.BindJSON(payload)
	player, err2 := c.getPlayerFromRequest(ctx)
	user, err3 := c.getAuthenticatedUserFromRequest(ctx)
	if err1 != nil || err2 != nil || err3 != nil {
		if err1 != nil {
			httputil.NewError(ctx, http.StatusBadRequest, "Incorrect body parameters")
		}
		return
	}

	isTeamOwner, isAdmin := player.Team.Owner.ID == user.ID, user.IsAdmin()
	if !isAdmin && !isTeamOwner {
		httputil.NewError(ctx, http.StatusUnauthorized, "Only administrators or owners can edit players")
		return
	}

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

	err := c.Repo.Update(&player)
	if err != nil {
		log.Println(err)
		httputil.NewError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	httputil.NoErrorEmpty(ctx)
}

// Handles a DELETE request to the player resource
// @Summary Delete a player
// @Description Deletes a player
// @Tags Players
// @Accept  json
// @Produce  json
// @Param id path int true "Player ID"
// @Success 200
// @Failure 401 {object} httputil.HTTPError
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /player/{id} [delete]
// @Security BearerAuth
func (c *Controller) DeletePlayer(ctx *gin.Context) {
	player, err := c.getPlayerFromRequest(ctx)
	if err != nil {
		return
	}

	err = c.Repo.Delete(player)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	httputil.NoErrorEmpty(ctx)
}

// Gets a Player model from the id in the request
func (c *Controller) getPlayerFromRequest(ctx *gin.Context) (models.Player, error) {
	id, err := c.parseIdFromRequest(ctx)
	if err != nil {
		return models.Player{}, err
	}

	player, err := c.Repo.GetPlayer(id)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, "Failed to find player")
		return models.Player{}, err
	}

	return player, nil
}

func (c *Controller) getPlayerPayload(p models.Player) models.ShowPlayer {
	return models.ShowPlayer{
		ID: p.ID,
		BasePlayer: models.BasePlayer{
			FirstName:   p.FirstName,
			LastName:    p.LastName,
			Country:     p.Country,
			Age:         p.Age,
			MarketValue: p.MarketValue,
			Position:    p.Position,
		},
	}
}
