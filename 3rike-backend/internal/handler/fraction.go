package handler

import (
	"errors"

	"github.com/3rike12/3rike-backend/internal/domain"
	"github.com/3rike12/3rike-backend/internal/service"
	"github.com/gofiber/fiber/v2"
)

type FractionHandler struct{ svc domain.FractionService }

// Buy handles POST /api/fractions.
//
// Body: { tricycle_id: uint, units: int }
// Auth: JWT (caller's user ID is taken from the token; no spoofing).
//
// Lazily creates an Investor profile for the caller if they don't have one,
// then records the Fraction. Returns the new Fraction on success.
func (h *FractionHandler) Buy(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var body struct {
		TricycleID uint `json:"tricycle_id"`
		Units      int  `json:"units"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_body"})
	}
	if body.TricycleID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing_tricycle_id"})
	}

	frac, err := h.svc.Buy(c.Context(), userID, body.TricycleID, body.Units)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidUnits):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_units"})
		case errors.Is(err, service.ErrTricycleNotAvailable):
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "tricycle_not_available"})
		case errors.Is(err, service.ErrInsufficientUnits):
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "insufficient_units"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal"})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(frac)
}
