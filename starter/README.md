# Campaign Events Service

A Go HTTP API that receives campaign events from messaging providers and provides campaign-level statistics.

Supported event types:
- `sent`
- `delivered`
- `opened`
- `clicked`

## Requirements

- Go 1.22 or later.
- No third-party dependencies are required.

## Run the Application

Open a terminal in the `starter` directory and run:

```bash
go run .
```

The server starts on:

```text
http://localhost:8080
```

## API Endpoints

### 1. POST /events

Accepts a JSON array of campaign events.

Example request:

```bash
curl -s -X POST localhost:8080/events \
  -H 'Content-Type: application/json' \
  -d '[{"event_id":"evt_1","campaign_id":"cmp_summer_sale","contact_id":"ct_001","type":"delivered","timestamp":"2026-08-10T06:15:00Z"}]'
```

Example response:

```json
{
  "accepted": 1,
  "duplicates": 0,
  "rejected": 0,
  "errors": []
}
```

**Validation rules:**
- `event_id`, `campaign_id`, and `contact_id` are required.
- Event type must be one of the four supported types.
- Timestamp must be in RFC3339 format.
- The request body must contain a JSON array.
- Request body size is limited to 2 MiB.

**Batch handling:**
- Valid events are accepted even if other events in the batch are invalid.
- Repeated event IDs are counted as duplicates.
- Invalid events are rejected, with error details returned in the response.
- A repeated event ID is treated as a duplicate even if its other fields differ.

### 2. GET /campaigns/{campaign_id}/stats

Returns event counts for the requested campaign.

Example:

```bash
curl -s localhost:8080/campaigns/cmp_summer_sale/stats
```

Example response:

```json
{
  "campaign_id": "cmp_summer_sale",
  "stats": {
    "sent": 10,
    "delivered": 8,
    "opened": 5,
    "clicked": 2
  }
}
```

The statistics count stored event records by event type. They do not represent unique contacts.

## Storage and Limitations

- Events are stored in an in-memory Go map.
- Event IDs are used as keys for duplicate detection.
- A mutex protects shared event data during concurrent requests.
- Events are not permanently stored.
- All stored events are lost when the application restarts.
- The implementation does not include pagination or persistent database storage.

## Running Tests

Run the automated tests:

```bash
go test -v ./...
```

Run tests with the Go race detector:

```bash
go test -race ./...
```

Run static analysis:

```bash
go vet ./...
```

The tests cover valid events, invalid events, duplicate detection, campaign statistics, concurrent duplicate requests, and rejection of a `null` batch.

## Seed Data

`seed/events.json` contains 235 event records, including retries, conflicts, and malformed records.

To load the seed data, first start the application and then run this command from the `starter` directory in another terminal:

```bash
curl -s -X POST localhost:8080/events \
  -H 'Content-Type: application/json' \
  --data-binary @seed/events.json
```

The service processes valid events and reports accepted, duplicate, and rejected counts.

## Design Decisions

- Used in-memory storage to keep the initial implementation simple.
- Used event IDs for duplicate detection.
- Used partial batch processing so one invalid event does not prevent valid events from being accepted.
- Used a read-write mutex to protect shared data.
- Accepted late and out-of-order events without assuming events arrive in timestamp order.
- Used event timestamps only for validation; campaign statistics count stored events.
