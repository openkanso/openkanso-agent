package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	serverBaseURL, ok := os.LookupEnv("SERVER_BASE_URL")
	if !ok {
		log.Fatal("Missing required environment variable: SERVER_BASE_URL")
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	stopCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pingServer(client, serverBaseURL)

		case <-stopCtx.Done():
			log.Println("Stopping ping")
			return
		}
	}
}

func pingServer(client *http.Client, baseURL string) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/node/ping", baseURL),
		nil,
	)
	if err != nil {
		log.Printf("Failed to build ping request: %v\n", err)
		return
	}

	req.Header.Set("X-NodeID", "abc123")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "abc123"))

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send ping request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Bad ping response status: %v\n", resp.Status)
		return
	}

	log.Println("Ping succeeded")
}
