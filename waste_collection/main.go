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
	"time"
)

// Struct to hold the response data from Simbio
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
	for {
		fmt.Println("##### RESTARTING #######")
		fmt.Println("Environment Variables:")
		for _, e := range os.Environ() {
			fmt.Println(e)
		}

		// Retrieve the Supervisor token from the environment variable
		haToken := os.Getenv("SUPERVISOR_TOKEN")
		if haToken == "" {
			log.Fatalf("SUPERVISOR_TOKEN is missing. Ensure the add-on is configured correctly.")
		}

		// Use the Supervisor API URL directly if it's not set as an environment variable
		haURL := os.Getenv("SUPERVISOR_API")
		if haURL == "" {
			haURL = "http://supervisor/core/api"
			fmt.Println("Set haURL to" + haURL)
		}
		haURL = haURL + "/states/sensor.waste_collection_ha"

		log.Println("Starting waste collection data fetch and push to Home Assistant...")

		// Define the URL for the waste collection service
		urlStr := "https://www.simbio.si/sl/moj-dan-odvoza-odpadkov"

		// Create form data for the POST request to Simbio
		data := url.Values{}
		data.Set("action", "simbioOdvozOdpadkov")
		data.Set("query", "začret 69") // Replace with the desired query

		// Create a POST request to Simbio
		req, err := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
		if err != nil {
			log.Fatalf("Error creating request to Simbio: %v", err)
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
			err = sendToHomeAssistant(data, haURL, haToken)
			if err != nil {
				log.Printf("Error sending data to Home Assistant: %v", err)
			}
		}

		// Wait for an hour before running the next scrape
		fmt.Println("Sleeping for 1 hour...")
		time.Sleep(3600 * time.Second)

	}
}

// Function to send data to Home Assistant
func sendToHomeAssistant(data WasteData, haURL, haToken string) error {
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

	// Log the data sent to Home Assistant
	log.Printf("Data successfully sent to Home Assistant!")
	log.Printf("State: updated")
	log.Printf("Address: %s", data.Name)
	log.Printf("City: %s", data.City)
	log.Printf("Next MKO Pickup: %s", data.NextMKO)
	log.Printf("Next EMB Pickup: %s", data.NextEMB)
	log.Printf("Next BIO Pickup: %s", data.NextBIO)
	log.Printf("Friendly Name: Waste Collection")
	return nil
}
