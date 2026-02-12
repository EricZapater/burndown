package app

import (
	"api/config"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// App conté tots els recursos de l'aplicació
type App struct {
	Router *gin.Engine
	DB     *sql.DB
	Config *config.Config
	server *http.Server
}

// Run inicia el servidor HTTP amb graceful shutdown
func (a *App) Run() error {
	// Configurar servidor HTTP
	a.server = &http.Server{
		Addr:           ":" + a.Config.ApiPort,
		Handler:        a.Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Canal per capturar senyals del sistema
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Goroutine per executar el servidor
	go func() {
		log.Printf("🚀 Server starting on port %s", a.Config.ApiPort)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Esperar senyal de shutdown
	<-quit
	log.Println("🛑 Shutting down server...")

	// Context amb timeout per graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("✅ Server exited gracefully")
	return nil
}

// Close tanca tots els recursos de l'aplicació
func (a *App) Close() error {
	log.Println("🧹 Cleaning up resources...")
	
	if a.DB != nil {
		if err := a.DB.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
		log.Println("✅ Database connection closed")
	}
	
	return nil
}
