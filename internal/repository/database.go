package repository

import "ReceiptApi/models"

type Database interface {
	SaveReciept(id string, receipt models.Receipt) error
	GetPoints(id string) int
}
