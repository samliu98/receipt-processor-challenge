package service

import (
	"ReceiptApi/internal/repository"
	"ReceiptApi/models"
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ReceiptService struct {
	db repository.Database
}

func NewReceiptService(db repository.Database) *ReceiptService {
	return &ReceiptService{db: db}
}

type InMemoryDB struct {
	reciepts1 map[string]models.Receipt
}

func (rs *ReceiptService) SaveReciept(id string, receipt models.Receipt) error {
	rs.db.SaveReciept(id, receipt)
	return nil
}

func (rs *ReceiptService) GetPoints(id string) int {
	points := rs.db.GetPoints(id)

	return points
}

func (rs *ReceiptService) CalculatePoints(receipt models.Receipt) int {
	points := 0

	// One point for every alphanumeric character in the retailer name.
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	processedString := reg.ReplaceAllString(receipt.Retailer, "")
	points += len(processedString)

	// 50 points if the total is a round dollar amount with no cents.
	total, _ := strconv.ParseFloat(receipt.Total, 64)
	if total-float64(int(total)) == 0.0 {
		points += 50
	}

	// 25 points if the total is a multiple of 0.25.
	if int(total*100)%25 == 0 {
		points += 25
	}

	// 5 points for every two items on the receipt.
	points += (len(receipt.Items) / 2) * 5

	// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer.
	for _, item := range receipt.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(price * 0.2))
		}
	}

	// 6 points if the day in the purchase date is odd.
	date, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if date.Day()%2 != 0 {
		points += 6
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	t, _ := time.Parse("15:04", receipt.PurchaseTime)
	if (t.Hour() == 14 && t.Minute() > 0) || (t.Hour() >= 15 && t.Hour() < 16) {
		points += 10
	}

	return points
}

func (rs *ReceiptService) ValidateReceipt(receipt models.Receipt) error {
	if receipt.Retailer == "" || receipt.PurchaseDate == "" || receipt.PurchaseTime == "" || receipt.Total == "" || receipt.Total == "0" || len(receipt.Items) == 0 {
		return errors.New("invalid receipt data")
	}
	var sum float64
	for _, item := range receipt.Items {
		if item.ShortDescription == "" || item.Price == "" || item.Price == "0" {
			return errors.New("invalid item data")
		}
		price, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return errors.New("invalid item price data")
		}
		sum += price
	}
	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		return errors.New("invalid total price data")
	}

	if int(sum*100) != int(total*100) {
		return errors.New("sum of item prices must equal total")
	}
	return nil
}
