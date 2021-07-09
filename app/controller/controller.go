package controller

import (
	"../httputil"
	"../models"
	"../repos"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// Controller example
type Controller struct {
	Repo repos.Repository
}

func NewController(repo repos.Repository) *Controller {
	return &Controller{Repo: repo}
}

func (c *Controller) getAuthenticatedUserFromRequest(ctx *gin.Context) (models.User, error) {
	val, ok := ctx.Get("user")
	if !ok {
		httputil.NewError(ctx, http.StatusUnauthorized, "Authenticated user not found")
		return models.User{}, fmt.Errorf("internal server error")
	}
	user, ok := val.(models.User)
	if !ok {
		httputil.NewError(ctx, http.StatusBadRequest, "Authenticated user not found")
		return models.User{}, fmt.Errorf("internal server error")
	}
	return user, nil
}

func (c *Controller) parseIdFromRequest(ctx *gin.Context) (uint, error) {
	id, ok := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if ok != nil {
		httputil.NewError(ctx, http.StatusBadRequest, "Invalid id")
		return 0, fmt.Errorf("invalid id")
	}
	return uint(id), nil
}
