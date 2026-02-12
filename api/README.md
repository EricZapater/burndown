# Burndown API

API backend per l'aplicació Burndown, construïda amb Go i Gin.

## 🏗️ Arquitectura

L'aplicació segueix una arquitectura de **Dependency Injection** amb separació clara de responsabilitats:

```
cmd/api/main.go          → Bootstrap mínim
internal/setup/setup.go  → Dependency Injection
internal/app/app.go      → Contenidor de l'aplicació
internal/*/              → Mòduls (users, auth, meals, etc.)
```

Per més detalls, consulta el [walkthrough de l'arquitectura](file:///.gemini/antigravity/brain/29858abc-e22b-44f6-b4de-76ea5666bd63/walkthrough.md).

## 🚀 Inici Ràpid

### Prerequisits

- Go 1.21+
- PostgreSQL
- Variables d'entorn configurades (`.env`)

### Configuració

1. Copia el fitxer `.env.example` a `.env`:

```bash
cp .env.example .env
```

2. Edita `.env` amb les teves credencials:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=your_password
DB_NAME=burndown
API_PORT=8080
JWT_KEY=your_secret_key
LLM_KEY=your_gemini_api_key
```

### Executar l'aplicació

```bash
# Desenvolupament
go run ./cmd/api/main.go

# Build
go build -o bin/api.exe ./cmd/api

# Executar el binari
./bin/api.exe
```

## 📡 API Endpoints

### Públics

- `GET /health` - Health check
- `POST /api/v1/users` - Registre d'usuari
- `POST /api/v1/auth/login` - Login
- `GET /api/v1/auth/refresh_token` - Refresh token

### Protegits (requereixen JWT)

#### Users

- `GET /api/v1/users/:id` - Obtenir usuari
- `PUT /api/v1/users/:id/password` - Canviar password
- `DELETE /api/v1/users/:id` - Desactivar usuari

#### Macros

- `POST /api/v1/users/:id/macros` - Afegir macros
- `GET /api/v1/users/:id/macros` - Obtenir macros actius
- `GET /api/v1/users/:id/macros/historical` - Historial de macros
- `PUT /api/v1/users/:id/macros/:macroId` - Actualitzar macros
- `DELETE /api/v1/users/:id/macros/:macroId` - Eliminar macros
- `PATCH /api/v1/users/:id/macros/:macroId/inactive` - Desactivar macros

## 🧪 Testing

```bash
# Executar tots els tests
go test ./...

# Amb coverage
go test -cover ./...

# Amb verbose
go test -v ./...
```

## 📦 Estructura del Projecte

```
api/
├── cmd/
│   └── api/
│       └── main.go              # Entry point
├── config/
│   └── config.go                # Configuració
├── internal/
│   ├── app/
│   │   └── app.go               # App container
│   ├── setup/
│   │   └── setup.go             # Dependency injection
│   ├── database/
│   │   └── postgres.go          # DB connection
│   ├── middleware/
│   │   └── jwt_middleware.go    # JWT config
│   ├── users/
│   │   ├── model.go             # Models
│   │   ├── repository.go        # Data access
│   │   ├── service.go           # Business logic
│   │   ├── handler.go           # HTTP handlers
│   │   ├── routes.go            # Route registration
│   │   └── errors.go            # Custom errors
│   ├── auth/
│   │   └── ...                  # Auth module
│   └── advisor/
│       └── gemini.go            # Gemini AI integration
└── .env                         # Environment variables
```

## 🔧 Afegir un Nou Mòdul

1. Crear directori `internal/newmodule/`
2. Crear fitxers: `model.go`, `repository.go`, `service.go`, `handler.go`, `routes.go`
3. Afegir al `setup.go`:

```go
// Repository
newRepo := newmodule.NewRepository(db)

// Service
newService := newmodule.NewService(newRepo)

// Handler
newHandler := newmodule.NewHandler(newService)

// Routes (dins setupRoutes)
newmodule.RegisterRoutes(protected, newHandler)
```

**No cal tocar `main.go`!**

## 🛠️ Tecnologies

- **Framework**: [Gin](https://github.com/gin-gonic/gin)
- **Database**: PostgreSQL
- **Authentication**: JWT ([gin-jwt](https://github.com/appleboy/gin-jwt))
- **AI**: Google Gemini
- **Config**: [godotenv](https://github.com/joho/godotenv) + [env](https://github.com/caarlos0/env)

## 📝 Convencions

- Cada mòdul té la seva pròpia carpeta a `internal/`
- Separació clara: Repository (data) → Service (logic) → Handler (HTTP)
- Tots els mètodes accepten `context.Context` com a primer paràmetre
- Errors customitzats per mòdul (`errors.go`)
- Routes registrades amb `RegisterRoutes(router, handler)`

## 🔐 Seguretat

- Passwords encriptats amb bcrypt
- JWT per autenticació
- CORS configurable
- Graceful shutdown per evitar pèrdua de dades

## 📚 Recursos

- [Documentació de Gin](https://gin-gonic.com/docs/)
- [Go Best Practices](https://golang.org/doc/effective_go)
- [Walkthrough de l'arquitectura](file:///.gemini/antigravity/brain/29858abc-e22b-44f6-b4de-76ea5666bd63/walkthrough.md)
