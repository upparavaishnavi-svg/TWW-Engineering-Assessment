# AI Usage Documentation

## AI Tool Used
ChatGPT

## How I Used AI
- To understand the assignment requirements and break the work into smaller tasks.
- To understand Go concepts and the HTTP API implementation.
- To get help with event validation, duplicate detection, and campaign statistics.
- To troubleshoot errors while running and testing the Go programs.
- To understand the four bugs in the debugging exercise.

## Suggestions I Changed or Reconsidered
I chose in-memory storage for the service to keep the implementation simple and complete the core requirements within the available time. I understand that persistent storage would be necessary for a production system.

I also reviewed the suggested code and made corrections while implementing and testing it rather than blindly accepting every suggestion.

## What AI Helped Me Understand
AI helped me understand why duplicate events must be identified using event IDs, why shared counters need mutex protection when goroutines update them, and why UTC should be used for consistent daily statistics.

## AI-Generated Code That Needed Debugging
During the debugging exercise, I encountered a syntax error caused by an extra closing brace while editing the tracking function. I corrected the code, formatted it using gofmt, and verified the result against the expected output.

## Verification
- Compared debugging output with expected_output.txt.
- Ran the debugging program with the Go race detector.
- Tested the service using POST and GET requests.
- Tested the service using the provided seed events.
