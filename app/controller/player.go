package controller

import (
	"../httputil"
	"../models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Handles GET requests to the players resource
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
// @Router /players/{id} [get]
func (c *Controller) ShowPlayer(ctx *gin.Context) {
	player, err := c.getPlayerFromRequest(ctx)
	if err != nil {
		return
	}

	payload := c.getPlayerPayload(player)

	httputil.NoError(ctx, payload)
}

// Handles a PATCH request to the player resource
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
// @Router /players/{id} [patch]
// @Security BearerAuth
func (c *Controller) UpdatePlayer(ctx *gin.Context) {
	player, err2 := c.getPlayerFromRequest(ctx)
	payload := c.fillDefaultPlayerPayload(player)
	err1 := ctx.BindJSON(&payload)
	user, err3 := c.getAuthenticatedUserFromRequest(ctx)

	if err1 != nil || err2 != nil || err3 != nil {
		log.Println(err1, err2, err3)
		if err1 != nil {
			httputil.NewError(ctx, http.StatusBadRequest, "Incorrect body parameters")
		}
		return
	}

	isTeamOwner, isAdmin := player.Team.UserID == user.ID, user.IsAdmin()
	if !isAdmin && !isTeamOwner {
		httputil.NewError(ctx, http.StatusUnauthorized, "Only administrators or owners can edit players")
		return
	}

	player.FirstName = payload.FirstName
	player.LastName = payload.LastName
	player.Country = payload.Country
	if isAdmin {
		team, err := c.Repo.GetTeam(uint(payload.Team))
		if err != nil {
			httputil.NewError(ctx, http.StatusNotFound, "Team not found")
			return
		}
		player.TeamID = uint(payload.Team)
		player.Team = team
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
// @Router /players/{id} [delete]
// @Security BearerAuth
func (c *Controller) DeletePlayer(ctx *gin.Context) {
	player, err := c.getPlayerFromRequest(ctx)
	if err != nil {
		return
	}

	err = c.Repo.DeletePlayer(&player)
	if err != nil {
		log.Println(err)
		httputil.NewError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	httputil.NoErrorEmpty(ctx)
}

// Gets a Player model from the id in the request
func (c *Controller) getPlayerFromRequest(ctx *gin.Context) (models.Player, error) {
	id, err := c.parseIdFromRequest(ctx, "playerId")
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

// Fill the player payload with default values
func (c* Controller) fillDefaultPlayerPayload(player models.Player) models.UpdatePlayer {
	var payload models.UpdatePlayer
	payload.Team = int(player.TeamID)
	payload.FirstName = player.FirstName
	payload.LastName = player.LastName
	payload.Country = player.Country
	payload.Age = player.Age
	payload.MarketValue = player.MarketValue
	payload.Position = player.Position
	return payload
}

// Create and fill the show player payload
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
