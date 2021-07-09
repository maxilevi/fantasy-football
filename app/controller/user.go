package controller

import (
	"../httputil"
	"../models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/mail"
)

// Handles GET request to the user resource when no ID is provided
// @Summary Get the logged in user
// @Description Get user by ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} models.ShowUser
// @Failure 401 {object} httputil.HTTPError
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /user/me [get]
// @Security BearerAuth[write, admin]
func (c *Controller) ShowMyself(ctx *gin.Context) {
	user, err := c.getAuthenticatedUserFromRequest(ctx)
	if err != nil {
		return
	}

	payload, err := c.getShowUserPayload(ctx, user)
	if err != nil {
		return
	}
	httputil.NoError(ctx, payload)
}

// Handles GET request to the user resource
// @Summary Get a user
// @Description Get user by ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} models.ShowUser
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /user/{id} [get]
func (c *Controller) ShowUser(ctx *gin.Context) {
	user, err := c.getUserFromRequest(ctx)
	if err != nil {
		return
	}

	payload, err := c.getShowUserPayload(ctx, user)
	if err != nil {
		log.Println(err)
		return
	}
	httputil.NoError(ctx, payload)
}

// Handles DELETE requests to the user's resource
// @Summary Delete a user
// @Description Delete user by ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200
// @Failure 401 {object} httputil.HTTPError
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /user/{id} [delete]
// @Security BearerAuth[write, admin]
func (c *Controller) DeleteUser(ctx *gin.Context) {
	user, err := c.getUserFromRequest(ctx)
	if err != nil {
		return
	}

	err = c.Repo.Delete(user)
	if err != nil {
		log.Println(err)
		httputil.NewError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}
	httputil.NoErrorEmpty(ctx)
}

// Handles PATCH requests to the user's resource
// @Summary Update a user
// @Description Update user by ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param email body models.UpdateUser true "UpdateUser data"
// @Success 200
// @Failure 401 {object} httputil.HTTPError
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /user/{id} [patch]
// @Security BearerAuth[admin]
func (c *Controller) UpdateUser(ctx *gin.Context) {
	user, err := c.getUserFromRequest(ctx)
	if err != nil {
		return
	}

	var t models.UpdateUser
	err = ctx.BindJSON(&t)
	if err != nil {
		log.Println(err)
		httputil.NewError(ctx, http.StatusBadRequest, "Incorrect body parameters")
		return
	}

	user.Email = t.Email
	err = c.Repo.Update(user)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, "Incorrect body parameters")
		return
	}

	httputil.NoErrorEmpty(ctx)
}

// Handles POST request to the user resource
// @Summary Create a user
// @Description Register a new user
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body models.CreateUser true "CreateUser"
// @Success 200
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /user [post]
// @Security BearerAuth
func (c *Controller) CreateUser(ctx *gin.Context) {
	var t models.CreateUser
	err := ctx.BindJSON(&t)
	log.Println(err, t)
	if err != nil {
		log.Println(err)
		httputil.NewError(ctx, http.StatusBadRequest, "Incorrect body parameters")
		return
	}
	if !c.validEmail(t.Email) {
		httputil.NewError(ctx, http.StatusBadRequest, "Invalid email")
		return
	}
	if c.emailExists(t.Email) {
		httputil.NewError(ctx, http.StatusBadRequest, "Provided email is already registered")
		return
	}
	if !c.validPassword(t.Password) {
		httputil.NewError(ctx, http.StatusBadRequest, "Password needs a minimum of at least 8 characters")
		return
	}

	log.Println("Registering a new user...")

	user, err := c.registerUser(t.Email, t.Password)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	httputil.NoError(ctx, map[string]interface{}{
		"id": user.ID,
	})
}

// Returns a bool to check if the password is valid
func (c *Controller) validPassword(password string) bool {
	return len(password) >= 8
}

// Returns a bool to check if the email is valid
func (c *Controller) validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Return a bool to see if the email is already registered
func (c *Controller) emailExists(email string) bool {
	_, err := c.Repo.GetUserByEmail(email)
	return err == nil
}

// Registers a new user with the given credentials
func (c *Controller) registerUser(email, password string) (models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	return c.Repo.CreateUser(email, hashedPassword, 0)
}

// Parse a user from the request parameters or return an error if not found
func (c *Controller) getUserFromRequest(ctx *gin.Context) (models.User, error) {
	id, err := c.parseIdFromRequest(ctx)
	if err != nil {
		return models.User{}, err
	}

	user, err := c.Repo.GetUserById(id)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, "User not found")
		return models.User{}, err
	}
	return user, err
}

func (c *Controller) getShowUserPayload(ctx *gin.Context, user models.User) (models.ShowUser, error) {
	team, err := c.Repo.GetUserTeam(user)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, "Internal server error")
		return models.ShowUser{}, err
	}

	return models.ShowUser{
		ID: user.ID,
		Email: user.Email,
		Team: c.getTeamPayload(team),
	}, nil
}
