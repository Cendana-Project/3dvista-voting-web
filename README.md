# VoteWeb - Simple & Secure Voting System

A production-ready, simple voting web application built with Go, featuring anti-abuse mechanisms and modern security practices.

## Features

- 🗳️ **Simple Voting**: One-click voting with confirmation
- 🔒 **Anti-Abuse**: IP-based rate limiting (1 vote per IP per innovation)
- 🛡️ **Security**: CSRF protection, security headers, HMAC-SHA256 IP hashing
- 🚀 **Production-Ready**: Docker support, graceful shutdown, health checks
- 📱 **Responsive Design**: Clean, modern UI that works on all devices
- 🎯 **No Authentication**: Frictionless voting experience

## Tech Stack

- **Backend**: Go 1.22 + Gin framework
- **Database**: PostgreSQL with pgx/v5 driver
- **Frontend**: HTML templates + Vanilla JavaScript
- **Deployment**: Docker + Docker Compose

## Project Structure

```
voteweb/
├── cmd/server/              # Application entry point
├── internal/
│   ├── app/                # Application initialization
│   ├── config/             # Configuration management
│   ├── domain/             # Business logic & entities
│   ├── http/               # HTTP handlers & middleware
│   │   ├── handlers/       # Request handlers
│   │   └── middleware/     # Middleware (CSRF, security, etc.)
│   ├── repo/               # Data access layer
│   └── util/               # Utility functions
├── web/
│   ├── static/             # CSS & JavaScript
│   └── templates/          # HTML templates
├── migrations/             # Database migrations
├── seed/                   # Seed data
└── docker-compose.yml      # Docker orchestration
```

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Or: Go 1.22+, PostgreSQL 14+

### Using Docker (Recommended)

1. **Clone the repository**
```bash
git clone <repository-url>
cd voteweb
```

2. **Start the application**
```bash
make docker-up
```

This will:
- Build the Docker images
- Start PostgreSQL database
- Run migrations
- Seed initial data
- Start the web application

3. **Access the application**

Open your browser and visit:
- http://localhost:8080/pemprov-jabar/jabar-digital-academy
- http://localhost:8080/pemprov-jabar/data-potensi-digital-desa-tapal-desa
- http://localhost:8080/bumn-bumd/simotip
- http://localhost:8080/kementrian-lembaga-pt/isopa-intelligent-solar-panel

4. **View logs**
```bash
make docker-logs
```

5. **Stop the application**
```bash
make docker-down
```

### Local Development

1. **Setup environment**
```bash
cp env.example .env
# Edit .env with your database credentials and settings
```

2. **Run migrations**
```bash
make migrate-up
```

3. **Seed database**
```bash
SEED=true go run ./cmd/server/main.go
```

4. **Run the application**
```bash
make run
```

## Configuration

Environment variables (`.env` file):

```env
# Database connection
DATABASE_URL=postgres://postgres:postgres@localhost:5432/voteweb?sslmode=disable

# Security: IP hashing salt (use a long random string in production)
IP_HASH_SALT=change-me-please-super-long-random-string-for-production

# Proxy settings (set to true if behind a reverse proxy)
TRUST_PROXY=false
ALLOWED_PROXY_CIDRS=10.0.0.0/8,172.16.0.0/12,192.168.0.0/16

# Application settings
APP_BASE_URL=http://localhost:8080
PORT=8080
GIN_MODE=debug
```

## Architecture

### Security Features

1. **IP-Based Anti-Abuse**
   - Each IP can vote only once per innovation
   - IPs are hashed with HMAC-SHA256 + salt before storage
   - Unique constraint ensures atomic vote deduplication

2. **CSRF Protection**
   - Double-submit cookie pattern
   - Token validation on all state-changing requests

3. **Security Headers**
   - Content Security Policy (CSP)
   - X-Content-Type-Options: nosniff
   - Referrer-Policy: no-referrer
   - X-Frame-Options: DENY
   - HSTS (when using HTTPS)

4. **Proxy-Aware IP Detection**
   - Configurable trusted proxy ranges
   - Proper X-Forwarded-For parsing

### Database Schema

**innovations** table:
- Stores innovation details (name, slug, division, etc.)
- Unique constraint on (group_slug, slug)

**votes** table:
- Stores vote records
- Unique constraint on (innovation_id, voter_ip_hash)
- Ensures one vote per IP per innovation atomically

### Vote Flow

1. User clicks "Vote" button
2. JavaScript confirmation dialog appears
3. On confirmation, POST request sent with CSRF token
4. Backend validates CSRF token and extracts client IP
5. IP is hashed with HMAC-SHA256
6. `INSERT ... ON CONFLICT DO NOTHING` ensures atomic deduplication
7. Vote count is retrieved and returned
8. Frontend displays success/already-voted modal

## API Endpoints

- `GET /:group/:slug` - Display innovation page
- `POST /api/vote/:group/:slug` - Submit vote
- `GET /healthz` - Health check endpoint

## Available Innovations

The application comes pre-seeded with 31 innovations across 6 categories:

1. **pemprov-jabar** (5 innovations)
2. **bumn-bumd** (4 innovations)
3. **kementrian-lembaga-pt** (5 innovations)
4. **smp-sma-sederajat** (5 innovations)
5. **pemda-kota** (5 innovations)
6. **pemda-kabupaten** (6 innovations)

## Development

### Running Tests

```bash
make test
```

### Building Binary

```bash
make build
```

### Database Operations

```bash
# Run migrations
make migrate-up

# Rollback migrations
make migrate-down

# Seed data
make seed
```

## Production Deployment

### Environment Variables for Production

1. **Change IP_HASH_SALT**: Use a cryptographically secure random string
```bash
openssl rand -base64 64
```

2. **Enable TRUST_PROXY**: If behind a reverse proxy (nginx, Cloudflare, etc.)
```env
TRUST_PROXY=true
```

3. **Set GIN_MODE to release**
```env
GIN_MODE=release
```

### Scaling Considerations

For ~1,000 concurrent users:
- Single instance is sufficient
- Database connection pool: 10-25 connections
- Consider adding Redis for rate limiting if needed
- Use a reverse proxy (nginx/Caddy) for TLS termination

### Monitoring

- Health check endpoint: `/healthz`
- Structured JSON logging
- Request ID tracking in headers and logs

## License

MIT

## Support

For issues or questions, please open an issue on the repository.


