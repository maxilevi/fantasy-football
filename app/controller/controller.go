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

// Return a new controller with a given repository
func NewController(repo repos.Repository) *Controller {
	return &Controller{Repo: repo}
}

// Get the user the request got authenticated with
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

// Parse an id into an int from a request parameter
func (c *Controller) parseIdFromRequest(ctx *gin.Context, paramName string) (uint, error) {
	id, ok := strconv.ParseInt(ctx.Param(paramName), 10, 32)
	if ok != nil {
		httputil.NewError(ctx, http.StatusBadRequest, "Invalid " + paramName)
		return 0, fmt.Errorf("invalid " + paramName)
	}
	return uint(id), nil
}
