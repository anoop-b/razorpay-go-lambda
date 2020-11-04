package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

// ErpNext payload struct
type ErpNext struct {
	CustomerID string `json:"customer_id"`
	Plan       string `json:"plan"`
	Pan        string `json:"pan"`
}

// RazorpayResponse to extract customer id only
type RazorpayResponse struct {
	CustomerID string `json:"customer_id"`
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var payload Customer
	json.Unmarshal([]byte(request.Body), &payload)
	response := Rzpay(payload)
	// convert map to json encoded bytes
	jsonResponse, err := json.Marshal(response)
	// check for error in json encoding
	if err != nil {
		fmt.Println(err.Error())
	}

	erpNext(payload, jsonResponse)

	return events.APIGatewayProxyResponse{Body: string(jsonResponse), StatusCode: 200, Headers: map[string]string{
		"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,x-requested-with",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "POST,GET,OPTIONS",
	}}, nil
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

func erpNext(payload Customer, response []byte) {
	var razorpayResponse RazorpayResponse
	token := os.Getenv("ERP_KEY")
	secret := os.Getenv("ERP_TOKEN")

	err := json.Unmarshal(response, &razorpayResponse)
	if err != nil {
		fmt.Println("error Unmarshalling razorpay response", err)
	}

	data := ErpNext{
		razorpayResponse.CustomerID,
		payload.Plan,
		payload.Pan,
	}

	jsonResponse, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error Marshalling ErpNext payload", err)
	}

	body := bytes.NewReader(jsonResponse)
	req, err := http.NewRequest("POST", "https://iff.erpnext.com/api/method/iff.api.create_member", body)
	if err != nil {
		fmt.Println("error creating request", err)
	}

	// req.Header.Set("Authorization", token)
	req.SetBasicAuth(token, secret) //fallback if auth header doesn't work
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error sending post request", err)
	}

	defer resp.Body.Close()
}

func main() {
	lambda.Start(handleRequest)
}
