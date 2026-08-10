package controller

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kirtanwyn/allPanelexch/config"
	"github.com/kirtanwyn/allPanelexch/dto"
	"github.com/kirtanwyn/allPanelexch/repository"
	"github.com/kirtanwyn/allPanelexch/service"
)

func GetAccountStatement(c *gin.Context) {
	session := sessions.Default(c)
	userIDRaw := session.Get("CLIENT_LOGIN_ID")
	if userIDRaw == nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Unauthorized or Session Expired"})
		return
	}
	userID := userIDRaw.(int)

	var req dto.AccountStatementRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Invalid request parameters"})
		return
	}

	repo := repository.NewAccountStatementRepository(config.DB)
	svc := service.NewAccountStatementService(repo)

	data, err := svc.GetAccountStatement(userID, req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   data,
	})
}

func GetAccountBetStatement(c *gin.Context) {
	session := sessions.Default(c)
	userIDRaw := session.Get("CLIENT_LOGIN_ID")
	if userIDRaw == nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Unauthorized or Session Expired"})
		return
	}
	userID := userIDRaw.(int)

	var req dto.AccountBetStatementRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Invalid request parameters"})
		return
	}

	repo := repository.NewAccountStatementRepository(config.DB)
	svc := service.NewAccountStatementService(repo)

	data, err := svc.GetAccountBetStatement(userID, req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   data,
	})
}
