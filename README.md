# Supply Chain Management System

供应链管理系统 - A full-stack supply chain management application.

## Tech Stack

**Frontend:** Vue 3 + TypeScript + Vite + Pinia + PrimeVue

**Backend:** Go + Gin + GORM + MySQL + JWT

## Project Structure

```
supply-chain/           # Frontend (Vue 3)
supply-chain-server/    # Backend (Go)
```

## Quick Start

### Prerequisites

- Node.js 18+
- Go 1.21+
- MySQL 5.7+

### Backend Setup

```bash
cd supply-chain-server

# Install dependencies
go mod tidy

# Copy config and modify (use configs/config.example.yaml as reference)
cp configs/config.example.yaml configs/config.yaml
# Edit configs/config.yaml with your database credentials

# Run database migrations (create supply_chain database first)
mysql -u root -p < docs/sql/schema.sql

# Run server
go run ./cmd/server/main.go
# Server runs on http://localhost:8080
```

### Frontend Setup

```bash
cd supply-chain

# Install dependencies
npm install

# Run dev server
npm run dev
# App runs on http://localhost:3000
```

### Default Login

- Username: admin
- Password: 123456

## API Endpoints

See `CLAUDE.md` for detailed API documentation.

## Features

- 供应商管理 (Supplier Management)
- 采购管理 (Procurement)
- 库存管理 (Inventory)
- 销售订单 (Sales Orders)
- 物流跟踪 (Logistics)
- 用户权限管理 (User & RBAC)
- 仪表盘 (Dashboard)

## License

MIT
