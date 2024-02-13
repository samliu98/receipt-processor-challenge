package controller

import (
	"ReceiptApi/internal/service"
	"ReceiptApi/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReceiptController struct {
	service *service.ReceiptService
}

func NewReceiptController(service *service.ReceiptService) *ReceiptController {
	return &ReceiptController{service: service}
}

func (rc *ReceiptController) SaveReciept(c *gin.Context) {
	var receipt models.Receipt
	if err := c.ShouldBindJSON(&receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := rc.service.ValidateReceipt(receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	receipt.ID = uuid.New().String()
	receipt.Points = rc.service.CalculatePoints(receipt)
	rc.service.SaveReciept(receipt.ID, receipt)

	c.JSON(http.StatusOK, gin.H{"id": receipt.ID})
}

func (rc *ReceiptController) GetPoints(c *gin.Context) {
	id := c.Param("id")
	point := rc.service.GetPoints(id)

	if point == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"points": point})
}
