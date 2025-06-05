package controllers

import (
	"Savannah_Screening_Test/clients"
	"Savannah_Screening_Test/dtos"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MessageRequest struct {
	Message string `json:"message" binding:"required"`
	To      string `json:"to" binding:"required"`
}

// SendTestMessage godoc
// @Summary Send a test SMS message
// @Description Sends an SMS message to a specified recipient using the SMS gateway
// @Tags messaging
// @Accept json
// @Produce json
// @Param request body MessageRequest true "Message payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /send-message [post]
func SendTestMessage(c *gin.Context) {
	var request MessageRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := clients.SendSMS(request.To, request.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "sending failed, contact backend admin"})
	}
	c.JSON(http.StatusOK, gin.H{"status": "message sent successfully"})
}

// SendTestEmail godoc
// @Summary Send a test email
// @Description Sends a test email to the system administrator
// @Tags messaging
// @Accept json
// @Produce json
// @Param request body dtos.EmailRequest true "Email payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /send-email [post]
func SendTestEmail(c *gin.Context) {
	var request dtos.EmailRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := clients.SendAdminEmail(request.Subject, request.Message)
	if err != nil {
		c.JSON(500, gin.H{"error": "sending email failed, contact backend admin"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "email sent successfully"})
}
