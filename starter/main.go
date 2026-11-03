package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	events = make(map[string]Event)
	mu     sync.RWMutex
)

type BatchResponse struct {
	Accepted   int      `json:"accepted"`
	Duplicates int      `json:"duplicates"`
	Rejected   int      `json:"rejected"`
	Errors     []string `json:"errors"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /events", handlePostEvents)
	mux.HandleFunc("GET /campaigns/{campaignID}/stats", handleGetStats)

	addr := ":8080"
	fmt.Printf("listening on %s\n", addr)

	log.Fatal(http.ListenAndServe(addr, mux))
}

func handlePostEvents(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, 2<<20)

	body, err := io.ReadAll(r.Body)

	if err != nil {
		writeJSON(w, http.StatusBadRequest,
			map[string]string{"error": "invalid request body"})
		return
	}

	var rawEvents []json.RawMessage

	err = json.Unmarshal(body, &rawEvents)

	if err != nil {
		writeJSON(w, http.StatusBadRequest,
			map[string]string{"error": "expected JSON array"})
		return
	}

	result := BatchResponse{
		Errors: []string{},
	}

	for i, raw := range rawEvents {

		var event Event

		err := json.Unmarshal(raw, &event)

		if err != nil {
			result.Rejected++
			result.Errors = append(result.Errors,
				fmt.Sprintf("event %d: invalid JSON", i))
			continue
		}

		if err := validateEvent(event); err != nil {
			result.Rejected++
			result.Errors = append(result.Errors,
				fmt.Sprintf("event %d: %s", i, err))
			continue
		}

		mu.Lock()

		if _, exists := events[event.EventID]; exists {
			result.Duplicates++
			mu.Unlock()
			continue
		}

		events[event.EventID] = event
		result.Accepted++

		mu.Unlock()
	}

	writeJSON(w, http.StatusOK, result)
}

func validateEvent(e Event) error {

	if strings.TrimSpace(e.EventID) == "" {
		return fmt.Errorf("event_id is required")
	}

	if strings.TrimSpace(e.CampaignID) == "" {
		return fmt.Errorf("campaign_id is required")
	}

	if strings.TrimSpace(e.ContactID) == "" {
		return fmt.Errorf("contact_id is required")
	}

	switch e.Type {
	case "sent", "delivered", "opened", "clicked":
	default:
		return fmt.Errorf("invalid event type")
	}

	if _, err := time.Parse(time.RFC3339, e.Timestamp); err != nil {
		return fmt.Errorf("invalid timestamp")
	}

	return nil
}

func handleGetStats(w http.ResponseWriter, r *http.Request) {

	campaignID := r.PathValue("campaignID")

	stats := map[string]int{
		"sent":      0,
		"delivered": 0,
		"opened":    0,
		"clicked":   0,
	}

	mu.RLock()

	for _, event := range events {

		if event.CampaignID == campaignID {
			stats[event.Type]++
		}
	}

	mu.RUnlock()

	writeJSON(w, http.StatusOK, map[string]any{
		"campaign_id": campaignID,
		"stats":       stats,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON: %v", err)
	}
}
