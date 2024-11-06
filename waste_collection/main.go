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

		haToken := os.Getenv("SUPERVISOR_TOKEN")
		if haToken == "" {
			log.Fatalf("SUPERVISOR_TOKEN is missing. Ensure the add-on is configured correctly.")
		}

		haURL := os.Getenv("SUPERVISOR_API")
		if haURL == "" {
			haURL = "http://supervisor/core/api"
			fmt.Println("Set haURL to " + haURL)
		}
		haURL = haURL + "/states/sensor.waste_collection_ha"

		log.Println("Starting waste collection data fetch and push to Home Assistant...")

		urlStr := "https://www.simbio.si/sl/moj-dan-odvoza-odpadkov"

		data := url.Values{}
		data.Set("action", "simbioOdvozOdpadkov")
		data.Set("query", "začret 69") // Replace with the desired query

		req, err := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
		if err != nil {
			log.Fatalf("Error creating request to Simbio: %v", err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error sending request to Simbio: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Non-OK HTTP status from Simbio: %s", resp.Status)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response from Simbio: %v", err)
		}

		var wasteData []WasteData
		err = json.Unmarshal(body, &wasteData)
		if err != nil {
			log.Fatalf("Error parsing JSON from Simbio: %v", err)
		}

		for _, data := range wasteData {
			success := false
			for i := 0; i < 5; i++ {
				err = sendToHomeAssistant(data, haURL, haToken)
				if err == nil {
					log.Println("Data successfully sent to Home Assistant.")
					success = true
					break
				}
				log.Printf("Error sending data to Home Assistant (attempt %d/5): %v", i+1, err)
				time.Sleep(1 * time.Minute)
			}
			if !success {
				log.Println("Failed to send data to Home Assistant after 5 attempts. Sleeping for 1 hour...")
				time.Sleep(3600 * time.Second)
			}
		}

		fmt.Println("Sleeping for 1 hour...")
		time.Sleep(3600 * time.Second)
	}
}

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

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", haURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+haToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request to Home Assistant: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("error response from Home Assistant: %s - %s", resp.Status, string(bodyBytes))
	}

	log.Printf("Data successfully sent to Home Assistant. Address: %s, City: %s, Next MKO Pickup: %s, Next EMB Pickup: %s, Next BIO Pickup: %s",
		data.Name, data.City, data.NextMKO, data.NextEMB, data.NextBIO)
	return nil
}
