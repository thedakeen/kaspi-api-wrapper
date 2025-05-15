# Kaspi Pay API Wrapper

This service acts as a middleware between your application and Kaspi Pay's payment processing system. It implements three integration schemes described in the Kaspi Pay API documentation:

- **Basic Scheme** - Simplified integration with API key authentication
- **Standard Scheme** - Enhanced features with certificate-based authentication
- **Enhanced Scheme** - Full feature set with certificate authentication

## Getting Started

### Prerequisites

- Go 1.23+
- Docker and Docker Compose
- PostgreSQL 16+ (provided via Docker)
- WSL / Linux OS

### With Docker Compose

1. Clone the repository
2. Run the setup command:

```bash
make setup
```
- Creates .env file from example if it doesn't exist
- Builds Docker images
- Starts the database
- Runs migrations
- Creates database dump
- Starts all services

3. Servers:
   - HTTP API: http://localhost:8081
   - gRPC API: http://localhost:8082


```
# .env example
HTTP_PORT=8081
GRPC_PORT=8082

# Integration scheme (basic, standard, enhanced)
KASPI_API_SCHEME=basic

# API endpoints
KASPI_API_BASE_URL_BASIC=http://mock-kaspi-api:1080/r1/v01
KASPI_API_BASE_URL_STANDARD=http://mock-kaspi-api:1080/r2/v01
KASPI_API_BASE_URL_ENHANCED=http://mock-kaspi-api:1080/r3/v01

# For basic scheme
KASPI_API_KEY=test_api_key

# For standard and enhanced schemes
KASPI_PFX_FILE=./certs/client.pfx
KASPI_KEY_PASSWORD=test123
KASPI_ROOT_CA_FILE=./certs/ca.crt

# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=kaspi_pay
DB_SSL_MODE=disable
```

## API Reference

### REST API Endpoints

The service provides a RESTful API with endpoints that correspond to the Kaspi Pay API. The base URL is `http://localhost:8081/api`.

#### Basic scheme endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/tradepoints` | Get trade points |
| POST | `/device/register` | Register device |
| POST | `/device/delete` | Delete device |
| POST | `/qr/create` | Create QR code for payment |
| POST | `/qr/create-link` | Create payment link |
| GET | `/payment/status/{qrPaymentId}` | Get payment status |

#### Standard scheme endpoints (all Basic endpoints plus)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/return/create` | Create refund QR |
| GET | `/return/status/{qrReturnId}` | Get refund status |
| POST | `/return/operations` | Get customer operations |
| GET | `/payment/details` | Get payment details |
| POST | `/payment/return` | Refund payment |

#### Enhanced scheme endpoints (all Standard endpoints plus)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/tradepoints/enhanced/{organizationBin}` | Get trade points |
| POST | `/device/register/enhanced` | Register device |
| POST | `/device/delete/enhanced` | Delete device |
| POST | `/qr/create/enhanced` | Create QR code |
| POST | `/qr/create-link/enhanced` | Create payment link |
| POST | `/enhanced/payment/return` | Refund payment without customer |
| GET | `/remote/client-info` | Get client info |
| POST | `/remote/create` | Create remote payment |
| POST | `/remote/cancel` | Cancel remote payment |

#### Test endpoints (all schemes)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/test/health` | Health check |
| POST | `/test/payment/scan` | Test QR scan |
| POST | `/test/payment/confirm` | Test payment confirmation |
| POST | `/test/payment/scanerror` | Test QR scan error |
| POST | `/test/payment/confirmerror` | Test payment confirmation error |

### gRPC API

The service also provides a gRPC API on port 8082. The proto files are located in the `pkg/protos/proto` directory:

- `device/device.proto` - Device management operations
- `payment/payment.proto` - Payment processing operations
- `refund/refund.proto` - Refund operations (standard scheme)
- `refund_enhanced/refund_enhanced.proto` - Enhanced refund operations
- `utility/utility.proto` - Utility operations