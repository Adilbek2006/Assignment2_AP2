package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"order-service/internal/domain"
	"order-service/internal/usecase"
)

type OrderHandler struct {
	UC *usecase.OrderUseCase
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	order := &domain.Order{
		CustomerID: req.CustomerID,
		ItemName:   req.ItemName,
		Amount:     req.Amount,
	}

	if err := h.UC.CreateOrder(order); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) Get(c *gin.Context) {
	id := c.Param("id")
	order, err := h.UC.GetOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	id := c.Param("id")
	if err := h.UC.CancelOrder(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cancelled"})
}

func (h *OrderHandler) GetStats(c *gin.Context) {
	stats, err := h.UC.GetOrderStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get stats"})
		return
	}

	if stats == nil {
		stats = &domain.OrderStats{}
	}

	resp := OrderStatsResponse{
		Total:     stats.Total,
		Pending:   stats.Pending,
		Paid:      stats.Paid,
		Failed:    stats.Failed,
		Cancelled: stats.Cancelled,
	}

	c.JSON(http.StatusOK, resp)
}
