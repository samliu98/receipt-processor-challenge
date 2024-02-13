package repository

import "ReceiptApi/models"

type InMemoryDB struct {
	reciepts1 map[string]models.Receipt
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		reciepts1: make(map[string]models.Receipt),
	}
}

func (db *InMemoryDB) SaveReciept(id string, receipt models.Receipt) error {
	db.reciepts1[id] = receipt
	return nil
}

func (db *InMemoryDB) GetPoints(id string) int {
	reciept, exist := db.reciepts1[id]
	if !exist {
		return -1
	}
	return reciept.Points
}
