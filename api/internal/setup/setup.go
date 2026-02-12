package setup

import (
	"api/config"
	"api/internal/advisor"
	"api/internal/app"
	"api/internal/auth"
	"api/internal/database"
	"api/internal/meals"
	"api/internal/middleware"
	"api/internal/users"
	"fmt"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// NewApp inicialitza tota l'aplicació amb dependency injection
func NewApp(cfg *config.Config) (*app.App, error) {
	// 1. Inicialitzar infraestructura
	db, err := database.NewDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Verificar connexió a la base de dades
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 2. Inicialitzar JWT middleware
	jwtMiddleware, err := middleware.SetupJWT(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to setup JWT middleware: %w", err)
	}

	// 3. Inicialitzar serveis externs (Gemini AI)
	advisorService, err := advisor.NewService(cfg.LLMKey)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize advisor service: %w", err)
	}

	// 4. Inicialitzar repositories
	userRepo := users.NewRepository(db)
	mealRepo := meals.NewRepository(db)

	// 5. Inicialitzar services (injectar dependencies)
	userService := users.NewService(userRepo)
	authService := auth.NewService(userService, jwtMiddleware)
	mealService := meals.NewService(mealRepo)

	// 6. Inicialitzar handlers
	userHandler := users.NewHandler(userService)
	authHandler := auth.NewHandler(authService, jwtMiddleware)
	advisorHandler := advisor.NewHandler(advisorService)
	mealHandler := meals.NewHandler(mealService)

	// 7. Configurar router i routes
	router := setupRouter(cfg)
	setupRoutes(router, userHandler, authHandler, advisorHandler, mealHandler, jwtMiddleware)

	// Log de serveis inicialitzats
	fmt.Println("✅ Database connected")
	fmt.Println("✅ JWT middleware configured")
	if advisorService != nil {
		fmt.Println("✅ Gemini AI advisor initialized")
	}
	fmt.Println("✅ All services initialized")

	return &app.App{
		Router: router,
		DB:     db,
		Config: cfg,
	}, nil
}

// setupRouter configura el router Gin amb middleware globals
func setupRouter(cfg *config.Config) *gin.Engine {
	// Mode de producció o debug
	if cfg.ApiPort == "8125" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	
	// Middleware globals
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	// CORS middleware - IMPORTANT per a React Native
	router.Use(middleware.SetupCORS())

	return router
}

// setupRoutes registra totes les rutes de l'aplicació
func setupRoutes(
	router *gin.Engine,
	userHandler *users.Handler,
	authHandler *auth.Handler,
	advisorHandler *advisor.Handler,
	mealHandler *meals.Handler,
	jwtMiddleware *jwt.GinJWTMiddleware,
) {
	// Health check endpoint (públic)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "burndown-api",
		})
	})

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Public routes (auth)
		authGroup := v1.Group("/auth")
		auth.RegisterRoutes(authGroup, authHandler, jwtMiddleware)

		// Public routes (users - només registre)
		v1.POST("/users", userHandler.Create)

		// Protected routes (requereixen autenticació)
		protected := v1.Group("")
		protected.Use(jwtMiddleware.MiddlewareFunc())
		{
			// User routes
			users.RegisterRoutes(protected, userHandler)
			advisor.RegisterRoutes(protected, advisorHandler)
			meals.RegisterRoutes(protected, mealHandler)			
			// Aquí es poden afegir més mòduls:
			// meals.RegisterRoutes(protected, mealHandler)
			// workouts.RegisterRoutes(protected, workoutHandler)
		}
	}
}
