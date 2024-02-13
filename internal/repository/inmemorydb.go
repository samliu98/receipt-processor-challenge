package repository

import "ReceiptApi/models"

type InMemoryDB struct {
	reciepts map[string]models.Receipt
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		reciepts: make(map[string]models.Receipt),
	}
}

func (db *InMemoryDB) SaveReciept(id string, receipt models.Receipt) error {
	db.reciepts[id] = receipt
	return nil
}

func (db *InMemoryDB) GetPoints(id string) int {
	reciept, exist := db.reciepts[id]
	if !exist {
		return -1
	}
	return reciept.Points
}
