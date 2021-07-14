package controller

import (
	"../httputil"
	"../models"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

// Handles GET requests to the transfers resource
// @Summary Show all transfers
// @Description Show all transfers and filter by country, team name, player name, age and value
// @Tags Transfers
// @Accept  json
// @Produce  json
// @Param country query string false "Filter by the player's country"
// @Param team_name query string false "Filter by the player's team name"
// @Param player_name query string false "Filter by the player's complete name"
// @Param min_age query string false "Filter by the player's age"
// @Param max_age query string false "Filter by the player's age"
// @Param min_value query string false "Filter by the transfer ask value"
// @Param max_value query string false "Filter by the transfer ask value"
// @Success 200 {array} models.ShowTransfer
// @Router /transfers [get]
func (c *Controller) ListTransfers(ctx *gin.Context) {
	filter := c.parseTransferFilters(ctx)
	transfers := c.Repo.GetTransfers()

	arr := make([]models.ShowTransfer, 0)
	for _, transfer := range transfers {
		if !filter.Matches(transfer) {
			continue
		}
		arr = append(arr, models.ShowTransfer{
			ID:     transfer.ID,
			Player: c.getPlayerPayload(transfer.Player),
			Ask:    transfer.Ask,
		})
	}

	httputil.NoError(ctx, map[string]interface{}{
		"transfers": arr,
	})
}

// Handles GET requests to the transfers resource
// @Summary Show a transfer
// @Description Get a transfer by ID
// @Tags Transfers
// @Accept  json
// @Produce  json
// @Param id path int true "Transfer ID"
// @Success 200 {object} models.ShowTransfer
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /transfers/{id} [get]
func (c *Controller) ShowTransfer(ctx *gin.Context) {
	transfer, err := c.getTransferFromRequest(ctx)
	if err != nil {
		return
	}

	payload := models.ShowTransfer{
		Player: c.getPlayerPayload(transfer.Player),
		Ask:    transfer.Ask,
	}

	httputil.NoError(ctx, payload)
}

// Handles a POST request to a transfer resource
// @Summary Create a new transfer
// @Description Create a new transfer
// @Tags Transfers
// @Accept  json
// @Produce  json
// @Success 200
// @Param transfer body models.CreateTransfer true "Create transfer"
// @Failure 400 {object} httputil.HTTPError
// @Failure 401 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /transfers [post]
// @Security BearerAuth
func (c *Controller) CreateTransfer(ctx *gin.Context) {
	user, err := c.getAuthenticatedUserFromRequest(ctx)
	if err != nil {
		return
	}

	var t models.CreateTransfer
	err = ctx.BindJSON(&t)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, "Invalid body parameters")
		return
	}

	player, err := c.Repo.GetPlayer(t.PlayerID)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, "Player not found")
		return
	}

	if !user.IsAdmin() && player.Team.UserID != user.ID {
		httputil.NewError(ctx, http.StatusUnauthorized, "Trying to create a transfer on a player not owned")
		return
	}
	transfer, err := c.Repo.GetTransferWithPlayer(&player)
	if err == nil {
		httputil.NewError(ctx, http.StatusBadRequest, "Player already has an open transfer")
		return
	}

	transfer = models.Transfer{
		PlayerID: t.PlayerID,
		Ask:      t.Ask,
	}
	err = c.Repo.Create(&transfer)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	httputil.NoError(ctx, map[string]interface{}{
		"id": transfer.ID,
	})
}

// Handles PATCH requests to the transfers resource
// @Summary Updates a existing transfer.
// @Description Updates a existing transfer by ID
// @Tags Transfers
// @Accept  json
// @Produce  json
// @Param id path int true "Transfer ID"
// @Param transfer body models.UpdateTransfer true "Update transfer"
// @Success 200
// @Failure 401 {object} httputil.HTTPError
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /transfers/{id} [patch]
// @Security BearerAuth[write, admin]
func (c *Controller) UpdateTransfer(ctx *gin.Context) {
	transfer, err1 := c.getTransferFromRequest(ctx)
	user, err2 := c.getAuthenticatedUserFromRequest(ctx)
	if err1 != nil || err2 != nil {
		return
	}

	t := c.fillDefaultTransferPayload(transfer)
	err := ctx.BindJSON(&t)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, "Invalid body parameters")
		return
	}

	if !user.IsAdmin() && user.ID != transfer.Player.Team.UserID {
		httputil.NewError(ctx, http.StatusUnauthorized, "Trying to update a not owned transfer")
		return
	}

	transfer.Ask = t.Ask

	err = c.Repo.Update(&transfer)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	httputil.NoErrorEmpty(ctx)
}

// @Summary Buy a transfer
// @Description Buys a transfer with a specific id and buys it. Updates budgets, values and finally deletes the transfer.
// @Tags Transfers
// @Accept  json
// @Produce  json
// @Param id path int true "Transfer ID"
// @Success 200
// @Failure 401 {object} httputil.HTTPError
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /transfers/{id}/buy [put]
// @Security BearerAuth
func (c *Controller) BuyTransfer(ctx *gin.Context) {
	transfer, err1 := c.getTransferFromRequest(ctx)
	user, err2 := c.getAuthenticatedUserFromRequest(ctx)
	if err1 != nil || err2 != nil {
		return
	}

	if  transfer.Player.Team.UserID == user.ID {
		log.Println("Trying to buy own player")
		httputil.NewError(ctx, http.StatusBadRequest, "Cannot buy your own player")
		return
	}

	buyer, err := c.Repo.GetUserTeam(user)
	if err != nil {
		log.Println(err)
		httputil.NewError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	if buyer.Budget < transfer.Ask {
		httputil.NewError(ctx, http.StatusBadRequest, fmt.Sprintf("Team does not have enough money to execute the purchase (%v < %v)", buyer.Budget, transfer.Ask))
		return
	}

	err = c.doExecuteTransfer(&transfer, buyer)
	if err != nil {
		log.Println(err)
		httputil.NewError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}
	p, _ := c.Repo.GetPlayer(transfer.PlayerID)
	fmt.Println(p.TeamID, buyer.ID)
	fmt.Println(len(c.Repo.GetPlayers(buyer.ID)))


	httputil.NoErrorEmpty(ctx)
}

// Handles DELETE requests to the transfers resource
// @Summary Delete a transfer
// @Description Delete a transfer by ID
// @Tags Transfers
// @Accept  json
// @Produce  json
// @Param id path int true "Transfer ID"
// @Success 200
// @Failure 401 {object} httputil.HTTPError
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /transfers/{id} [delete]
// @Security BearerAuth[write, admin]
func (c *Controller) DeleteTransfer(ctx *gin.Context) {
	transfer, err1 := c.getTransferFromRequest(ctx)
	user, err2 := c.getAuthenticatedUserFromRequest(ctx)
	if err1 != nil || err2 != nil {
		return
	}

	if !user.IsAdmin() && user.ID != transfer.Player.Team.UserID {
		httputil.NewError(ctx, http.StatusUnauthorized, "Trying to delete a not owned transfer")
		return
	}

	err := c.Repo.Delete(&transfer)
	if err != nil {
		log.Println(err)
		httputil.NewError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	httputil.NoErrorEmpty(ctx)
}

// Execute a transfer and update the records if successful
func (c *Controller) doExecuteTransfer(transfer *models.Transfer, buyer models.Team) error {
	seller := transfer.Player.Team
	// Randomly update the player value
	player := transfer.Player
	player.MarketValue = int32(float64(player.MarketValue) * (1.1 + rand.Float64()*0.9))

	// Actually do the transfer
	player.TeamID = buyer.ID
	player.Team = buyer

	seller.Budget += transfer.Ask
	buyer.Budget -= transfer.Ask

	return c.Repo.RunInTransaction(func() error {
		err1 := c.Repo.Update(&player)
		err2 := c.Repo.Update(&buyer)
		err3 := c.Repo.Update(&seller)
		err4 := c.Repo.Delete(&transfer)
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			return fmt.Errorf("failed to save models")
		}
		return nil
	})
}

// Gets transfers from the request
func (c *Controller) getTransferFromRequest(ctx *gin.Context) (models.Transfer, error) {
	id, err := c.parseIdFromRequest(ctx, "transferId")
	transfer, err := c.Repo.GetTransfer(id)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, "Not found")
		return models.Transfer{}, err
	}
	return transfer, nil
}

// Parse URL parameters into a transferFilters object
func (c *Controller) parseTransferFilters(ctx *gin.Context) transferFilters {
	filter := transferFilters{
		Country:    ctx.Param("country"),
		TeamName:   ctx.Param("team_name"),
		PlayerName: ctx.Param("player_name"),
	}

	ageFilter := ctx.Param("min_age")
	if age, err := strconv.ParseInt(ageFilter, 10, 32); err != nil {
		filter.MinAgeFilter = -1
	} else {
		filter.MinAgeFilter = int(age)
	}

	valueFilter := ctx.Param("min_value")
	if value, err := strconv.ParseInt(valueFilter, 10, 32); err != nil {
		filter.MinValueFilter = -1
	} else {
		filter.MinValueFilter = int(value)
	}

	ageFilter = ctx.Param("max_age")
	if age, err := strconv.ParseInt(ageFilter, 10, 32); err != nil {
		filter.MaxAgeFilter = math.MaxInt32
	} else {
		filter.MaxAgeFilter = int(age)
	}

	valueFilter = ctx.Param("max_value")
	if value, err := strconv.ParseInt(valueFilter, 10, 32); err != nil {
		filter.MaxValueFilter = math.MaxInt32
	} else {
		filter.MaxValueFilter = int(value)
	}

	return filter
}

// A group of filters to apply to matches
type transferFilters struct {
	Country     string
	TeamName    string
	PlayerName  string
	MinAgeFilter   int
	MinValueFilter int
	MaxAgeFilter   int
	MaxValueFilter int
}

// Returns a bool that tells if the transfer matches with the filter
func (f *transferFilters) Matches(transfer models.Transfer) bool {
	return strings.Contains(strings.ToLower(transfer.Player.FirstName+" "+transfer.Player.LastName), strings.ToLower(f.PlayerName)) &&
		strings.Contains(strings.ToLower(transfer.Player.Team.Name), strings.ToLower(f.TeamName)) &&
		strings.Contains(strings.ToLower(transfer.Player.Country), strings.ToLower(f.Country)) &&
		transfer.Ask >= f.MinValueFilter && transfer.Player.Age >= f.MinAgeFilter &&
		transfer.Ask <= f.MaxValueFilter && transfer.Player.Age <= f.MaxAgeFilter
}

/// Fill the transfer payload with default values
func (c* Controller) fillDefaultTransferPayload(transfer models.Transfer) models.UpdateTransfer {
	var payload models.UpdateTransfer
	payload.Ask = transfer.Ask
	return payload
}
