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

3. Set up environment configuration:
   ```bash
   cp .env.example .env
   ```
   Edit the `.env` file to set your admin credentials (required for admin endpoints).

4. Build and start the application:
   ```bash
   docker-compose up --build
   ```
   Note: Initial startup may take 1-2 minutes as the system builds 3 backend instances and waits for health checks.

5. Access the application:
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

## Admin Interface Testing

Admin endpoints require authentication using credentials from your `.env` file:

```bash
# Test admin status (replace credentials with your .env values)
curl -u admin:your-admin-password http://localhost/admin/status

# Test admin metrics
curl -u admin:your-admin-password http://localhost/admin/metrics
```

## Troubleshooting

**Port 80 already in use**: If you get port conflicts, stop any web servers using port 80 or modify the docker-compose.yml port mapping to `"8080:80"` and access via http://localhost:8080/

For detailed testing and troubleshooting, see [TESTING_GUIDE.md](TESTING_GUIDE.md)