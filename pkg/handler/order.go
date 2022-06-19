package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getOrderById(c *gin.Context) {

	orderId, _ := c.GetQuery("id")
	orderData := h.services.Cache[orderId]

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": orderData,
	})
}
