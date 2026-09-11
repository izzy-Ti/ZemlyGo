package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/izzy-Ti/ZemlyGo/internals/handler"
	"github.com/izzy-Ti/ZemlyGo/internals/middleware"
	"github.com/izzy-Ti/ZemlyGo/internals/repository/postgres"
	"github.com/izzy-Ti/ZemlyGo/internals/routes"
	"github.com/izzy-Ti/ZemlyGo/internals/service"
)

type Providers struct {
	WebsocketHandler *handler.WebsocketHandler
	AuthHandler      *handler.AuthHandler
	UserHandler      *handler.UserHandler
	DriverHandler    *handler.DriverHandler
	VehicleHandler   *handler.VehicleHandler
	RideHandler      *handler.RideHandler
	PaymentHandler   *handler.PaymentHandler
	RatingHandler    *handler.RatingHandler
	MapHandler       *handler.MapHandler
	ChatHandler      *handler.ChatHandler
	SafetyHandler    *handler.SafetyHandler
	MetricsHandler   *handler.MetricsHandler
	AuthMiddleware   gin.HandlerFunc
}

func NewProvider(app *App) *Providers {
	// Repositories
	userRepo := postgres.NewUserRepository(app.DB)
	driverRepo := postgres.NewDriverRepository(app.DB)
	vehicleRepo := postgres.NewVehicleRepo(app.DB)
	rideRepo := postgres.NewRideRepository(app.DB)
	paymentRepo := postgres.NewPaymentRepo(app.DB)
	ratingRepo := postgres.NewRatingRepo(app.DB)
	chatRepo := postgres.NewChatRepository(app.DB)
	emergencyRepo := postgres.NewEmergencyRepository(app.DB)

	// Services
	mapService := service.NewMapService(app.Config.MapsProvider, app.Config.MapsAPIKey)
	pricingService := service.NewPricingService()
	matchingService := service.NewMatchingService(driverRepo, rideRepo, vehicleRepo, app.Hub)
	locationService := service.NewLocationService(driverRepo, rideRepo, app.Hub)
	rideService := service.NewRideService(
		rideRepo,
		driverRepo,
		userRepo,
		vehicleRepo,
		paymentRepo,
		pricingService,
		matchingService,
		app.Hub,
	)
	authService := service.NewAuthService(userRepo, app.Config.JWTSecret)
	userService := service.NewUserService(userRepo)
	driverService := service.NewDriverService(driverRepo, userRepo, vehicleRepo, rideRepo, paymentRepo)
	vehicleService := service.NewVehicleService(vehicleRepo, driverRepo)
	paymentService := service.NewPaymentService(paymentRepo, rideRepo, app.Hub)
	ratingService := service.NewRatingService(ratingRepo, rideRepo, driverRepo, app.Hub)
	chatService := service.NewChatService(chatRepo, rideRepo, driverRepo, app.Hub)
	safetyService := service.NewSafetyService(emergencyRepo, rideRepo, app.Hub)

	// Neon Auth Middleware
	authMiddleware := middleware.NeonAuthMiddleware(app.NeonAuth, userRepo)

	// Handlers
	wsHandler := handler.NewWebSocketHandler(app.Hub, app.Config.JWTSecret)
	authHandler := handler.NewAuthHandler(authService, userService)
	userHandler := handler.NewUserHandler(userService)
	driverHandler := handler.NewDriverHandler(driverService, locationService)
	vehicleHandler := handler.NewVehicleHandler(vehicleService, driverService)
	rideHandler := handler.NewRideHandler(rideService, pricingService, driverService)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	ratingHandler := handler.NewRatingHandler(ratingService)
	mapHandler := handler.NewMapHandler(mapService)
	chatHandler := handler.NewChatHandler(chatService)
	safetyHandler := handler.NewSafetyHandler(safetyService)
	metricsHandler := handler.NewMetricsHandler(app.DB, pricingService)

	return &Providers{
		WebsocketHandler: wsHandler,
		AuthHandler:      authHandler,
		UserHandler:      userHandler,
		DriverHandler:    driverHandler,
		VehicleHandler:   vehicleHandler,
		RideHandler:      rideHandler,
		PaymentHandler:   paymentHandler,
		RatingHandler:    ratingHandler,
		MapHandler:       mapHandler,
		ChatHandler:      chatHandler,
		SafetyHandler:    safetyHandler,
		MetricsHandler:   metricsHandler,
		AuthMiddleware:   authMiddleware,
	}
}

func (p *Providers) ToRouteProviders() *routes.Providers {
	return &routes.Providers{
		WebsocketHandler: p.WebsocketHandler,
		AuthHandler:      p.AuthHandler,
		UserHandler:      p.UserHandler,
		DriverHandler:    p.DriverHandler,
		VehicleHandler:   p.VehicleHandler,
		RideHandler:      p.RideHandler,
		PaymentHandler:   p.PaymentHandler,
		RatingHandler:    p.RatingHandler,
		MapHandler:       p.MapHandler,
		ChatHandler:      p.ChatHandler,
		SafetyHandler:    p.SafetyHandler,
		MetricsHandler:   p.MetricsHandler,
		AuthMiddleware:   p.AuthMiddleware,
	}
}
