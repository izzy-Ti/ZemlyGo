# ZemlyGo — Project Progress

> Uber-style ride-hailing backend (Go + Gin + PostgreSQL + Supabase + WebSockets)  
> Last reviewed: July 13, 2026

---

## Overall status

| Area | Status | Notes |
|------|--------|-------|
| Architecture plan | Done | See `docs/structure.md` |
| App bootstrap & server | Partial | Server starts; only 2 routes wired |
| Domain models | Done | 7 entities defined |
| Database / migrations | Partial | GORM AutoMigrate; incomplete & missing tables |
| Repositories | Partial | Interfaces + Postgres stubs; many bugs & gaps |
| Services | Minimal | Only ride notification helper exists |
| HTTP handlers (REST) | Not started | No auth/user/driver/ride APIs |
| WebSocket / realtime | Partial | Hub + `/ws` endpoint working at basic level |
| Auth (Supabase) | Partial | Middleware written; not wired; config broken |
| Payments & ratings | Not started | Models + repos only |
| DevOps (Docker, env, scripts) | Not started | Planned in structure doc but not in repo |

**Rough completion: ~20–25%** — good skeleton and data model, but most business logic and API surface are still missing. The repo does **not fully compile** yet due to incomplete repository methods and an empty vehicle repository interface.

---

## What is done

### 1. Project foundation

- [x] Go module (`go.mod`) with Gin, GORM, PostgreSQL driver, Gorilla WebSocket, Supabase client, godotenv
- [x] Entry point: `cmd/main.go` — loads config, connects DB, runs migration, starts HTTP server
- [x] Config loader: `configs/config.go` reads `PORT` and `DATABASE_URL` from `.env`
- [x] Bootstrap layer: `App`, `Database`, `Router`, `Providers` packages exist

### 2. Domain layer (data models)

All core entities are defined under `internals/domain/`:

| Model | Purpose |
|-------|---------|
| `Users` | Rider/driver accounts linked to Supabase |
| `Drivers` | Driver profile, approval, online status |
| `Vehicle` | Driver vehicles (plate, model, color) |
| `Ride` | Pickup/destination, status, fare, ride type |
| `Payment` | Ride payment record |
| `Rating` | Post-ride ratings |
| `DriverLocation` | Live driver coordinates |

### 3. Repository layer (partial)

**Interfaces defined** for: User, Driver, Ride, Payment, Rating  
**Postgres implementations exist** for all of the above + Vehicle (but Vehicle interface file is empty).

Implemented concepts include:
- CRUD-style methods for users, drivers, rides, payments, ratings, vehicles
- Driver online toggle
- Driver location update
- Nearby driver search (Haversine via `internals/utils/distance.go`)
- Ride status updates, assign driver, cancel, complete
- Ride history by rider/driver

### 4. Realtime / WebSocket

- [x] `internals/realtime/` — Hub, Client, Event types
- [x] `GET /ws?user_id=&role=` — WebSocket upgrade, register client, read/write pumps
- [x] `RideService.NotifyDriverNewRide()` — calculates fare and pushes `ride.requested` event to a driver

### 5. Middleware (written, not used)

- [x] `AuthMiddleware` in `internals/middleware/supabase_user.go` — validates Supabase JWT from `Authorization: Bearer` header

### 6. Utilities

- [x] `Haversine()` distance helper in `internals/utils/distance.go`

### 7. Working HTTP routes (today)

| Method | Path | Status |
|--------|------|--------|
| GET | `/ping` | Works — health check |
| GET | `/ws` | Works — WebSocket connect (requires `user_id` query param) |

---

## What is NOT done (vs planned architecture)

These are listed in `docs/structure.md` but **do not exist in the codebase yet**:

- DTOs (`internals/dto/`)
- Constants (`roles`, `ride_status`, `errors`)
- REST handlers: auth, user, driver, vehicle, ride, payment, rating
- Services: auth, user, driver, vehicle, pricing, matching, location, payment, notification, rating
- Centralized routes file (`internals/routes/routes.go`)
- Extra middleware: role check, logger, recovery, rate limit
- Supabase infrastructure package (`client`, `auth`, `storage`)
- SQL migration files (`migrations/*.sql`)
- Docker / docker-compose
- `.env.example`, seed script, API docs
- Tests

---

## Known issues to fix first

Before building new features, these block compilation or cause runtime bugs:

1. **Build errors**
   - `user_repository.go`: `Create()` and `UpdateRole()` have no function body
   - `vehicle_repository.go` (interface): file is empty — Postgres impl references undefined interface

2. **Config bug**
   - `configs/config.go`: Supabase env vars use `os.Getenv("")` — always empty

3. **Incomplete migrations**
   - `RunMigration()` migrates Payment, Rating, Ride, Vehicle only
   - Missing: `Users`, `Drivers`, `DriverLocation`

4. **GORM usage bugs** (examples)
   - `First(ride)` instead of `First(&ride)` in several repos
   - `Save(ride)` without `.Error` check in ride repo `Update()`
   - Wrong column names: `supabase_id` vs model field `SupabaseUserID`
   - Rating repo queries use `from_user_iD` / `to_user_iD` (typos)
   - Vehicle `Delete()` deletes `DriverLocation` instead of `Vehicle`

5. **Ride status inconsistency**
   - Domain default: `requested`
   - Repo active queries use: `ACTIVE`, `CANCLED`, `COMPLETE` (typos / mismatch)

6. **Security**
   - WebSocket `CheckOrigin` allows all origins
   - Auth middleware not applied to any route
   - WebSocket auth uses plain `user_id` query param (no token verification)

7. **Naming / structure drift**
   - Plan says `internal/`; code uses `internals/`
   - `Vehcle.go` typo; `ride_respository.go` typo

---

## Recommended build order

Follow this sequence so each phase unlocks the next. Do not skip Phase 0.

### Phase 0 — Make it compile & runnable (1–2 days)

1. Fix config (`SUPABASE_URL`, `SUPABASE_ANON_KEY`, `SUPABASE_SERVICE_KEY`)
2. Complete `VehicleRepository` interface
3. Finish stub methods in `UserRepository`
4. Fix GORM pointer/`.Error` bugs across all repos
5. Add all domain models to `RunMigration()`
6. Add `.env.example` with required variables
7. Verify: `go build ./...` and server starts with Postgres

### Phase 1 — Auth & user identity (2–3 days)

Uber flow starts with knowing who is calling.

1. Supabase client bootstrap in `internals/infrastructure/supabase/`
2. Wire `AuthMiddleware` on protected route groups
3. **User service** — sync Supabase user → local `Users` row on first login
4. **User handler** — `GET /me`, `PATCH /me`, role assignment
5. **Constants** — roles (`rider`, `driver`, `admin`), standard error responses
6. **DTOs** — auth + user request/response shapes

**Deliverable:** Rider can sign in via Supabase token and get/create their profile.

### Phase 2 — Driver & vehicle onboarding (2–3 days)

Drivers cannot take rides without approval and a vehicle.

1. **Driver service + handler**
   - Register as driver (license number)
   - Admin approve driver (`IsApproved`)
   - Go online / go offline
2. **Vehicle service + handler**
   - Add vehicle, list vehicles, set active vehicle
3. **Role middleware** — restrict driver routes to `driver` role
4. Fix driver location: upsert if not exists (currently fails on first update)

**Deliverable:** Approved driver can go online with an active vehicle.

### Phase 3 — Location & matching (2–3 days)

Core Uber mechanic: find nearby available drivers.

1. **Location service**
   - `POST /drivers/location` — driver sends lat/lng (REST or WS)
   - Persist to `DriverLocation`
2. **Matching service**
   - Given pickup lat/lng + ride type → find nearest online, approved, available drivers
   - Filter drivers without active ride
3. **Pricing service**
   - Move fare logic out of `RideService` into dedicated service
   - Support ride types: `standard`, `van`, `family`
4. **Endpoint:** `POST /rides/estimate` — return fare + ETA before booking

**Deliverable:** Rider sees price estimate; system can find candidate drivers.

### Phase 4 — Ride lifecycle (3–5 days)

The main product loop.

| Step | Rider action | Driver action | Status flow |
|------|--------------|---------------|-------------|
| 1 | Request ride | — | `requested` |
| 2 | — | Accept ride | `accepted` |
| 3 | — | Arrived at pickup | `arrived` |
| 4 | — | Start trip | `in_progress` |
| 5 | — | End trip | `completed` |
| Either | Cancel | Cancel | `cancelled` |

1. Define ride status constants (single source of truth)
2. **Ride service** — full CRUD + state machine with validation (no invalid transitions)
3. **Ride handler** REST API:
   - `POST /rides` — request ride
   - `POST /rides/:id/accept` — driver accepts
   - `POST /rides/:id/arrive` / `start` / `complete` / `cancel`
   - `GET /rides/active` — current ride for rider or driver
   - `GET /rides/history` — paginated history
4. Integrate matching: on request → notify nearby drivers via WebSocket (`ride.requested`)
5. WebSocket events: `ride.accepted`, `ride.driver_location`, `ride.completed`, etc.

**Deliverable:** End-to-end ride from request to completion over REST + WebSocket.

### Phase 5 — Payments (2–3 days)

1. **Payment service** — create payment on ride complete, mark paid/failed
2. **Payment handler** — `POST /rides/:id/pay`, `GET /payments/:id`
3. Stub or integrate real gateway (Stripe / local mobile money) later

**Deliverable:** Payment record tied to completed ride.

### Phase 6 — Ratings & notifications (1–2 days)

1. **Rating service + handler** — rate after completed ride (rider ↔ driver)
2. **Notification service** — push WS events + optional email/SMS hooks
3. Average rating on driver profile

**Deliverable:** Post-ride feedback loop.

### Phase 7 — Production hardening (ongoing)

1. Logger, recovery, rate-limit middleware
2. Docker + docker-compose (app + Postgres)
3. SQL migrations (replace or supplement GORM AutoMigrate)
4. API documentation (`docs/api.md`)
5. Seed script for dev data
6. Unit + integration tests for ride state machine and matching
7. Secure WebSocket (JWT on connect, origin check)
8. Admin endpoints (approve drivers, view rides)

---

## Suggested priority summary

```
Phase 0  Fix compile + DB          ← START HERE
Phase 1  Auth & users
Phase 2  Driver & vehicle onboarding
Phase 3  Location & matching & pricing
Phase 4  Ride lifecycle (core product)
Phase 5  Payments
Phase 6  Ratings & notifications
Phase 7  DevOps & hardening
```

---

## Quick reference: file checklist

| Planned (structure.md) | Exists? | Compiles? | Wired to server? |
|------------------------|---------|-----------|------------------|
| `cmd/main.go` | Yes | Partial | Yes |
| Domain models | Yes | Yes | — |
| Repository interfaces | Mostly | No (vehicle empty) | No |
| Repository postgres | Yes | No (bugs/stubs) | No |
| Services | 1 of ~11 | Partial | No |
| Handlers | 1 of ~8 | Yes (WS only) | WS only |
| Middleware auth | Yes | Yes | No |
| Realtime hub | Yes | Yes | Yes |
| Routes | Inline in router | Yes | 2 routes |
| DTOs | No | — | — |
| Constants | No | — | — |
| Supabase infra | No | — | — |
| Migrations SQL | No | — | — |
| Docker | No | — | — |

---

## Next immediate tasks (this week)

If you are picking up work now, do these in order:

1. Fix Phase 0 compile blockers (config, vehicle interface, user repo stubs, GORM fixes)
2. Migrate `Users`, `Drivers`, `DriverLocation` in `RunMigration()`
3. Add Supabase client + wire auth middleware
4. Implement `POST /users/sync` or `GET /me` as first real authenticated endpoint
5. Standardize ride statuses in `internals/constants/ride_status.go`

Once Phase 0–1 are done, you have a solid base to build the actual Uber-like ride flow in Phases 2–4.
