package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type STKPushPayload struct {
	BusinessShortCode string `json:"BusinessShortCode"`
	Password          string `json:"Password"`
	Timestamp         string `json:"Timestamp"`
	TransactionType   string `json:"TransactionType"`
	Amount            int    `json:"Amount"`
	PartyA            string `json:"PartyA"`
	PartyB            string `json:"PartyB"`
	PhoneNumber       string `json:"PhoneNumber"`
	CallBackURL       string `json:"CallBackURL"`
	AccountReference  string `json:"AccountReference"`
	TransactionDesc   string `json:"TransactionDesc"`
}

func main() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("env load error", err)
	}
	log.Println("env file loaded")
	// get token
	authToken, err := getToken(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	if err != nil {
		log.Fatal(err)
	}

	// Create an example request
	now := time.Now()
	timestamp := now.Format("20060102150405")
	shortCode := "174379"

	password := generatePassword(shortCode, os.Getenv("PASSKEY"), timestamp)

	requestPayload := STKPushPayload{
		BusinessShortCode: shortCode,
		Password:          password,
		Timestamp:         timestamp,
		TransactionType:   "CustomerPayBillOnline", //CustomerBuyGoodsOnline for till number
		Amount:            1,
		PartyA:            "254799962084",
		PartyB:            shortCode,
		PhoneNumber:       "254799962084",
		CallBackURL:       "https://mydomain.com/path",
		AccountReference:  "CompanyXLTD",
		TransactionDesc:   "Payment of X",
	}

	err = reqSTKPayment(authToken, requestPayload)
	if err != nil {
		panic(err)
	}
}

func getToken(consumerKey, consumerSecret string) (string, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", os.Getenv("TOKEN_URL"), nil)
	req.SetBasicAuth(consumerKey, consumerSecret)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	return result["access_token"].(string), nil
}

func reqSTKPayment(token string, request STKPushPayload) error {
	url := os.Getenv("URL")
	method := "POST"
	// Convert struct to JSON
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	// Convert JSON to io.Reader for request body
	reqBody := strings.NewReader(string(payload))

	// Create a new HTTP request

	req, err := http.NewRequest(method, url, reqBody)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	// Execute the HTTP request
	client := &http.Client{}
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
	return nil
}

func generatePassword(shortcode, passkey, timestamp string) string {
	rawPassword := shortcode + passkey + timestamp
	password := base64.StdEncoding.EncodeToString([]byte(rawPassword))
	return password
}

// {
// 	"MerchantRequestID": "9b85-4b30-ba72-7f3656486f92357920",
// 	"CheckoutRequestID": "ws_CO_06112024153808431799962084",
// 	"ResponseCode": "0",
// 	"ResponseDescription": "Success. Request accepted for processing",
// 	"CustomerMessage": "Success. Request accepted for processing"
//   }
