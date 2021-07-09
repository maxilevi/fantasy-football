package controller

import (
	"../httputil"
	"../models"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
)

func (c *Controller) ShowAllTransfer(ctx *gin.Context) {
	transfers := c.Repo.GetTransfers()

	arr := make([]models.ShowTransfer, 0)
	for _, transfer := range transfers {
		arr = append(arr, models.ShowTransfer{
			ID:     transfer.ID,
			Player: c.getPlayerPayload(transfer.Player),
			Ask:    transfer.Ask,
		})
	}

	httputil.NoError(ctx, map[string]interface{}{
		"transfers": transfers,
	})
}

func (c *Controller) ShowTransfer(ctx *gin.Context) {
	transfer, err := c.getTransferFromRequest(ctx)
	if err != nil {
		return
	}

	httputil.NoError(ctx, transfer)
}

// Handle creating a transfer
func (c *Controller) CreateTransfer(w http.ResponseWriter, req *http.Request) {

}

// Handle editing a transfer
func (c *Controller) UpdateTransfer(ctx *gin.Context) {

}

// Handle deleting a transfer
func (c *Controller) DeleteTransfer(ctx *gin.Context) {
	transfer, err1 := c.getTransferFromRequest(ctx)
	user, err2 := c.getAuthenticatedUserFromRequest(ctx)
	if err1 != nil || err2 != nil {
		return
	}

	if !user.IsAdmin() && user.ID != transfer.Player.Team.OwnerID {
		httputil.NewError(ctx, http.StatusUnauthorized, "Trying to delete a not owned transfer")
		return
	}

	err := c.Repo.Delete(transfer)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	httputil.NoErrorEmpty(ctx)
}

// Execute a transfer and update the records if successful
func (c *Controller) executeTransfer(transfer models.Transfer, seller, buyer models.Team) error {
	// Randomly update the player value
	player := transfer.Player
	player.MarketValue = int32(float64(player.MarketValue) * (1.1 + rand.Float64()*0.9))

	// Actually do the transfer
	player.TeamID = buyer.ID

	seller.Budget += transfer.Ask
	buyer.Budget -= transfer.Ask

	return c.Repo.RunInTransaction(func() error {
		err1 := c.Repo.Update(&player)
		err2 := c.Repo.Update(&buyer)
		err3 := c.Repo.Update(&seller)
		if err1 != nil || err2 != nil || err3 != nil {
			return fmt.Errorf("failed to save models")
		}
		return nil
	})
}

// Gets transfers from the request
func (c *Controller) getTransferFromRequest(ctx *gin.Context) (models.Transfer, error) {
	id, err := c.parseIdFromRequest(ctx)
	transfer, err := c.Repo.GetTransfer(id)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, "Not found")
		return models.Transfer{}, err
	}
	return transfer, nil
}
