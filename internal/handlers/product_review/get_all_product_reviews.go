package product_review

import (
	"net/http"
	"phoenix/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAll(c *gin.Context) {
	var reviews []models.ProductReview
	if err := h.DB.Preload("User").Preload("Product").Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": reviews})
}