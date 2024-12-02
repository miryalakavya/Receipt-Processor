package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Receipt structure to hold the receipt data
type Receipt struct {
	RetName       string  `json:"retailer"`
	PurchaseDate  string  `json:"purchaseDate"`
	PurchaseTime  string  `json:"purchaseTime"`
	Items         []Item  `json:"items"`
	Total         string  `json:"total"`
}

// Item structure for each item on the receipt
type Item struct {
	ShortDesc string `json:"shortDescription"`
	Price     string `json:"price"`
}

// PointsResponse structure to return points in the response
type PointsResponse struct {
	Points int `json:"points"`
}

// ReceiptIDResponse structure to return the receipt ID in the response
type ReceiptIDResponse struct {
	ID string `json:"id"`
}

// In-memory storage for receipts
var receiptDB = make(map[string]Receipt)

func main() {
	// Add a log statement to confirm server start
	log.Println("Starting server on port 8080...")

	// Set up the HTTP routes
	http.HandleFunc("/receipts/process", ProcessReceipt) // POST /receipts/process
	http.HandleFunc("/receipts/", GetPoints)            // GET /receipts/{id}/points

	// Start the server on port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}


// ProcessReceipt handles the POST request to process a receipt
func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	log.Println("Received a request to process a receipt")

	var receipt Receipt
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&receipt); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate a unique ID for the receipt
	receiptID := uuid.New().String()
	// Save the receipt in the in-memory database
	receiptDB[receiptID] = receipt

	// Log receipt ID and the stored database
	log.Printf("Processed receipt with ID: %s\n", receiptID)
	log.Printf("Current receiptDB: %+v\n", receiptDB)

	// Return the ID in the response
	response := ReceiptIDResponse{ID: receiptID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


// GetPoints handles the GET request to fetch points for a receipt by ID
func GetPoints(w http.ResponseWriter, r *http.Request) {
    // Extract the receipt ID from the URL path, ignoring '/points' part
    path := strings.TrimPrefix(r.URL.Path, "/receipts/")
    receiptID := strings.TrimSuffix(path, "/points") // Remove '/points' from the end

    log.Println("Received request for receipt ID:", receiptID)

    // Find the receipt by ID in the in-memory database
    receipt, found := receiptDB[receiptID]
    if !found {
        log.Println("Receipt not found:", receiptID)
        http.Error(w, "Receipt not found", http.StatusNotFound)
        return
    }

    // Calculate the points for the receipt
    points := calculatePoints(receipt)

    log.Println("Points for receipt ID:", receiptID, "Points:", points)

    // Return the points in the response
    response := PointsResponse{Points: points}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}


// calculatePoints calculates the points based on the receipt rules
func calculatePoints(receipt Receipt) int {
	points := 0

	// Rule 1: Points for alphanumeric characters in retailer name
	points += len(receipt.RetName)

	// Rule 2: 50 points if total is a round dollar amount
	total, _ := strconv.ParseFloat(receipt.Total, 64)
	if total == math.Floor(total) {
		points += 50
	}

	// Rule 3: 25 points if total is a multiple of 0.25
	if math.Mod(total, 0.25) == 0 {
		points += 25
	}

	// Rule 4: 5 points for every 2 items on the receipt
	points += len(receipt.Items) / 2 * 5

	// Rule 5: Price calculation for item description length multiple of 3
	for _, item := range receipt.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDesc)
		if len(trimmedDesc)%3 == 0 {
			itemPrice, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(itemPrice * 0.2))
		}
	}

	// Rule 6: 6 points if the day of the purchase date is odd
	purchaseDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if purchaseDate.Day()%2 != 0 {
		points += 6
	}

	// Rule 7: 10 points if the purchase time is between 2:00pm and 4:00pm
	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	if purchaseTime.Hour() >= 14 && purchaseTime.Hour() < 16 {
		points += 10
	}

	return points
}
