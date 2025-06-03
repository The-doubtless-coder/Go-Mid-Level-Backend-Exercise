package controllers

import (
	"Savannah_Screening_Test/clients"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MessageRequest struct {
	Message string `json:"message" binding:"required"`
	To      string `json:"to" binding:"required"`
}

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

type EmailRequest struct {
	Subject string `json:"subject" binding:"required"`
	Message string `json:"message" binding:"required"`
}

func SendTestEmail(c *gin.Context) {
	var request EmailRequest
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
