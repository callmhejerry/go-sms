# School SMS – Multi-Tenant School Management System

Production-oriented school management backend built with:
- Go (`net/http`)
- PostgreSQL + sqlc
- Docker

## Features

- Multi-tenant architecture (shared schema)
- Authentication (JWT) + RBAC (owner, admin, teacher, accountant...)
- Academic management (sessions, classes, arms, students, admissions)
- Billing & payments
- Grading & report cards
- Inventory & stock movements
- Audit logging
- Security headers + rate limiting

## Quick Start (Development)

### Prerequisites
- Go 1.22+
- Docker
- golang-migrate
- sqlc

### Setup

```bash
cp .env.example .env
make up
make migrate-up
go run ./cmd/api