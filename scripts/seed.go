package main

import (
	"log"
	"time"

	"github.com/izzy-Ti/ZemlyGo/configs"
	"github.com/izzy-Ti/ZemlyGo/internals/bootstrap"
	"github.com/izzy-Ti/ZemlyGo/internals/constants"
	"github.com/izzy-Ti/ZemlyGo/internals/domain"
	"github.com/izzy-Ti/ZemlyGo/internals/infrastructure/database"
)

func main() {
	log.Println("Seeding ZemlyGo development database...")

	cfg := configs.Load()
	db := bootstrap.NewDatabase(cfg)
	if err := database.RunMigration(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// 1. Seed Riders
	riderUser := domain.Users{
		SupabaseUserID: "sb-rider-001",
		Name:           "Alice Rider",
		Email:          "alice@example.com",
		Phone:          "+1555123456",
		Role:           constants.RoleRider,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	db.Where(domain.Users{Email: riderUser.Email}).FirstOrCreate(&riderUser)
	log.Printf("Seeded Rider: ID=%d, Name=%s\n", riderUser.ID, riderUser.Name)

	// 2. Seed Driver 1
	driverUser1 := domain.Users{
		SupabaseUserID: "sb-driver-001",
		Name:           "Bob Driver",
		Email:          "bob@example.com",
		Phone:          "+1555654321",
		Role:           constants.RoleDriver,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	db.Where(domain.Users{Email: driverUser1.Email}).FirstOrCreate(&driverUser1)

	driver1 := domain.Drivers{
		UserId:     driverUser1.ID,
		LicenseNo:  "DL-987654321",
		IsApproved: true,
		IsOnline:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	db.Where(domain.Drivers{UserId: driverUser1.ID}).FirstOrCreate(&driver1)

	vehicle1 := domain.Vehicle{
		DriverID:    driver1.ID,
		PlateNumber: "UBER-101",
		Model:       "Toyota Camry Hybrid",
		Color:       "Silver",
		Status:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	db.Where(domain.Vehicle{PlateNumber: vehicle1.PlateNumber}).FirstOrCreate(&vehicle1)

	loc1 := domain.DriverLocation{
		DriverID:  driver1.ID,
		Lat:       40.7128,
		Lng:       -74.0060,
		UpdatedAt: time.Now(),
	}
	db.Where(domain.DriverLocation{DriverID: driver1.ID}).FirstOrCreate(&loc1)
	log.Printf("Seeded Driver: ID=%d, Name=%s, Vehicle=%s\n", driver1.ID, driverUser1.Name, vehicle1.Model)

	// 3. Seed Admin
	adminUser := domain.Users{
		SupabaseUserID: "sb-admin-001",
		Name:           "System Admin",
		Email:          "admin@zemlygo.com",
		Phone:          "+1555999999",
		Role:           constants.RoleAdmin,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	db.Where(domain.Users{Email: adminUser.Email}).FirstOrCreate(&adminUser)
	log.Printf("Seeded Admin: ID=%d, Name=%s\n", adminUser.ID, adminUser.Name)

	log.Println("Database seed completed successfully!")
}
