// pkg/server/server.go
package server

import (
	"ReceiptApi/internal/controller"
	"ReceiptApi/internal/repository"
	"ReceiptApi/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, receiptController *controller.ReceiptController) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/receipts", receiptController.SaveReciept)
	r.GET("/receipts/:id/points", receiptController.GetPoints)
}

func SetupRouter() *gin.Engine {
	db := repository.NewInMemoryDB()
	receiptService := service.NewReceiptService(db)
	receiptController := controller.NewReceiptController(receiptService)
	r := gin.Default()
	SetupRoutes(r, receiptController)
	return r
}
