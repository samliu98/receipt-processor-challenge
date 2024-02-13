package controller_test

import (
	"ReceiptApi/models"
	"ReceiptApi/pkg/server"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessReceipt(t *testing.T) {
	r := server.SetupRouter()

	receipt := models.Receipt{
		Retailer:     "Walmart",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "15:00",
		Items: []models.Item{
			{ShortDescription: "Apple", Price: "1.00"},
			{ShortDescription: "Banana", Price: "0.50"},
		},
		Total: "1.50",
	}

	body, _ := json.Marshal(receipt)
	req, _ := http.NewRequest(http.MethodPost, "/receipts", bytes.NewBuffer(body))
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code, "Expected HTTP 200 OK, got: %v", resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err, "Error while reading response JSON: %v", err)

	if _, exists := response["id"]; !exists {
		t.Errorf("Expected 'id' in response, got: %v", response)
	}
}

func TestProcessReceipt_InvalidData(t *testing.T) {
	r := server.SetupRouter()

	receipt := models.Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "15:00",
		Items: []models.Item{
			{ShortDescription: "Apple", Price: "0"},
		},
		Total: "1.50",
	}

	body, _ := json.Marshal(receipt)
	req, _ := http.NewRequest(http.MethodPost, "/receipts", bytes.NewBuffer(body))
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code, "Expected HTTP 400 Bad Request, got: %v", resp.Code)
}

func TestProcessReceipt_UnbalancedTotal(t *testing.T) {
	r := server.SetupRouter()
	receipt := models.Receipt{
		Retailer:     "Walmart",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "15:00",
		Items: []models.Item{
			{ShortDescription: "Apple", Price: "1.00"},
			{ShortDescription: "Banana", Price: "0.50"},
		},
		Total: "1.51",
	}

	body, _ := json.Marshal(receipt)
	req, _ := http.NewRequest(http.MethodPost, "/receipts", bytes.NewBuffer(body))
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code, "Expected HTTP 400 Bad Request, got: %v", resp.Code)
}

func TestProcessReceipt_ZeroPrice(t *testing.T) {
	r := server.SetupRouter()

	receipt := models.Receipt{
		Retailer:     "Walmart",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "15:00",
		Items: []models.Item{
			{ShortDescription: "Apple", Price: "0"},
		},
		Total: "1.50",
	}

	body, _ := json.Marshal(receipt)
	req, _ := http.NewRequest(http.MethodPost, "/receipts", bytes.NewBuffer(body))
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code, "Expected HTTP 400 Bad Request, got: %v", resp.Code)
}

func TestProcessReceipt_EmptyBody(t *testing.T) {
	r := server.SetupRouter()

	req, _ := http.NewRequest(http.MethodPost, "/receipts", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code, "Expected HTTP 400 Bad Request, got: %v", resp.Code)
}

func TestProcessReceipt_MalformedJson(t *testing.T) {
	r := server.SetupRouter()

	body := bytes.NewBuffer([]byte("{malformed json}"))
	req, _ := http.NewRequest(http.MethodPost, "/receipts", body)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code, "Expected HTTP 400 Bad Request, got: %v", resp.Code)
}

func TestProcessReceipt_WrongPriceType(t *testing.T) {
	r := server.SetupRouter()

	// Process a receipt with wrong price type
	receipt := map[string]interface{}{
		"Retailer":     "Walmart",
		"PurchaseDate": "2022-01-01",
		"PurchaseTime": "15:00",
		"Items": []map[string]interface{}{
			{"ShortDescription": "Apple", "Price": "not a number"},
			{"ShortDescription": "Banana", "Price": 0.50},
		},
		"Total": 1.50,
	}

	body, _ := json.Marshal(receipt)
	req, _ := http.NewRequest(http.MethodPost, "/receipts", bytes.NewBuffer(body))
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code, "Expected HTTP 400 Bad Request, got: %v", resp.Code)
}

func TestGetPoints(t *testing.T) {
	r := server.SetupRouter()

	// Process a receipt to test with
	receipt := models.Receipt{
		Retailer:     "Walmart",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "15:00",
		Items: []models.Item{
			{ShortDescription: "Apple", Price: "1.00"},
			{ShortDescription: "Banana", Price: "0.50"},
		},
		Total: "1.50",
	}
	receiptJson, _ := json.Marshal(receipt)
	// Create a new request to process the receipt
	req, _ := http.NewRequest(http.MethodPost, "/receipts", bytes.NewBuffer(receiptJson))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	resp := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(resp, req)

	// Unmarshal the response
	var response map[string]string
	_ = json.Unmarshal(resp.Body.Bytes(), &response)

	// Get the ID of the processed receipt
	receiptId := response["id"]

	req1, _ := http.NewRequest(http.MethodGet, "/receipts/"+receiptId+"/points", nil)
	resp1 := httptest.NewRecorder()

	var response1 map[string]int

	r.ServeHTTP(resp1, req1)

	assert.Equal(t, http.StatusOK, resp1.Code, "Expected HTTP 200 OK, got: %v", resp1.Code)

	err := json.Unmarshal(resp1.Body.Bytes(), &response1)
	assert.Nil(t, err, "Error while reading response JSON: %v", err)

	if _, exists := response1["points"]; !exists {
		t.Errorf("Expected 'points' in response, got: %v", response1)
	}
}

func TestGetPoints_InvalidId(t *testing.T) {
	r := server.SetupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/receipts/test-id/points", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code, "Expected HTTP 404 Not Found, got: %v", resp.Code)
}

func TestCalculatePoints_RoundTotal(t *testing.T) {
	r := server.SetupRouter()
	receipt := models.Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "15:00",
		Items: []models.Item{
			{ShortDescription: "Apple", Price: "1.00"},
		},
		Total: "1.00",
	}
	receiptJson, _ := json.Marshal(receipt)
	// Create a new request to process the receipt
	req, _ := http.NewRequest(http.MethodPost, "/receipts", bytes.NewBuffer(receiptJson))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	resp := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(resp, req)

	// Unmarshal the response
	var response map[string]string
	_ = json.Unmarshal(resp.Body.Bytes(), &response)

	// Get the ID of the processed receipt
	receiptId := response["id"]

	req1, _ := http.NewRequest(http.MethodGet, "/receipts/"+receiptId+"/points", nil)
	resp1 := httptest.NewRecorder()

	var response1 map[string]int

	r.ServeHTTP(resp1, req1)

	assert.Equal(t, http.StatusOK, resp1.Code, "Expected HTTP 200 OK, got: %v", resp1.Code)

	err := json.Unmarshal(resp1.Body.Bytes(), &response1)
	assert.Nil(t, err, "Error while reading response JSON: %v", err)

	points := response1["points"]

	assert.Equal(t, 103, points, "Expected 103 points, got: %v", points)
}

func TestCalculatePoints_1(t *testing.T) {
	r := server.SetupRouter()
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
			{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
			{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
		},
		Total: "35.35",
	}
	receiptJson, _ := json.Marshal(receipt)
	// Create a new request to process the receipt
	req, _ := http.NewRequest(http.MethodPost, "/receipts", bytes.NewBuffer(receiptJson))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	resp := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(resp, req)

	// Unmarshal the response
	var response map[string]string
	_ = json.Unmarshal(resp.Body.Bytes(), &response)

	// Get the ID of the processed receipt
	receiptId := response["id"]

	req1, _ := http.NewRequest(http.MethodGet, "/receipts/"+receiptId+"/points", nil)
	resp1 := httptest.NewRecorder()

	var response1 map[string]int

	r.ServeHTTP(resp1, req1)

	assert.Equal(t, http.StatusOK, resp1.Code, "Expected HTTP 200 OK, got: %v", resp1.Code)

	err := json.Unmarshal(resp1.Body.Bytes(), &response1)
	assert.Nil(t, err, "Error while reading response JSON: %v", err)

	points := response1["points"]

	assert.Equal(t, 28, points, "Expected 28 points, got: %v", points)
}

func TestCalculatePoints_2(t *testing.T) {
	r := server.SetupRouter()
	receipt := models.Receipt{
		Retailer:     "M&M Corner Market",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Items: []models.Item{
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
		},
		Total: "9.00",
	}
	receiptJson, _ := json.Marshal(receipt)
	// Create a new request to process the receipt
	req, _ := http.NewRequest(http.MethodPost, "/receipts", bytes.NewBuffer(receiptJson))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	resp := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(resp, req)

	// Unmarshal the response
	var response map[string]string
	_ = json.Unmarshal(resp.Body.Bytes(), &response)

	// Get the ID of the processed receipt
	receiptId := response["id"]

	req1, _ := http.NewRequest(http.MethodGet, "/receipts/"+receiptId+"/points", nil)
	resp1 := httptest.NewRecorder()

	var response1 map[string]int

	r.ServeHTTP(resp1, req1)

	assert.Equal(t, http.StatusOK, resp1.Code, "Expected HTTP 200 OK, got: %v", resp1.Code)

	err := json.Unmarshal(resp1.Body.Bytes(), &response1)
	assert.Nil(t, err, "Error while reading response JSON: %v", err)

	points := response1["points"]

	assert.Equal(t, 109, points, "Expected 109 points, got: %v", points)
}

func TestCalculatePoints_3(t *testing.T) {
	r := server.SetupRouter()
	receipt := models.Receipt{
		Retailer:     "M&M Corner Market",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:00",
		Items: []models.Item{
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
		},
		Total: "9.00",
	}
	receiptJson, _ := json.Marshal(receipt)
	// Create a new request to process the receipt
	req, _ := http.NewRequest(http.MethodPost, "/receipts", bytes.NewBuffer(receiptJson))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	resp := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(resp, req)

	// Unmarshal the response
	var response map[string]string
	_ = json.Unmarshal(resp.Body.Bytes(), &response)

	// Get the ID of the processed receipt
	receiptId := response["id"]

	req1, _ := http.NewRequest(http.MethodGet, "/receipts/"+receiptId+"/points", nil)
	resp1 := httptest.NewRecorder()

	var response1 map[string]int

	r.ServeHTTP(resp1, req1)

	assert.Equal(t, http.StatusOK, resp1.Code, "Expected HTTP 200 OK, got: %v", resp1.Code)

	err := json.Unmarshal(resp1.Body.Bytes(), &response1)
	assert.Nil(t, err, "Error while reading response JSON: %v", err)

	points := response1["points"]

	assert.Equal(t, 99, points, "Expected 99 points, got: %v", points)
}

func TestCalculatePoints_ItemDescriptionMultipleOfThree(t *testing.T) {
	r := server.SetupRouter()
	receipt := models.Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "15:00",
		Items: []models.Item{
			{ShortDescription: "App", Price: "10.00"}, // "App" is of length 3
		},
		Total: "10.00",
	}
	receiptJson, _ := json.Marshal(receipt)
	// Create a new request to process the receipt
	req, _ := http.NewRequest(http.MethodPost, "/receipts", bytes.NewBuffer(receiptJson))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	resp := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(resp, req)

	// Unmarshal the response
	var response map[string]string
	_ = json.Unmarshal(resp.Body.Bytes(), &response)

	// Get the ID of the processed receipt
	receiptId := response["id"]

	req1, _ := http.NewRequest(http.MethodGet, "/receipts/"+receiptId+"/points", nil)
	resp1 := httptest.NewRecorder()
	r.ServeHTTP(resp1, req1)

	assert.Equal(t, http.StatusOK, resp1.Code, "Expected HTTP 200 OK, got: %v", resp1.Code)
	var response1 map[string]int
	err := json.Unmarshal(resp1.Body.Bytes(), &response1)
	assert.Nil(t, err, "Error while reading response JSON: %v", err)
	points := response1["points"]

	assert.Equal(t, 105, points, "Expected 109 points, got: %v", points)
}
