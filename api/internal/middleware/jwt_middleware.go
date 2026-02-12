// middleware/jwt.go
package middleware

import (
	"api/config"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func SetupJWT(cfg *config.Config) (*jwt.GinJWTMiddleware, error) {
    return jwt.New(&jwt.GinJWTMiddleware{
        Realm:       "api",
        Key:         []byte(cfg.JWTKey),
        Timeout:     time.Hour * 8,
        MaxRefresh:  time.Hour * 24,
        IdentityKey: "id",
        PayloadFunc: func(data interface{}) jwt.MapClaims {
            if v, ok := data.(string); ok {
                return jwt.MapClaims{
                    "id": v,
                }
            }
            return jwt.MapClaims{}
        },
        IdentityHandler: func(c *gin.Context) interface{} {
            claims := jwt.ExtractClaims(c)
            if id, exists := claims["id"]; exists && id != nil {
                return id.(string)
            }
            return nil
        },
        // Aquest mètode s'utilitzarà per autenticar amb usuari/contrasenya
        // Però en el nostre cas serà gestionat pel servei d'autenticació
        Authenticator: func(c *gin.Context) (interface{}, error) {
            // Aquest codi no s'utilitzarà directament, sinó a través del servei d'autenticació
            return nil, jwt.ErrFailedAuthentication
        },
        // Verificar si l'usuari té accés a una ruta específica
        Authorizator: func(data interface{}, c *gin.Context) bool {
            // Aquí pots implementar una lògica més complexa de verificació de permisos
            // Per exemple, verificar el rol de l'usuari per a rutes específiques
            return data != nil
        },
        Unauthorized: func(c *gin.Context, code int, message string) {
            c.JSON(code, gin.H{
                "code":    code,
                "message": message,
            })
        },
        TokenLookup:   "header: Authorization, query: token, cookie: jwt",
        TokenHeadName: "Bearer",
        TimeFunc:      time.Now,
    })
}