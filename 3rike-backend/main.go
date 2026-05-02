package main

import (
	"log"

	"github.com/3rike12/3rike-backend/cache"
	"github.com/3rike12/3rike-backend/config"
	"github.com/3rike12/3rike-backend/db"
	"github.com/3rike12/3rike-backend/internal/handler"
	"github.com/3rike12/3rike-backend/internal/service"
	"github.com/3rike12/3rike-backend/pkg/canton"
	"github.com/3rike12/3rike-backend/router"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.Load()

	database, err := db.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}

	rdb, err := cache.New(cfg.RedisURL)
	if err != nil {
		log.Printf("redis unavailable, sessions/caching disabled: %v", err)
		rdb = nil
	}

	cantonClient := canton.New(cfg.CantonURL, cfg.CantonToken)
	if cfg.OIDCTokenURL != "" && cfg.OIDCUsername != "" {
		tp := canton.NewTokenProvider(cfg.OIDCTokenURL, cfg.OIDCClientID, cfg.OIDCUsername, cfg.OIDCPassword)
		cantonClient = canton.NewWithTokenProvider(cfg.CantonURL, tp)
		log.Println("canton: using OIDC token provider")
	}
	svc := service.New(database, cantonClient, cfg.JWTSecret, rdb)
	h := handler.New(svc, cantonClient, cfg.CantonValidatorURL)

	app := fiber.New(fiber.Config{AppName: "3riKE API v1"})
	router.Register(app, h, cfg, rdb)

	log.Printf("3riKE API starting on :%s", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
