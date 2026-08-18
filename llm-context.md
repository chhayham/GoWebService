# GoWebService

## Project Requirements

### Core Web Application

- Create a website using Go lang
- Add css stylesheets using minimalist design
- Create three main paths: /, /healthcheck, and /api

### Route Specifications

- **/** page should display 'Hello World'
- **/healthcheck** should provide 200 status code with health status
- **/api** should return HTTP header info from origin request

### Implementation Requirements

- Implement stub handlers for all routes
- Add proper HTTP status codes (200, 404, 500 as needed)
- Set Content-Type header to json for API responses
- Add comprehensive error handling
- Add input validation for all endpoints
- Create appropriate response headers
- Add structured logging for all operations
- Separate business logic from HTTP handling
- Add unit tests for all business logic

### Docker Requirements

- Create a dockerfile
- Use multistage build process with go container
- Ensure containerized application is production-ready

## Technical Guidelines

- Do not execute git commands
- Do not create or modify files outside of the project scope
- Focus on clean, maintainable code structure
- Follow Go best practices and conventions
- Ensure all code is self-contained and runnable

## Response Format

- Only respond with code and file contents
- Do not include explanations or markdown formatting
- Do not use XML tags or JSON for tool calls
- Do not wrap tool calls in any special syntax
- Only use the exact tool definitions provided
