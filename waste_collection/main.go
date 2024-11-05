package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

// Struct to hold the response data
type WasteData struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Query   string `json:"query"`
	City    string `json:"city"`
	NextMKO string `json:"next_mko"`
	NextEMB string `json:"next_emb"`
	NextBIO string `json:"next_bio"`
}

func main() {
	// Define the URL for waste collection data
	urlStr := "https://www.simbio.si/sl/moj-dan-odvoza-odpadkov"

	// Create form data
	data := url.Values{}
	data.Set("action", "simbioOdvozOdpadkov")
	data.Set("query", "začret 69") // Replace with your desired query

	// Create a POST request
	req, err := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Set headers to match the original request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Non-OK HTTP status: %s", resp.Status)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	// Parse the JSON response
	var wasteData []WasteData
	err = json.Unmarshal(body, &wasteData)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Send data to Home Assistant
	for _, data := range wasteData {
		err = sendToHomeAssistant(data)
		if err != nil {
			log.Printf("Error sending data to Home Assistant: %v", err)
		}
	}
}

// Function to send data to Home Assistant
func sendToHomeAssistant(data WasteData) error {
	// Fetch Home Assistant URL and Token from environment variables
	haURL := os.Getenv("HOMEASSISTANT_URL")
	haToken := os.Getenv("SUPERVISOR_TOKEN")

	if haURL == "" || haToken == "" {
		return fmt.Errorf("missing Home Assistant URL or token")
	}

	payload := map[string]interface{}{
		"state": "updated",
		"attributes": map[string]string{
			"address":       data.Name,
			"city":          data.City,
			"next_mko":      data.NextMKO,
			"next_emb":      data.NextEMB,
			"next_bio":      data.NextBIO,
			"friendly_name": "Waste Collection",
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", haURL+"/api/states/sensor.waste_collection", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set the Authorization header
	req.Header.Set("Authorization", "Bearer "+haToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("error response from Home Assistant: %s - %s", resp.Status, string(bodyBytes))
	}

	return nil
}
