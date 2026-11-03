# Debugging Exercise – Four Bugs

## Bug 1: Duplicate Detection

Problem:
Duplicate events were counted again when they appeared
in different batches.

Reason:
The seen map was created inside processBatch().

Fix:
Moved event ID tracking to a global map named
seenEventIDs.

Verification:
Compared actual.txt with expected_output.txt.
The outputs matched exactly.

## Bug 2: Unique Opens

Problem:
Unique opens were counted globally by contact ID.

Reason:
The same contact opening in two different campaigns
was counted only once.

Fix:
Used campaign ID and contact ID together as the key.

Verification:
Unique open counts matched expected output.

## Bug 3: UTC Date Calculation

Problem:
Daily delivered statistics used local time.

Reason:
Timestamp.Local() converted timestamps into
the system timezone.

Fix:
Replaced Local() with UTC().

Verification:
Daily delivered statistics matched expected output.

## Bug 4: Race Condition

Problem:
Multiple goroutines updated shared campaign counters
without synchronization.

Reason:
The apply() function modified shared integers
concurrently.

Fix:
Added sync.Mutex to protect counter updates.

Verification:
Ran go run -race . events.jsonl.
No data races were reported.
The output matched expected_output.txt.
