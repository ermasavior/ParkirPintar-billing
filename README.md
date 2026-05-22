# billing-service

[![Golang CI/CD](https://github.com/ermasavior/parkirpintar-billing/actions/workflows/cicd.yml/badge.svg)](https://github.com/ermasavior/parkirpintar-billing/actions/workflows/cicd.yml)

Calculates parking fees, creates invoices, and handles parking-payment outcomes via NATS.

## Responsibilities

- `CalculateAndCreateInvoice` — runs the pricing engine, creates an invoice (`PENDING_PAYMENT`), calls Payment Service for a parking-fee QRIS code. Idempotent via `idempotency_key`.
- `GetInvoice` — returns invoice state and fee breakdown
- `RetryPayment` — resets a `PAYMENT_FAILED` invoice to `PENDING_PAYMENT` and issues a new QRIS code
- Consumes `payment.parking.done` from NATS → transitions invoice to `PAID` (success) or `PAYMENT_FAILED` (failure/expired)

## Pricing

| Component | Rule |
|---|---|
| Parking fee | `ceil(duration_minutes / 60) × 5,000 IDR` |
| Overnight fee | `20,000 IDR × number of midnights crossed` |
| Booking fee | Always 5,000 IDR (stored on invoice, charged at reservation time) |

Pricing engine: [`pkg/pricing/engine.go`](pkg/pricing/engine.go)

## gRPC API

```
service BillingService {
  rpc CalculateAndCreateInvoice (CreateInvoiceRequest)  returns (CreateInvoiceResponse);
  rpc GetInvoice                (GetInvoiceRequest)     returns (GetInvoiceResponse);
  rpc RetryPayment              (RetryPaymentRequest)   returns (RetryPaymentResponse);
}
```

Proto: [`proto/billing/v1/billing.proto`](proto/billing/v1/billing.proto)

## Dependencies

| Dependency | Purpose |
|---|---|
| PostgreSQL | Invoices |
| NATS JetStream | Consume `payment.parking.done` |
| Payment Service (gRPC) | Create parking-fee QRIS payment |

## Configuration

```bash
cp .env.example .env
```

Key variables: `POSTGRES_DSN`, `NATS_URL`, `PAYMENT_SERVICE_URL`

## Development

```bash
make run              # run locally
make build            # compile binary → bin/billing
make test             # all tests
make test-unit        # unit tests only
make unit-test-coverage
make proto            # regenerate gRPC code from .proto
make mock             # regenerate mocks
```
