# Tax Calculator API

A Go API for calculating Canadian income tax using marginal tax brackets.

## Quick Start

### Prerequisites

- Go 1.25+
- Docker (for running the mock tax API)

### 1. Start the Mock Tax API

```bash
docker pull ptsdocker16/interview-test-server
docker run --init -p 5001:5001 -it ptsdocker16/interview-test-server
```

Follow instructions from [the README on Github](https://github.com/Points/interview-test-server?tab=readme-ov-file#get-up-and-running)

### 2. Run the Tax Calculator

```bash
# Clone and setup
git clone https://github.com/alaaeelsayed/tax-calculator
cd tax-calculator

# Install dependencies
make install

# Run tests
make test

# Start development server
make dev
```

## API Endpoints

### Calculate Tax

```bash
GET /taxes/{year}?salary={amount}

# Examples
curl "http://localhost:5002/taxes/2022?salary=50000"
curl "http://localhost:5002/taxes/2019?salary=100000"
```

**Response:**

```json
{
  "total_tax": 18141.11,
  "effective_rate": 0.18141105,
  "tax_by_bracket": [
    {
      "min": 0,
      "max": 47630,
      "rate": 0.15,
      "amount_taxable": 47630,
      "tax_payable": 7144.5
    },
    {
      "min": 47630,
      "max": 95259,
      "rate": 0.205,
      "amount_taxable": 47629,
      "tax_payable": 9763.945
    },
    {
      "min": 95259,
      "max": 147667,
      "rate": 0.26,
      "amount_taxable": 4741,
      "tax_payable": 1232.66
    }
  ]
}
```

## Configuration

Create a `.env` file:

```bash
cp .env.example .env
```

## Make Commands

```bash
make install       # Install dependencies
make test          # Run tests
make dev           # Run in development mode
make build         # Build binary
make clean         # Clean build artifacts
make lint          # Format and vet code
make help          # Show all commands
```

## Test Cases

**Running Tests:**

```bash
make test           # All tests
```

## Production Deployment

1. **Build:** `make build`
2. **Configure:** Set environment variables or use `.env`
3. **Run:** `./bin/tc-server`
