package product

import (
	"net/http"
	"phoenix/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAll(c *gin.Context) {
	var products []models.Product
	if err := h.DB.Preload("Images").Preload("Tags").Preload("Reviews").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}