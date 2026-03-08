mini-uber/
├── cmd/
│   └── server/
│       └── main.go
├── configs/
│   ├── config.go
│   └── config.yaml
├── internal/
│   ├── bootstrap/
│   │   ├── app.go
│   │   ├── database.go
│   │   ├── router.go
│   │   └── providers.go
│   ├── domain/
│   │   ├── user.go
│   │   ├── driver.go
│   │   ├── vehicle.go
│   │   ├── ride.go
│   │   ├── payment.go
│   │   └── rating.go
│   ├── dto/
│   │   ├── auth_dto.go
│   │   ├── user_dto.go
│   │   ├── driver_dto.go
│   │   ├── ride_dto.go
│   │   └── common_dto.go
│   ├── repository/
│   │   ├── interfaces/
│   │   │   ├── user_repository.go
│   │   │   ├── driver_repository.go
│   │   │   ├── vehicle_repository.go
│   │   │   └── ride_repository.go
│   │   └── postgres/
│   │       ├── user_repository.go
│   │       ├── driver_repository.go
│   │       ├── vehicle_repository.go
│   │       └── ride_repository.go
│   ├── service/
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── driver_service.go
│   │   ├── vehicle_service.go
│   │   ├── ride_service.go
│   │   ├── pricing_service.go
│   │   ├── matching_service.go
│   │   ├── location_service.go
│   │   ├── payment_service.go
│   │   ├── notification_service.go
│   │   └── rating_service.go
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── driver_handler.go
│   │   ├── vehicle_handler.go
│   │   ├── ride_handler.go
│   │   ├── payment_handler.go
│   │   └── rating_handler.go
│   ├── middleware/
│   │   ├── auth_middleware.go
│   │   ├── role_middleware.go
│   │   ├── logger_middleware.go
│   │   ├── recovery_middleware.go
│   │   └── rate_limit_middleware.go
│   ├── routes/
│   │   └── routes.go
│   ├── utils/
│   │   ├── jwt.go
│   │   ├── password.go
│   │   ├── response.go
│   │   ├── validator.go
│   │   ├── pagination.go
│   │   └── distance.go
│   ├── constants/
│   │   ├── roles.go
│   │   ├── ride_status.go
│   │   └── errors.go
│   └── infrastructure/
│       ├── supabase/
│       │   ├── client.go
│       │   ├── auth.go
│       │   └── storage.go
│       └── database/
│           └── migrations.go
├── pkg/
│   └── common/
│       └── helpers.go
├── migrations/
│   ├── 001_create_users.sql
│   ├── 002_create_drivers.sql
│   ├── 003_create_vehicles.sql
│   ├── 004_create_rides.sql
│   └── 005_create_ratings.sql
├── docs/
│   └── api.md
├── scripts/
│   ├── seed.go
│   └── dev.sh
├── Dockerfile
├── docker-compose.yml
├── .env
├── .env.example
├── go.mod
└── go.sum