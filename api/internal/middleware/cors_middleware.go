// middleware/cors.go
package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupCORS() gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	
	// Per a desenvolupament amb React Native, permetre tots els orígens
	// En producció, hauries de restringir això a dominis específics
	corsConfig.AllowAllOrigins = true
	
	// Si vols més control, descomenta i modifica això:
	// corsConfig.AllowOrigins = []string{
	//     "http://localhost:5173",           // Web local
	//     "https://orkestra.zenith.ovh",     // Web producció
	//     "exp://192.168.1.*",               // Expo local network
	// }
	
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept"}
	corsConfig.ExposeHeaders = []string{"Content-Length", "Content-Type"}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 12 * time.Hour
	
	return cors.New(corsConfig)
}