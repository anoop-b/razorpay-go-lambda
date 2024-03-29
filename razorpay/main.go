package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	razorpay "github.com/razorpay/razorpay-go"
)

// Customer struct from frontend
type Customer struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Contact   string `json:"contact"`
	MaxAmount int    `json:"max_amount"`
	Plan      string `json:"plan"`
	Pan       string `json:"pan"`
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var payload Customer
	if err := json.Unmarshal([]byte(request.Body), &payload); err != nil {
		log.Fatal("[ERROR]", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 406}, err
	}
	response := Rzpay(payload)
	// convert map to json encoded bytes
	jsonResponse, err := json.Marshal(response)
	// check for error in json encoding
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(jsonResponse)

	return events.APIGatewayProxyResponse{Body: string(jsonResponse), StatusCode: 200}, nil
}

// Rzpay function handles API calls to razorpay
func Rzpay(payload Customer) map[string]interface{} {
	key := os.Getenv("KEY")
	secret := os.Getenv("SECRET")
	client := razorpay.NewClient(key, secret)

	finalPayload := map[string]interface{}{
		"customer": map[string]interface{}{
			"name":    payload.Name,
			"email":   payload.Email,
			"contact": payload.Contact,
		},
		"type":        "link",
		"amount":      0,
		"currency":    "INR",
		"description": "Donation for Internet Freedom Foundation",
		"subscription_registration": map[string]interface{}{
			"method":               "emandate",
			"first_payment_amount": payload.MaxAmount,
			"max_amount":           payload.MaxAmount,
		},
	}

	b, err := client.Subscription.Create(finalPayload, nil)
	if err != nil {
		fmt.Println(err)
	}
	return b
}

func main() {
	lambda.Start(handleRequest)
}
