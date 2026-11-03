package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func resetEvents() {
	mu.Lock()
	defer mu.Unlock()
	events = make(map[string]Event)
}

func TestValidEvents(t *testing.T) {
	resetEvents()

	body := `[
		{
			"event_id":"evt1",
			"campaign_id":"cmp1",
			"contact_id":"ct1",
			"type":"sent",
			"timestamp":"2026-08-10T09:00:00Z"
		},
		{
			"event_id":"evt2",
			"campaign_id":"cmp1",
			"contact_id":"ct2",
			"type":"delivered",
			"timestamp":"2026-08-10T09:10:00Z"
		}
	]`

	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	handlePostEvents(rec, req)

	var result BatchResponse
	json.Unmarshal(rec.Body.Bytes(), &result)

	if result.Accepted != 2 || result.Rejected != 0 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestInvalidEvent(t *testing.T) {
	resetEvents()

	body := `[
		{
			"event_id":"evt1",
			"campaign_id":"cmp1",
			"contact_id":"ct1",
			"type":"unknown",
			"timestamp":"2026-08-10T09:00:00Z"
		}
	]`

	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	handlePostEvents(rec, req)

	var result BatchResponse
	json.Unmarshal(rec.Body.Bytes(), &result)

	if result.Rejected != 1 {
		t.Fatalf("expected one rejected event, got %+v", result)
	}
}

func TestDuplicateEvent(t *testing.T) {
	resetEvents()

	body := `[
		{
			"event_id":"evt1",
			"campaign_id":"cmp1",
			"contact_id":"ct1",
			"type":"opened",
			"timestamp":"2026-08-10T09:00:00Z"
		}
	]`

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		handlePostEvents(rec, req)

		var result BatchResponse
		json.Unmarshal(rec.Body.Bytes(), &result)

		if i == 0 && result.Accepted != 1 {
			t.Fatalf("first request should accept event")
		}

		if i == 1 && result.Duplicates != 1 {
			t.Fatalf("second request should detect duplicate")
		}
	}
}

func TestCampaignStats(t *testing.T) {
	resetEvents()

	events["evt1"] = Event{
		EventID: "evt1", CampaignID: "cmp1", Type: "sent",
	}

	events["evt2"] = Event{
		EventID: "evt2", CampaignID: "cmp1", Type: "opened",
	}

	req := httptest.NewRequest(http.MethodGet, "/campaigns/cmp1/stats", nil)
	req.SetPathValue("campaignID", "cmp1")

	rec := httptest.NewRecorder()

	handleGetStats(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestConcurrentDuplicateRequests(t *testing.T) {
	resetEvents()

	body := `[
		{
			"event_id":"evt1",
			"campaign_id":"cmp1",
			"contact_id":"ct1",
			"type":"sent",
			"timestamp":"2026-08-10T09:00:00Z"
		}
	]`

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(body))
			rec := httptest.NewRecorder()

			handlePostEvents(rec, req)
		}()
	}

	wg.Wait()

	mu.RLock()
	defer mu.RUnlock()

	if len(events) != 1 {
		t.Fatalf("expected exactly one stored event, got %d", len(events))
	}
}
func TestNullBatchRejected(t *testing.T) {
	resetEvents()

	req := httptest.NewRequest(
		http.MethodPost,
		"/events",
		bytes.NewBufferString("null"),
	)

	rec := httptest.NewRecorder()
	handlePostEvents(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for null batch, got %d", rec.Code)
	}
}
