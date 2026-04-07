package shipping

import (
	"net/http"
	"phoenix/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAll(c *gin.Context) {
	var shippings []models.Shipping
	if err := h.DB.Preload("FromUser").Preload("ToUser").Find(&shippings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": shippings})
}