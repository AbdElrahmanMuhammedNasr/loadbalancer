package routes

import (
	"context"
	"encoding/json"
	"errors"
	"loadbalancer/db"
	"loadbalancer/dto"
	"log"
	"net/http"
	"strings"
)

func ProcessKeysAndSendRequests() {
	ctx := context.Background()
	keys, err := db.RedisClient.Keys(ctx, "*").Result()
	if err != nil {
		log.Println("Error fetching keys:", err)
		return
	}

	// Iterate over each key
	for _, key := range keys {
		value, err := db.RedisClient.Get(ctx, key).Result()
		if err != nil {
			log.Println("Error fetching value for key:", key, err)
			continue
		}

		// Deserialize JSON if the value is a JSON string
		var urls []dto.ServiceUrlDTO
		if err := json.Unmarshal([]byte(value), &urls); err != nil {
			log.Println("Error decoding JSON for key:", key, err)
			continue
		}
		// Send request for each URL
		for _, url := range urls {
			go sendRequest(url) // Use goroutine for concurrent execution
		}
	}
}

// Function to send an HTTP request
func sendRequest(url dto.ServiceUrlDTO) (dto.ServiceUrlDTO, error) {
	if !strings.HasPrefix(url.Url, "http") {
		log.Println("Skipping invalid URL:", url)
		return url, errors.New("invalid URL format")
	}

	resp, err := http.Get(url.Url)
	if err != nil {
		log.Println("Request failed for URL:", url, err)
		url.Active = false
		return url, err
	}
	url.Active = true

	defer resp.Body.Close()
	log.Printf("Request sent to %s, Status: %d\n", url, resp.StatusCode)
	return url, nil
}
