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

// Supervisor API endpoint
const haURL = "http://supervisor/core/api/states/sensor.waste_collection"

func main() {
	// Get the Supervisor token from environment variables
	haToken := os.Getenv("SUPERVISOR_TOKEN")
	if haToken == "" {
		log.Fatalf("SUPERVISOR_TOKEN is missing. Ensure the add-on is configured correctly.")
	}

	// Define the Simbio URL
	urlStr := "https://www.simbio.si/sl/moj-dan-odvoza-odpadkov"

	// Create form data for the POST request to Simbio
	data := url.Values{}
	data.Set("action", "simbioOdvozOdpadkov")
	data.Set("query", "začret 69") // Replace with your desired query

	// Create a POST request
	req, err := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Set headers for Simbio request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request to Simbio: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Non-OK HTTP status from Simbio: %s", resp.Status)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response from Simbio: %v", err)
	}

	// Parse the JSON response from Simbio
	var wasteData []WasteData
	err = json.Unmarshal(body, &wasteData)
	if err != nil {
		log.Fatalf("Error parsing JSON from Simbio: %v", err)
	}

	// Send each piece of data to Home Assistant
	for _, data := range wasteData {
		err = sendToHomeAssistant(data, haToken)
		if err != nil {
			log.Printf("Error sending data to Home Assistant: %v", err)
		}
	}
}

// Function to send data to Home Assistant
func sendToHomeAssistant(data WasteData, haToken string) error {
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

	// Marshal the payload into JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Create a POST request to the Home Assistant API
	req, err := http.NewRequest("POST", haURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set the Authorization header with the Supervisor token
	req.Header.Set("Authorization", "Bearer "+haToken)
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request to Home Assistant: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("error response from Home Assistant: %s - %s", resp.Status, string(bodyBytes))
	}

	log.Println("Data successfully sent to Home Assistant!")
	return nil
}
