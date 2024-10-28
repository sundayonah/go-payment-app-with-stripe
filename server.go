package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/paymentintent"
)

// "The maximum lenght of a function is inversely proportional to the complexity and indentation level of that function"
// Never try to explain How your code works in a comment
// sk_test_51OJVxhBfYy21A59VIBTOCXFU0UvhJ2UnQgDRPy4OWAR872vq2gGwD3sNhpIpUs1vJGSd5fsusfxZ0DsVbfLs6gxd00fmuWFMwM

func main() {

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	http.HandleFunc("/create-payment-intent", handleCreatePaymentIntent)
	http.HandleFunc("/health", handleHealth)

	log.Println("Listening on http://localhost:4242...")
	var err = http.ListenAndServe("localhost:4242", nil)
	if err != nil {
		log.Fatal("Something went wrong")
	}

}

func handleCreatePaymentIntent(w http.ResponseWriter, r *http.Request) {

	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle OPTIONS request for preflight
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ProductId string `json:"product_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Address1  string `json:"address_1"`
		Address2  string `json:"address_2"`
		City      string `json:"city"`
		State     string `json:"state"`
		Zip       string `json:"zip"`
		Country   string `json:"country"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(calculateOrderAmount(req.ProductId)), // $10.99
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}
	paymentIntent, err := paymentintent.New(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(paymentIntent.ClientSecret)

	var response struct {
		ClientSecret string `json:"clientSecret"`
	}
	response.ClientSecret = paymentIntent.ClientSecret

	var buf bytes.Buffer

	err = json.NewEncoder(&buf).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	_, err = io.Copy(w, &buf)
	if err != nil {
		log.Println("Error writing response")
	}

}

func handleHealth(w http.ResponseWriter, r *http.Request) {

	response := []byte("Server is up and running!")

	_, err := w.Write(response)

	if err != nil {
		log.Println("Error writing response")
	}

}

func calculateOrderAmount(productId string) int64 {
	// Assuming product prices are stored in a database or some other storage
	// and retrieved based on the provided productId.
	// For demonstration purposes, let's assume a simple price lookup.
	switch productId {
	case "Forever Pants":
		return 1099
	case "Forever Shirt":
		return 2199
	case "Forever Shorts":
		return 30000
	default:
		return 0
	}
}
