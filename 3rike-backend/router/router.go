package router

import (
	"time"

	_ "github.com/3rike12/3rike-backend/docs"
	"github.com/3rike12/3rike-backend/config"
	"github.com/3rike12/3rike-backend/internal/handler"
	"github.com/3rike12/3rike-backend/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	swagger "github.com/gofiber/swagger"
	"github.com/redis/go-redis/v9"
)

func Register(app *fiber.App, h *handler.Handlers, cfg *config.Config, rdb *redis.Client) {
	app.Use(cors.New())
	app.Use(logger.New())

	jwtAuth := middleware.Auth(cfg.JWTSecret, rdb)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Swagger — JWT protected
	app.Get("/docs/*", jwtAuth, swagger.HandlerDefault)

	// ── Public ───────────────────────────────────────────────────────────────
	auth := app.Group("/auth")
	auth.Post("/register", h.Auth.Register)
	auth.Post("/login", h.Auth.Login)

	// ── User self-management (JWT required) ──────────────────────────────────
	me := app.Group("/auth", jwtAuth)
	me.Post("/logout", h.Auth.Logout)
	me.Get("/sessions", h.Auth.ListSessions)
	me.Delete("/sessions/:sessionID", h.Auth.RevokeSession)
	me.Get("/me", h.User.Me)
	me.Put("/profile", h.User.UpdateProfile)
	me.Put("/password", h.User.ChangePassword)
	me.Delete("/account", h.User.DeleteAccount)

	// ── Protected API ─────────────────────────────────────────────────────────
	api := app.Group("/api", jwtAuth)

	drivers := api.Group("/drivers")
	drivers.Post("/", h.Driver.Create)
	drivers.Get("/", middleware.Cache(rdb, 5*time.Minute), h.Driver.List)
	drivers.Get("/:id", middleware.Cache(rdb, 5*time.Minute), h.Driver.GetByID)

	investors := api.Group("/investors")
	investors.Post("/", h.Investor.Create)
	investors.Get("/", middleware.Cache(rdb, 5*time.Minute), h.Investor.List)
	investors.Get("/:id", middleware.Cache(rdb, 5*time.Minute), h.Investor.GetByID)
	investors.Get("/:id/investments", middleware.Cache(rdb, 2*time.Minute), h.Investor.ListInvestments)

	tricycles := api.Group("/tricycles")
	tricycles.Post("/", h.Tricycle.Create)
	tricycles.Get("/", middleware.Cache(rdb, 5*time.Minute), h.Tricycle.List)
	tricycles.Get("/:id", middleware.Cache(rdb, 5*time.Minute), h.Tricycle.GetByID)
	tricycles.Post("/:id/tokenize", h.Tricycle.Tokenize)
	tricycles.Post("/:id/fractionalize", h.Tricycle.Fractionalize)

	payments := api.Group("/payments")
	payments.Post("/", h.Payment.Record)
	payments.Get("/driver/:driverID", middleware.Cache(rdb, 2*time.Minute), h.Payment.GetByDriver)

	loans := api.Group("/loans")
	loans.Post("/", h.Loan.Apply)
	loans.Get("/:id", middleware.Cache(rdb, 5*time.Minute), h.Loan.GetByID)
	loans.Put("/:id/repay", h.Loan.Repay)

	savings := api.Group("/savings")
	savings.Post("/deposit", h.Savings.Deposit)
	savings.Get("/:driverID/balance", middleware.Cache(rdb, 2*time.Minute), h.Savings.GetBalance)

	yields := api.Group("/yield")
	yields.Get("/investor/:investorID", middleware.Cache(rdb, 5*time.Minute), h.Yield.GetByInvestor)
}
