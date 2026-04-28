package handler

import (
	"errors"

	"github.com/3rike12/3rike-backend/internal/domain"
	"github.com/3rike12/3rike-backend/internal/service"
	"github.com/gofiber/fiber/v2"
)

type WaitlistHandler struct{ svc domain.WaitlistService }

// Join handles POST /waitlist/join.
// Idempotent on email — duplicate signups return the existing entry with 200.
func (h *WaitlistHandler) Join(c *fiber.Ctx) error {
	var body struct {
		Email      string  `json:"email"`
		Phone      string  `json:"phone,omitempty"`
		ReferredBy *string `json:"referredBy,omitempty"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_body"})
	}

	entry, total, err := h.svc.Join(c.Context(), body.Email, body.Phone, body.ReferredBy)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidEmail),
			errors.Is(err, service.ErrInvalidPhone),
			errors.Is(err, service.ErrInvalidReferrer):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal"})
		}
	}

	return c.JSON(fiber.Map{
		"entry":      entry,
		"totalCount": total,
	})
}

// Stats handles GET /waitlist/stats. Cached at the router middleware layer.
func (h *WaitlistHandler) Stats(c *fiber.Ctx) error {
	total, err := h.svc.Stats(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal"})
	}
	return c.JSON(fiber.Map{"totalCount": total})
}

// GetByCode handles GET /waitlist/:code — used by returning visitors and
// referrers to look up their status (position, referrals collected).
func (h *WaitlistHandler) GetByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing_code"})
	}
	entry, total, refs, err := h.svc.GetByCode(c.Context(), code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not_found"})
	}
	return c.JSON(fiber.Map{
		"entry":          entry,
		"totalCount":     total,
		"referralCount":  refs,
	})
}
