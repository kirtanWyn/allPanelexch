package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kirtanwyn/allPanelexch/config"
	"github.com/kirtanwyn/allPanelexch/dto"
	"github.com/kirtanwyn/allPanelexch/repository"
	"github.com/kirtanwyn/allPanelexch/service"
)

func GetCasinoResult(c *gin.Context) {
	var req dto.CasinoResultRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repo := repository.NewCasinoResultRepository(config.DB)
	svc := service.NewCasinoResultService(repo)

	res, err := svc.GetCasinoResults(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch casino results"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetResultCards(c *gin.Context) {
	var req dto.GetResultCardsRequest
	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	repo := repository.NewCasinoResultRepository(config.DB)
	svc := service.NewCasinoResultService(repo)

	htmlContent, err := svc.GetResultCardsHTML(req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error retrieving cards")
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
}
