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

I followed this implementation order to establish correctness before considering scalability.

1. **Understand the event requirements:** Identify the four event types, required fields, and expected API behaviour before writing code.

2. **Implement POST /events:** Build the event ingestion endpoint first because it is the entry point for receiving provider events.

3. **Validate incoming events:** Check required identifiers, supported event types, and timestamps to prevent invalid data from entering the system.

4. **Implement duplicate detection:** Use event IDs to prevent provider retries from increasing campaign statistics more than once.

5. **Implement campaign statistics:** Build GET /campaigns/{campaign_id}/stats using accepted, unique event records.

6. **Test correctness and concurrency:** Test valid events, invalid events, duplicate events, campaign statistics, and concurrent requests. Use the Go race detector to identify unsafe shared access.

7. **Complete debugging and documentation:** Fix the four supplied bugs, compare the output against the expected result, and document assumptions, limitations, and scaling considerations.

This order prioritizes reliable event ingestion and accurate statistics before introducing more complex infrastructure.
## 6. Scale Memo – 100 Million Events Per Day

### What Breaks First?

At 100 million events per day, the current in-memory implementation would not be sufficient.

- Memory usage would increase continuously because events are stored in an in-memory map.
- A single server could become a bottleneck during peak traffic.
- Application restarts would cause data loss.
- Campaign statistics would become slower as the number of stored events increases.
- Concurrent requests would create additional database and processing pressure in a production system.

100 million events per day is approximately 1,157 events per second on average. Actual peak traffic could be significantly higher.

### What Changes Would I Make and In What Order?

1. **Introduce persistent storage:** Move events from in-memory storage to a database so data survives restarts.
2. **Add database indexing:** Index event IDs, campaign IDs, and timestamps to improve duplicate detection and statistics queries.
3. **Introduce a message queue:** Use Kafka or a similar system to absorb traffic spikes and separate ingestion from processing.
4. **Add background workers:** Process queued events asynchronously and update campaign statistics.
5. **Scale horizontally:** Run multiple API instances behind a load balancer.
6. **Optimize statistics:** Maintain pre-aggregated campaign counters instead of calculating every statistic from all raw events.
7. **Introduce retention and archival:** Move older events to cheaper storage according to business and compliance requirements.

### Would the API Contract Change?

The existing POST /events and GET /campaigns/{campaign_id}/stats endpoints can remain mostly unchanged.

However, I would define clear behaviour for:
- Partial batch acceptance.
- Duplicate event IDs.
- Conflicting payloads with the same event ID.
- Retryable failures.
- Event processing delays.

If processing becomes asynchronous, the POST response should clearly communicate whether an event has been accepted for processing or fully processed.

### Would Storage Design Change?

Yes.

I would use persistent database storage instead of an in-memory map.

The design would include:
- A unique constraint on event_id to enforce idempotency.
- Indexes on campaign_id and timestamp.
- Partitioning by event date when the event volume requires it.
- Pre-aggregated campaign statistics.
- Retention and archival policies for old events.

Raw events should remain the source of truth, while aggregated counters should be reconcilable against them.

### Would I Introduce a Message Queue?

Yes.

I would introduce Kafka or a similar message queue between event ingestion and background processing.

The API would validate incoming events and publish accepted events to the queue. Background workers would consume events, store them, and update statistics.

### What Problems Would the Queue Introduce?

- Duplicate message delivery.
- Consumer failures.
- Processing delays.
- Queue backlog.
- Message ordering issues.
- Retry and dead-letter queue management.
- More operational complexity.

I would use idempotent processing so repeated queue messages do not increase statistics more than once.

### Where Can Duplicates Appear?

Duplicates can occur during:
- Provider webhook retries.
- Network timeouts.
- Queue redelivery.
- Worker crashes and restarts.
- Database transaction retries.

A unique event ID constraint and idempotent processing would help prevent duplicate counting.

### How Would I Keep Counters Trustworthy?

- Use event IDs for idempotency.
- Update stored events and counters transactionally where possible.
- Make counter updates safe against repeated processing.
- Reconcile aggregated statistics against raw events periodically.
- Monitor differences between raw event counts and aggregated counters.

### What Would I Monitor?

- Events received per second.
- API response time and error rate.
- Queue backlog and processing delay.
- Consumer failures and retry counts.
- Database query performance.
- Duplicate and rejected event counts.
- Storage growth.
- Campaign statistics processing delays.

### What Would I Deliberately Leave Unsolved?

I would avoid introducing complex distributed infrastructure before actual traffic requires it.

I would not immediately build multi-region active-active processing, custom distributed counter systems, or complicated real-time analytics.

I would introduce these only when measured traffic, availability requirements, or business needs justify the additional complexity. 
## 7. Angry Marketer Scenario

The marketer reports that delivered events increased from 970,000 to 975,000, while opened events decreased from 250,000 to 248,000.

These changes should be investigated before deciding whether they represent expected behaviour or a system bug.

### Why Could Opened Statistics Decrease?

Possible explanations include:

- Duplicate opened events may have been removed during deduplication.
- Bot-generated opens may have been filtered.
- Previously accepted events may have been corrected or deleted.
- Historical statistics may have been recalculated.
- Dashboard filters or reporting time windows may have changed.
- Aggregation or caching problems may have produced inconsistent results.

### Could Late or Out-of-Order Events Explain This?

Late events can arrive after the original campaign activity.

Out-of-order events can arrive in a different sequence from the actual event timestamps.

These events can change statistics for historical time periods when reports are recalculated.

However, late or out-of-order events alone should not decrease a cumulative lifetime opened count in an append-only system. A decrease requires an additional operation, such as filtering, deletion, correction, deduplication, or recalculation.

### Why Could Delivered Counts Increase?

Delivery notifications may arrive late because providers can retry webhook requests or experience delays.

Therefore, delivered counts can increase even after the campaign has finished sending messages.

### What Would I Check First?

1. Compare raw opened event records at both reporting times.
2. Check whether any event IDs were removed as duplicates.
3. Compare event timestamps with event ingestion times.
4. Check whether bot filtering or event correction rules changed.
5. Verify dashboard date filters and reporting time zones.
6. Inspect aggregation jobs, database queries, and cached statistics.
7. Check processing delays, retries, and failed background jobs.

### How Would I Distinguish Expected Behaviour from a Bug?

If the system intentionally removes invalid, duplicate, or bot-generated events, a decrease may be expected.

If the system promises append-only lifetime counts of valid, unique events, opened counts should never decrease without an explicit correction or deletion.

I would compare the raw event records with aggregated statistics and review processing logs before declaring the behaviour a bug.

### Conclusion

The increase in delivered counts can be explained by late-arriving provider notifications.

The decrease in opened counts requires further investigation because late or out-of-order events alone do not explain a reduction in cumulative lifetime counts.
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
