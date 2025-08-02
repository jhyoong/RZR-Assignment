# Razer Take Home Assignment

A web application that allows users to check if their email addresses have been compromised in known data breaches.

## Quick Start with Docker

### Prerequisites
- Docker and Docker Compose installed

### Running the Application locally

1. Clone or download this repository
2. Navigate to the src directory:
   ```bash
   cd src
   ```

3. Build and start the application:
   ```bash
   docker-compose up --build
   ```

4. Access the application:
   - Web Interface: http://localhost/
   - API Health Check: http://localhost/api/health

### Stopping the Application

```bash
docker-compose down
```

## Usage

1. Open http://localhost/ in your browser
2. Enter an email address to check
3. View the results to see if the email has been found in known data breaches

## Test Emails

The following emails are pre-loaded for testing:
- `test@example.com` (compromised)
- `user@domain.com` (compromised)  
- `safe@example.com` (not compromised)