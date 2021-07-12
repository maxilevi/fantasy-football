package controller

import (
	"../httputil"
	"../models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// @Summary List players of a team
// @Description List all the players of a team
// @Tags Teams
// @Accept  json
// @Produce  json
// @Param id path int true "Team ID"
// @Success 200 {array} models.ShowPlayer
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /teams/{id}/players [get]
func (c *Controller) ListTeamPlayers(ctx *gin.Context) {
	team, err := c.getTeamFromRequest(ctx)
	if err != nil {
		return
	}

	players := c.Repo.GetPlayers(team.ID)

	playerModels := make([]models.ShowPlayer, 0)
	for _, p := range players {
		playerModels = append(playerModels, c.getPlayerPayload(p))
	}

	httputil.NoError(ctx, gin.H{"players": playerModels})
}

// Handles a GET request to a team resource
// @Summary Get a team
// @Description Get team by ID
// @Tags Teams
// @Accept  json
// @Produce  json
// @Param id path int true "Team ID"
// @Success 200 {object} models.ShowTeam
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /teams/{id} [get]
func (c *Controller) ShowTeam(ctx *gin.Context) {
	team, err := c.getTeamFromRequest(ctx)
	if err != nil {
		return
	}

	httputil.NoError(ctx, c.getTeamPayload(team))
}

// Handles a POST request to a team resource
// @Summary Post a team
// @Description Create a new team
// @Tags Teams
// @Accept  json
// @Produce  json
// @Success 200
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /teams [post]
// @Security BearerAuth
func (c *Controller) CreateTeam(ctx *gin.Context) {
	var t models.CreateTeam
	t.Owner = -1
	err := ctx.BindJSON(&t)
	owner := t.Owner
	if _, ok := ctx.Get("TeamOwner"); owner == -1 && ok {
		owner = ctx.GetInt("TeamOwner")
	}
	if err != nil || owner == -1 {
		log.Println(err, owner, ctx.GetInt("TeamOwner"))
		httputil.NewError(ctx, http.StatusBadRequest, "Bad request")
		return
	}

	team := models.Team{
		OwnerID: uint(owner),
		Name:    t.Name,
		Country: t.Country,
		Budget:  t.Budget,
	}

	err = c.Repo.Create(&team)
	if err != nil {
		log.Println(err)
		httputil.NewError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	httputil.NoError(ctx, map[string]interface{}{
		"id": team.ID,
	})
}

// Handles a DELETE request to a team resource
// @Summary Delete a team and all of it's players
// @Description Delete a team and all of it's players
// @Tags Teams
// @Accept  json
// @Produce  json
// @Param id path int true "Team ID"
// @Success 200
// @Failure 400 {object} httputil.HTTPError
// @Failure 401 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /teams/{id} [delete]
// @Security BearerAuth
func (c *Controller) DeleteTeam(ctx *gin.Context) {
	team, err := c.getTeamFromRequest(ctx)
	if err != nil {
		return
	}

	err = c.Repo.DeleteTeam(&team)
	if err != nil {
		log.Println(err)
		httputil.NewError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	httputil.NoErrorEmpty(ctx)
}

// Handles a PATCH request to a team resource
// @Summary Update a team
// @Description Update a team
// @Tags Teams
// @Accept  json
// @Produce  json
// @Param id path int true "Team ID"
// @Success 200
// @Failure 400 {object} httputil.HTTPError
// @Failure 401 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /teams/{id} [patch]
// @Security BearerAuth
func (c *Controller) UpdateTeam(ctx *gin.Context) {
	user, err := c.getAuthenticatedUserFromRequest(ctx)
	team, err := c.getTeamFromRequest(ctx)
	if err != nil || (!user.IsAdmin() && !c.validateTeamOwner(ctx, user, team)) {
		return
	}

	t := c.fillDefaultTeamPayload(team)
	err = ctx.BindJSON(&t)
	if err != nil {
		log.Println(err)
		httputil.NewError(ctx, http.StatusBadRequest, "Bad request")
		return
	}

	team.Country = t.Country
	team.Name = t.Name
	if user.IsAdmin() {
		team.Budget = t.Budget
	}

	err = c.Repo.Update(&team)
	if err != nil {
		log.Println(err)
		httputil.NewError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	httputil.NoErrorEmpty(ctx)
}

func (c *Controller) RedirectToPlayers(ctx *gin.Context) {
	team, err := c.getTeamFromRequest(ctx)
	if err != nil {
		return
	}

	ctx.Set("PlayerOwner", team.OwnerID)
	ctx.Redirect(http.StatusTemporaryRedirect, "/api/players/" + ctx.Param("action"))
}

// Generate a json from a team model
func (c *Controller) getTeamPayload(team models.Team) models.ShowTeam {
	marketValue := 0
	players := c.Repo.GetPlayers(team.ID)

	playerModels := make([]models.ShowPlayer, 0)
	for _, p := range players {
		playerModels = append(playerModels, c.getPlayerPayload(p))
		marketValue += int(p.MarketValue)
	}
	return models.ShowTeam{
		ID:          team.ID,
		Name:        team.Name,
		Country:     team.Country,
		Budget:      team.Budget,
		Players:     playerModels,
		MarketValue: marketValue,
	}
}

// Get team from request
func (c *Controller) getTeamFromRequest(ctx *gin.Context) (models.Team, error) {
	id, err := c.parseIdFromRequest(ctx, "teamId")
	if err != nil {
		return models.Team{}, err
	}

	team, err := c.Repo.GetTeam(id)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, "Not found")
		return models.Team{}, err
	}
	return team, nil
}

// Validate a team owner
func (c *Controller) validateTeamOwner(ctx *gin.Context, user models.User, team models.Team) bool {
	if user.ID != team.OwnerID {
		httputil.NewError(ctx, http.StatusUnauthorized, "unauthorized")
		return false
	}
	return true
}

/// Fill the team payload with default values
func (c* Controller) fillDefaultTeamPayload(team models.Team) models.UpdateTeam {
	var payload models.UpdateTeam
	payload.Owner = int(team.OwnerID)
	payload.Budget = team.Budget
	payload.Name = team.Name
	payload.Country = team.Country
	return payload
}