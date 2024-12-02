# Receipt Processor API

## Description

A **Go-based Receipt Processor API** that processes receipts and calculates points based on specific business rules. This API provides two main endpoints:
- **POST /receipts/process** to process a receipt.
- **GET /receipts/{id}/points** to retrieve the points for a specific receipt.

## Key Features:
- Submit a receipt in JSON format to get a unique ID.
- Retrieve the points associated with a receipt using its ID.
- Points calculation is based on several rules like retailer name length, total price, item descriptions, and more.

## How to Run Locally

1. Clone the repository:
   git clone https://github.com/miryalakavya/Receipt-Processor.git
