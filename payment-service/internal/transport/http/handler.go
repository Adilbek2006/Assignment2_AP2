package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"payment-service/internal/usecase"
)

type PaymentHandler struct {
	UC *usecase.PaymentUseCase
}

func (h *PaymentHandler) HandleAuthorize(c *gin.Context) {
	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	p, err := h.UC.Process(req.OrderID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         p.Status,
		"transaction_id": p.TransactionID,
	})
}

func (h *PaymentHandler) HandleGetStatus(c *gin.Context) {
	orderID := c.Param("order_id")
	p, err := h.UC.GetStatus(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}
