package main

import (
	"api/config"
	"api/internal/setup"
	"log"
)

func main() {
	// 1. Carregar configuració
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// 2. Inicialitzar aplicació amb dependency injection
	app, err := setup.NewApp(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to initialize app: %v", err)
	}
	defer app.Close()

	// 3. Executar servidor
	if err := app.Run(); err != nil {
		log.Fatalf("❌ Failed to run server: %v", err)
	}
}