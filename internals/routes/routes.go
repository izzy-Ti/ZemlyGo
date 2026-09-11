package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/izzy-Ti/ZemlyGo/configs"
	"github.com/izzy-Ti/ZemlyGo/internals/constants"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
	"github.com/izzy-Ti/ZemlyGo/internals/middleware"
)

type Providers struct {
	WebsocketHandler middlewareHandler
	AuthHandler      authHandler
	UserHandler      userHandler
	DriverHandler    driverHandler
	VehicleHandler   vehicleHandler
	RideHandler      rideHandler
	PaymentHandler   paymentHandler
	RatingHandler    ratingHandler
	MapHandler       mapHandler
	ChatHandler      chatHandler
	SafetyHandler    safetyHandler
	MetricsHandler   metricsHandler
	AuthMiddleware   gin.HandlerFunc
}

type middlewareHandler interface {
	Handle(c *gin.Context)
}

type authHandler interface {
	Sync(c *gin.Context)
	GetMe(c *gin.Context)
	DevToken(c *gin.Context)
}

type userHandler interface {
	GetMe(c *gin.Context)
	UpdateMe(c *gin.Context)
	List(c *gin.Context)
}

type driverHandler interface {
	Register(c *gin.Context)
	GetMe(c *gin.Context)
	SetOnline(c *gin.Context)
	SetOffline(c *gin.Context)
	UpdateLocation(c *gin.Context)
	GetEarnings(c *gin.Context)
	Approve(c *gin.Context)
	GetNearby(c *gin.Context)
}

type vehicleHandler interface {
	Create(c *gin.Context)
	List(c *gin.Context)
	SetActive(c *gin.Context)
}

type rideHandler interface {
	Estimate(c *gin.Context)
	RequestRide(c *gin.Context)
	AcceptRide(c *gin.Context)
	Arrive(c *gin.Context)
	Start(c *gin.Context)
	Complete(c *gin.Context)
	Cancel(c *gin.Context)
	GetActive(c *gin.Context)
	GetByID(c *gin.Context)
	GetHistory(c *gin.Context)
}

type paymentHandler interface {
	Pay(c *gin.Context)
	AddTip(c *gin.Context)
	GetByID(c *gin.Context)
}

type ratingHandler interface {
	Rate(c *gin.Context)
	GetDriverRatings(c *gin.Context)
}

type mapHandler interface {
	GetRoute(c *gin.Context)
	ReverseGeocode(c *gin.Context)
	Geocode(c *gin.Context)
	SearchPlaces(c *gin.Context)
}

type chatHandler interface {
	SendMessage(c *gin.Context)
	GetMessages(c *gin.Context)
}

type safetyHandler interface {
	TriggerEmergency(c *gin.Context)
	GetAlerts(c *gin.Context)
	ResolveAlert(c *gin.Context)
}

type metricsHandler interface {
	GetStats(c *gin.Context)
	GetSurgeHeatmap(c *gin.Context)
}

func SetupRoutes(r *gin.Engine, p *Providers, cfg *configs.Config) {
	// Global Middlewares
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.RecoveryMiddleware())

	// Rate limiter (100 tokens capacity, 20 tokens refill per second)
	limiter := middleware.NewRateLimiter(20.0, 100.0)
	r.Use(middleware.RateLimitMiddleware(limiter))

	// System Health & Metrics
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "PONG"})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, dto.SuccessResponse(gin.H{
			"status": "healthy",
			"env":    cfg.Env,
			"auth":   "neon-auth",
		}, "ZemlyGo API is online"))
	})

	r.GET("/metrics", p.MetricsHandler.GetStats)

	// Realtime WebSocket
	r.GET("/ws", p.WebsocketHandler.Handle)

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Public Auth routes (Neon Auth sync & Dev tokens)
		auth := v1.Group("/auth")
		{
			auth.POST("/sync", p.AuthHandler.Sync)
			auth.POST("/dev-token", p.AuthHandler.DevToken)
		}

		// Public Maps & Pricing routes
		mapsGroup := v1.Group("/maps")
		{
			mapsGroup.POST("/route", p.MapHandler.GetRoute)
			mapsGroup.GET("/reverse", p.MapHandler.ReverseGeocode)
			mapsGroup.GET("/geocode", p.MapHandler.Geocode)
			mapsGroup.GET("/places", p.MapHandler.SearchPlaces)
		}

		pricing := v1.Group("/pricing")
		{
			pricing.GET("/heatmap", p.MetricsHandler.GetSurgeHeatmap)
		}

		// Authenticated routes
		authenticated := v1.Group("")
		authenticated.Use(p.AuthMiddleware)
		{
			// Current user profile
			authenticated.GET("/auth/me", p.AuthHandler.GetMe)
			authenticated.GET("/users/me", p.UserHandler.GetMe)
			authenticated.PATCH("/users/me", p.UserHandler.UpdateMe)

			// Rides (Rider & Shared actions)
			rides := authenticated.Group("/rides")
			{
				rides.POST("/estimate", p.RideHandler.Estimate)
				rides.POST("", p.RideHandler.RequestRide)
				rides.GET("/active", p.RideHandler.GetActive)
				rides.GET("/history", p.RideHandler.GetHistory)
				rides.GET("/:id", p.RideHandler.GetByID)
				rides.POST("/:id/cancel", p.RideHandler.Cancel)

				// In-App Trip Chat
				rides.POST("/:id/messages", p.ChatHandler.SendMessage)
				rides.GET("/:id/messages", p.ChatHandler.GetMessages)

				// Safety & SOS Emergency
				rides.POST("/:id/emergency", p.SafetyHandler.TriggerEmergency)
				rides.GET("/:id/emergency", p.SafetyHandler.GetAlerts)

				// Payments & Tips
				rides.POST("/:id/pay", p.PaymentHandler.Pay)
				rides.POST("/:id/tip", p.PaymentHandler.AddTip)

				// Ratings
				rides.POST("/:id/rate", p.RatingHandler.Rate)
			}

			// Payments detail
			authenticated.GET("/payments/:id", p.PaymentHandler.GetByID)

			// Driver discovery & profile
			authenticated.GET("/drivers/:id/ratings", p.RatingHandler.GetDriverRatings)
			authenticated.POST("/drivers/register", p.DriverHandler.Register)
			authenticated.GET("/drivers/nearby", p.DriverHandler.GetNearby)

			// Driver-specific operations (Guarded by Driver Role)
			driverOnly := authenticated.Group("")
			driverOnly.Use(middleware.RequireRole(constants.RoleDriver))
			{
				driverOnly.GET("/drivers/me", p.DriverHandler.GetMe)
				driverOnly.POST("/drivers/online", p.DriverHandler.SetOnline)
				driverOnly.POST("/drivers/offline", p.DriverHandler.SetOffline)
				driverOnly.POST("/drivers/location", p.DriverHandler.UpdateLocation)
				driverOnly.GET("/drivers/earnings", p.DriverHandler.GetEarnings)

				// Vehicles
				driverOnly.POST("/vehicles", p.VehicleHandler.Create)
				driverOnly.GET("/vehicles", p.VehicleHandler.List)
				driverOnly.POST("/vehicles/:id/activate", p.VehicleHandler.SetActive)

				// Ride Execution Lifecycle
				driverOnly.POST("/rides/:id/accept", p.RideHandler.AcceptRide)
				driverOnly.POST("/rides/:id/arrive", p.RideHandler.Arrive)
				driverOnly.POST("/rides/:id/start", p.RideHandler.Start)
				driverOnly.POST("/rides/:id/complete", p.RideHandler.Complete)
			}

			// Admin-only operations
			adminOnly := authenticated.Group("")
			adminOnly.Use(middleware.RequireRole(constants.RoleAdmin))
			{
				adminOnly.POST("/drivers/:id/approve", p.DriverHandler.Approve)
				adminOnly.GET("/users", p.UserHandler.List)
				adminOnly.GET("/admin/stats", p.MetricsHandler.GetStats)
				adminOnly.POST("/admin/emergency/:id/resolve", p.SafetyHandler.ResolveAlert)
			}
		}
	}
}
