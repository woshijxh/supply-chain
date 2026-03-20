# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Supply Chain Management System with two main components:
- `supply-chain/` - Frontend (Vue 3 + TypeScript + Vite)
- `supply-chain-server/` - Backend (Go + Gin + GORM + MySQL)

## Build Commands

### Frontend (supply-chain/)
```bash
cd supply-chain
npm install        # Install dependencies
npm run dev        # Start dev server (port 3000)
npm run build      # TypeScript check + production build
npm run preview    # Preview production build
npx vue-tsc --noEmit # Type-only check
```

### Backend (supply-chain-server/)
```bash
cd supply-chain-server
go mod tidy        # Install dependencies
go build -o bin/server ./cmd/server/main.go  # Compile
./bin/server       # Run (port 8080)
go run ./cmd/server/main.go  # Run without compiling
```

### Test Commands
```bash
cd supply-chain-server
go test ./...                    # Run all tests
go test ./internal/service -run TestName  # Run single test
go test -v ./...                 # Verbose output
```

### Database Setup
```bash
mysql -u root -p < supply-chain-server/docs/sql/schema.sql
mysql -u root -p < supply-chain-server/docs/sql/test_data.sql  # Optional test data
```

## Development Workflow

1. Start backend: `cd supply-chain-server && go run ./cmd/server/main.go`
2. Start frontend: `cd supply-chain && npm run dev`
3. Access app: http://localhost:3000
4. API base: http://localhost:8080/api

## Architecture

### Frontend Architecture

**Tech Stack:** Vue 3.4 + TypeScript 5.3 + Vite 5 + Pinia + Vue Router 4 + PrimeVue 4

**Key Files:**
- `src/main.ts` - App initialization, PrimeVue config, component registration
- `src/router/index.ts` - Route definitions with auth guards
- `src/stores/supply.ts` - Central Pinia store for all data
- `src/api/request.ts` - Axios config with JWT interceptor
- `src/api/index.ts` - API endpoint definitions

**Component Pattern:**
```vue
<script setup lang="ts">
// Order: imports → types → composables → reactive state → computed → methods → lifecycle
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>
```

**State Management:**
- Single Pinia store (`supply.ts`) holds all entities
- Components use store methods to update state
- Computed getters for filtered lists (e.g., `activeSuppliers`, `lowStockItems`)

**Authentication:**
- JWT stored in localStorage
- 401 responses redirect to /login
- Route meta `requiresAuth` guards protected pages

**Styling:**
- SCSS with scoped styles per component
- PrimeVue custom theme in `src/styles/primevue-theme.ts`
- Dark mode toggle via `.dark-mode` class
- Remixicon + Primeicons for icons

### Backend Architecture

**Tech Stack:** Go 1.24 + Gin + GORM + MySQL + JWT + Viper

**Layer Structure:**
```
cmd/server/main.go              → Entry point
internal/router/router.go       → Dependency injection, route setup
internal/handler/               → HTTP handlers (gin.Context)
internal/service/               → Business logic
internal/repository/            → Data access (GORM)
internal/model/                 → GORM models
pkg/response/                   → Unified JSON responses
pkg/database/                   → DB connection + auto-migrate
```

**Handler Pattern:**
```go
func (h *Handler) List(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    items, total, err := h.service.List(page, pageSize, keyword)
    if err != nil {
        response.ServerError(c, err.Error())
        return
    }
    response.PageSuccess(c, items, total, page, pageSize)
}
```

**Response Format:**
- Success: `{code: 0, message: "success", data: {...}}`
- Page: `{code: 0, message: "success", data: {list: [], total, page, pageSize}}`
- Error: `{code: 400|401|404|500, message: "..."}`

**Authentication:**
- JWT middleware on protected routes
- Claims: UserID, Username, Role
- Token from `Authorization: Bearer <token>` header

**Configuration:**
- Viper loads from `configs/config.yaml`
- Env: `SERVER_PORT`, `DATABASE_HOST`, etc.
- Defaults in `internal/config/config.go`

## Key Conventions

### Frontend
- **Pages:** PascalCase + `Page` suffix (`SupplierPage.vue`)
- **Imports:** Use `@/` alias for src directory
- **Types:** Interfaces in `src/types/index.ts`, union types for status literals
- **i18n:** `zh-CN` default, keys in `src/i18n/locales/`

### Backend
- **Packages:** lowercase (repository, service, handler)
- **Types:** PascalCase (SupplierService, UserRepository)
- **Methods:** PascalCase exported, camelCase private
- **Errors:** Wrap with `fmt.Errorf("context: %w", err)`

## API Endpoints

```
POST   /api/auth/login
POST   /api/auth/register
GET    /api/auth/profile

GET    /api/suppliers?page=&pageSize=&keyword=
POST   /api/suppliers
GET    /api/suppliers/:id
PUT    /api/suppliers/:id
DELETE /api/suppliers/:id

GET    /api/inventory
POST   /api/inventory/stock-in
POST   /api/inventory/stock-out
GET    /api/inventory/stats

GET    /api/procurement
POST   /api/procurement
PUT    /api/procurement/:id/status

GET    /api/sales
POST   /api/sales
PUT    /api/sales/:id/status

GET    /api/logistics
POST   /api/logistics
PUT    /api/logistics/:id/status

GET    /api/dashboard/stats
```

## Database Models

Auto-migrated on startup via `pkg/database/database.go`:
- User, Supplier, Product, Inventory
- ProcurementOrder, ProcurementItem
- SalesOrder, SalesOrderItem
- LogisticsOrder, LogisticsTimeline

See `docs/sql/schema.sql` for full schema.

## Reference Project

`nginxpulse/` is a reference project (Nginx log analyzer in Go+Vue) - do not modify unless specifically requested.
