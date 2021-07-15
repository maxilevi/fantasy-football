package controller

import (
	"../httputil"
	"../models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/mail"
	"strconv"
)

// Handles GET request to the user resource when no ID is provided
// @Summary Get the logged in user
// @Description Get user by ID
// @Tags Me
// @Accept  json
// @Produce  json
// @Success 200 {object} models.ShowUser
// @Failure 401 {object} httputil.HTTPError
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /me [get]
// @Security BearerAuth[write, admin]
func (c *Controller) RedirectMyself(ctx *gin.Context) {
	user, err := c.getAuthenticatedUserFromRequest(ctx)
	if err != nil {
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, "/api/users/"+strconv.Itoa(int(user.ID)))

}

// @Summary Get the logged in user's team
// @Description Get the logged in user's team
// @Tags Me
// @Accept  json
// @Produce  json
// @Success 200 {object} models.ShowUser
// @Failure 401 {object} httputil.HTTPError
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /me/team [get]
// @Security BearerAuth
func (c *Controller) GetMyTeam(ctx *gin.Context) {
	c.RedirectMyTeam(ctx, "")
}

// @Summary Edit the logged in user's team
// @Description Edit the logged in user's team
// @Tags Me
// @Accept  json
// @Produce  json
// @Success 200
// @Param team body models.UpdateTeam true "Update team payload"
// @Failure 400 {object} httputil.HTTPError
// @Failure 401 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /me/team [patch]
// @Security BearerAuth
func (c *Controller) EditMyTeam(ctx *gin.Context) {
	c.RedirectMyTeam(ctx, "")
}

// @Summary Get the logged in user's team players
// @Description Get the logged in user's team players
// @Tags Me
// @Accept  json
// @Produce  json
// @Success 200 {array} models.ShowPlayer
// @Failure 401 {object} httputil.HTTPError
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /me/team/players [get]
// @Security BearerAuth
func (c *Controller) GetMyPlayers(ctx *gin.Context) {
	c.RedirectMyTeam(ctx, "/players")
}

// @Summary Edit the logged in user's team players
// @Description Get the logged in user's team players
// @Tags Me
// @Accept  json
// @Produce  json
// @Success 200
// @Param player body models.UpdatePlayer true "Update player"
// @Failure 401 {object} httputil.HTTPError
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /team/players/{id} [patch]
// @Security BearerAuth
func (c *Controller) EditMyPlayer(ctx *gin.Context) {
	c.RedirectMyPlayers(ctx)
}

// @Summary Get the logged in user's team player
// @Description Get the logged in user's team player
// @Tags Me
// @Accept  json
// @Produce  json
// @Success 200 {object} models.ShowPlayer
// @Failure 401 {object} httputil.HTTPError
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /me/team/players/{id} [get]
// @Security BearerAuth
func (c *Controller) GetMyPlayer(ctx *gin.Context) {
	c.RedirectMyPlayers(ctx)
}

// Redirect to the players resource
func (c *Controller) RedirectMyPlayers(ctx *gin.Context) {
	id, err := c.parseIdFromRequest(ctx, "playerId")
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, "A bad player id provided")
		return
	}
	c.RedirectMyTeam(ctx, "/players/"+strconv.Itoa(int(id)))
}

// Redirect to the team resource
func (c *Controller) RedirectMyTeam(ctx *gin.Context, postfix string) {
	user, err := c.getAuthenticatedUserFromRequest(ctx)
	if err != nil {
		return
	}

	team, err := c.Repo.GetUserTeam(user)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, "Not found")
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, "/api/teams/"+strconv.Itoa(int(team.ID))+postfix)
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
// @Router /users/{id} [get]
func (c *Controller) ShowUser(ctx *gin.Context) {
	authUser, err1 := c.getAuthenticatedUserFromRequest(ctx)
	user, err2 := c.getUserFromRequest(ctx)
	if err1 != nil || err2 != nil {
		return
	}
	if !authUser.IsAdmin() && authUser.ID != user.ID {
		httputil.NewError(ctx, http.StatusUnauthorized, "Can't query another user's information")
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
// @Summary Delete a user and all of it's associated resources
// @Description Delete a user by ID and all of it's associated resources
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200
// @Failure 401 {object} httputil.HTTPError
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /users/{id} [delete]
// @Security BearerAuth[write, admin]
func (c *Controller) DeleteUser(ctx *gin.Context) {
	user, err := c.getUserFromRequest(ctx)
	if err != nil {
		return
	}
	err = c.Repo.RunInTransaction(func() error {
		team, err := c.Repo.GetUserTeam(user)
		if err != nil {
			return err
		}

		err = c.Repo.DeleteTeam(&team)
		if err != nil {
			return err
		}

		return c.Repo.Delete(&user)
	})
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
// @Router /users/{id} [patch]
// @Security BearerAuth[admin]
func (c *Controller) UpdateUser(ctx *gin.Context) {
	user, err := c.getUserFromRequest(ctx)
	if err != nil {
		return
	}

	t := c.fillDefaultUserPayload(user)
	err = ctx.ShouldBindJSON(&t)
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
// @Router /users [post]
// @Security BearerAuth
func (c *Controller) CreateUser(ctx *gin.Context) {
	var t models.CreateUser
	err := ctx.ShouldBindJSON(&t)
	if err != nil {
		log.Println(err)
		httputil.NewError(ctx, http.StatusBadRequest, "Incorrect body parameters")
		return
	}
	if !c.validEmail(t.Email) {
		log.Println(err)
		httputil.NewError(ctx, http.StatusBadRequest, "Invalid email")
		return
	}
	if c.emailExists(t.Email) {
		log.Println(err)
		httputil.NewError(ctx, http.StatusBadRequest, "Provided email is already registered")
		return
	}
	if !c.validPassword(t.Password) {
		log.Println(err)
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

func (c *Controller) RedirectToTeam(ctx *gin.Context) {
	user, err := c.getUserFromRequest(ctx)
	if err != nil {
		return
	}

	team, err := c.Repo.GetUserTeam(user)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, "Team not found")
		return
	}

	ctx.Set("TeamOwner", user.ID)
	action := ctx.Param("action")
	if ctx.Request.Method == "POST" && (action == "" || action == "/") {
		ctx.Redirect(http.StatusTemporaryRedirect, "/api/teams/")
	} else {
		ctx.Redirect(http.StatusTemporaryRedirect, "/api/teams/"+strconv.Itoa(int(team.ID))+action)
	}
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
	id, err := c.parseIdFromRequest(ctx, "userId")
	if err != nil {
		return models.User{}, err
	}

	user, err := c.Repo.GetUserById(id)
	if err != nil {
		log.Println(err)
		httputil.NewError(ctx, http.StatusNotFound, "User not found")
		return models.User{}, err
	}
	return user, err
}

// Get the payload for showing an user
func (c *Controller) getShowUserPayload(ctx *gin.Context, user models.User) (models.ShowUser, error) {
	team, err := c.Repo.GetUserTeam(user)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, "Internal server error")
		return models.ShowUser{}, err
	}

	return models.ShowUser{
		ID:    user.ID,
		Email: user.Email,
		Team:  c.getTeamPayload(team, c.Repo.GetPlayers(team.ID)),
	}, nil
}

// Fill the user payload with default values
func (c *Controller) fillDefaultUserPayload(user models.User) models.UpdateUser {
	var payload models.UpdateUser
	payload.Email = user.Email
	return payload
}
