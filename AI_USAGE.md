# AI Usage Documentation

## AI Tool Used
ChatGPT

## How I Used AI
- To understand the assignment requirements and divide the work into manageable tasks.
- To understand Go concepts, HTTP routing, request handling, and JSON processing.
- To get guidance on implementing event validation, duplicate detection, and campaign statistics.
- To create and improve automated backend tests.
- To troubleshoot compilation errors, runtime issues, and test failures.
- To understand the four bugs in the debugging exercise and their fixes.
- To improve project documentation, including README.md and NOTES.md.

## Suggestions I Changed or Reconsidered
I chose in-memory storage instead of adding a database because I wanted to complete and test the core requirements within the available time.

I also reviewed the suggested implementation and tested the behaviour rather than accepting suggestions without verification.

## What AI Helped Me Understand
AI helped me understand:
- Why event IDs are important for duplicate detection.
- Why mutex protection is needed when multiple goroutines access shared data.
- Why UTC timestamps are useful for consistent daily statistics.
- How partial batch processing allows valid events to be accepted even when other events are invalid.
- How the Go race detector helps identify concurrency problems.

## AI-Generated Code That Needed Debugging
During the debugging exercise, I encountered a syntax error caused by an extra closing brace while editing the tracking function. I corrected the code, formatted it using gofmt, and verified the output against expected_output.txt.

I also reviewed and tested the backend implementation, including its validation and concurrent request handling.

## Verification
- Ran automated backend tests using go test -v ./...
- Ran the race detector using go test -race ./...
- Ran static analysis using go vet ./...
- Compared debugging output with expected_output.txt.
- Ran the debugging program with the Go race detector.
- Tested the service using POST and GET requests.
- Tested the service using the provided seed events.
