# Engineering Assignment – Campaign Events

## 1. Understanding of the Problem

The application receives campaign events from third-party
message delivery providers.

The events include sent, delivered, opened and clicked.

The backend must store events, prevent duplicate counting
and provide campaign statistics through REST APIs.

## 2. Assumptions

- Every event has a unique event_id.
- Only four event types are supported.
- Invalid events should be rejected.
- Valid events should be processed even if other events
  in the same batch are invalid.
- Event timestamps represent the actual event time.
- Events are stored in memory in the initial implementation.

## 3. Ambiguities

- Should duplicate events return success or conflict?
- Should statistics count events or unique contacts?
- Should events with the same ID but different payloads
  be treated as duplicates or conflicts?
- Should statistics include a date range?
- Should events be permanently stored?

## 4. Questions for the Product Manager

- Is event_id globally unique across providers?
- How should conflicting duplicate events be handled?
- Should statistics be based on event timestamps?
- What retention period is expected?

## 5. Implementation Priorities

1. Build POST /events.
2. Validate incoming events.
3. Prevent duplicate events.
4. Build campaign statistics API.
5. Test malformed, duplicate, and concurrent events.
6. Run automated tests, race detection, and static analysis.
7. Debug the provided Go program and document scaling considerations.
## 6. Scale Memo – 100 Million Events Per Day

### What breaks first?

- In-memory storage consumes too much RAM.
- A single server cannot handle all incoming traffic.
- Concurrent requests increase processing pressure.
- Restarting the application loses stored events.

### What would I change?

1. Introduce a message queue such as Kafka.
2. Use scalable persistent storage.
3. Add multiple backend instances.
4. Create database indexes for faster queries.
5. Separate event ingestion from statistics processing.

### Would the API contract change?

The existing API can remain mostly unchanged.
Internal processing and storage architecture would change.

### Would I introduce a queue?

Yes. A queue can handle incoming events and allow
background workers to process them.

New problems include:
- Message delivery failures.
- Duplicate processing.
- Processing delays.
- Queue backlog.

### Where can duplicates appear?

Duplicates can occur during provider retries,
queue redelivery, worker failures and database retries.

### How would I keep counters trustworthy?

- Use event_id for idempotency.
- Use database transactions.
- Reconcile aggregated statistics against stored events.
- Monitor duplicate and rejected event counts.

### What would I monitor?

- Events received per second.
- API response time.
- Error rate.
- Queue backlog.
- Database performance.
- Duplicate event count.
- Processing delays.

### What would I not solve yet?

I would not introduce complex distributed
infrastructure until actual traffic requires it.
 
## 7. Angry Marketer Scenario

### Possible Explanations

- Late events may have arrived.
- Duplicate events may have been removed.
- Invalid events may have been rejected.
- Bot-generated opens may have been filtered.
- Statistics may have been recalculated.
- Dashboard time filters may have changed.
- Previously counted events may have been corrected.
- Database aggregation or caching may have caused
  inconsistent results.

### What Would I Check First?

1. Check raw opened events at both timestamps.
2. Check duplicate event IDs.
3. Check event timestamps and arrival times.
4. Check dashboard filters.
5. Check database aggregation logic.
6. Check whether bot filtering or corrections occurred.

### Bug or Expected Behaviour?

If the dashboard displays lifetime counts of valid,
unique events in an append-only system, opened counts
should never decrease.

If events are removed, filtered or recalculated,
a decrease may be expected.

I would compare raw events with aggregated statistics
before declaring it a bug.

### Is Increasing Delivered Suspicious?

No. Providers can send delivery events late.
Therefore, delivered counts can increase even when
the campaign has already sent all messages.
## 8. What I Completed / What I Intentionally Skipped

### Completed
- Implemented POST /events and GET /campaigns/{campaign_id}/stats.
- Added event validation and duplicate detection.
- Implemented partial batch processing.
- Processed the provided seed data.
- Added six automated backend tests.
- Verified tests using the Go race detector.
- Ran go vet successfully.
- Fixed all four bugs in the debugging exercise.
- Verified debugging output against expected_output.txt.
- Documented assumptions, design decisions, and scaling considerations.

### Intentionally Skipped
- Persistent database storage.
- Event pagination.
- Frontend implementation.
- Distributed deployment and message queue infrastructure.

These were skipped to prioritize the core backend requirements within the available time.
## 9. If I Had Another Day

- Implement persistent storage using SQLite or PostgreSQL.
- Improve duplicate conflict handling when the same event ID arrives with different payloads.
- Add pagination for retrieving events.
- Improve API error responses and validation.
- Add structured logging and monitoring.
- Add more integration tests for API behaviour and edge cases.
- Test performance with larger event batches.
