package controller

import (
	"../httputil"
	"../middleware"
	"../models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// Handles creating a new session
// @Summary Create a new session
// @Description Creates a new session for a given set of credentials, returns a JWT token to be used as Bearer token.
// @Tags Session
// @Accept  json
// @Produce  json
// @Param user body models.CreateSession true "Credentials"
// @Success 200 {object} models.SessionToken
// @Failure 401 {object} httputil.HTTPError
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /session [post]
func (c *Controller) CreateSession(ctx *gin.Context) {
	var t models.CreateSession
	err := ctx.BindJSON(&t)

	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, "Incorrect body parameters")
		return
	}

	user, err := c.loginAndGetUser(t.Email, t.Password)
	if err != nil {
		httputil.NewError(ctx, http.StatusUnauthorized, "Invalid email or password")
		return
	}
	token, err := c.createToken(user)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}
	httputil.NoError(ctx, map[string]interface{}{
		"token": token,
	})
}

// Tries to login the user and returns it if successful
func (c *Controller) loginAndGetUser(email, password string) (*models.User, error) {
	user, err := c.Repo.GetUserByEmail(email)
	if err != nil {
		return &user, err
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
	if err != nil {
		return &user, err
	}

	return &user, nil
}

// Creates a JWT token for a specific user
func (c *Controller) createToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &middleware.AuthClaims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(middleware.JWTKey)
}
