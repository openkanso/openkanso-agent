package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	stopCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pingController(client)

		case <-stopCtx.Done():
			log.Println("Stopping ping")
			return
		}
	}
}

func pingController(client *http.Client) {
	req, err := http.NewRequest(
		http.MethodGet,
		"http://127.0.0.1:9440/ping",
		nil,
	)
	if err != nil {
		log.Printf("Failed to build ping request: %v\n", err)
		return
	}

	req.Header.Add("X-NodeID", "abc123")

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
