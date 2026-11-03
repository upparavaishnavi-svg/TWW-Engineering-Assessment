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
I considered adding persistent database storage, such as SQLite. However, I chose in-memory storage because I wanted to complete and test the core requirements within the available time. I understand that persistent storage would be necessary for a production system.

I also reviewed the suggested implementation and tested its behaviour rather than accepting suggestions without verification.

## What AI Helped Me Understand
AI helped me understand:
- Why event IDs are important for duplicate detection.
- Why mutex protection is needed when multiple goroutines access shared data.
- Why UTC timestamps are useful for consistent daily statistics.
- How partial batch processing allows valid events to be accepted even when other events are invalid.
- How the Go race detector helps identify concurrency problems.

## AI-Assisted Code That Needed Debugging
I used ChatGPT to help draft and understand parts of the Go backend implementation, including event validation, duplicate detection, campaign statistics, and automated tests.

During implementation, I encountered a syntax error caused by an extra closing brace while editing the tracking function. I corrected the brace, formatted the code using gofmt, and verified the program against the expected output.

I also reviewed the AI-assisted implementation, corrected issues encountered during integration, and verified the final behaviour using automated tests and the Go race detector.

## Verification
- Ran automated backend tests using go test -v ./...
- Ran the race detector using go test -race ./...
- Ran static analysis using go vet ./...
- Compared debugging output with expected_output.txt.
- Ran the debugging program with the Go race detector.
- Tested the service using POST and GET requests.
- Tested the service using the provided seed events.
