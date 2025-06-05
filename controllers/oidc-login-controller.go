package controllers

import (
	"Savannah_Screening_Test/clients"
	"Savannah_Screening_Test/dtos"
	"github.com/gin-gonic/gin"
	"net/http"
)

// LoginHandler godoc
// @Summary Customer login
// @Description Authenticates a user via password grant and returns an access and refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.LoginRequest true "Login credentials"
// @Success 200 {object} dtos.LoginResponse "Login success"
// @Failure 400 {object} dtos.ErrorResponse "Invalid input"
// @Failure 401 {object} dtos.ErrorResponse "Invalid credentials"
// @Router /users/signin [post]
func LoginHandler(c *gin.Context) {
	var request dtos.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//var token dtos.TokenResponse
	token, err := clients.LoginWithPasswordGrant(request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	response := dtos.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.ExpiresIn,
	}

	c.JSON(http.StatusOK, response)
}
