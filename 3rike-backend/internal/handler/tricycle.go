package handler

import (
	"strconv"

	"github.com/3rike12/3rike-backend/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type TricycleHandler struct{ svc domain.TricycleService }

// CreateTricycle godoc
// @Summary      Register a new tricycle
// @Tags         tricycles
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      domain.Tricycle  true  "Tricycle"
// @Success      201   {object}  domain.Tricycle
// @Router       /api/tricycles [post]
func (h *TricycleHandler) Create(c *fiber.Ctx) error {
	var t domain.Tricycle
	if err := c.BodyParser(&t); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	result, err := h.svc.Create(c.Context(), &t)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(result)
}

// GetTricycle godoc
// @Summary      Get tricycle by ID
// @Tags         tricycles
// @Security     BearerAuth
// @Param        id   path      int  true  "Tricycle ID"
// @Success      200  {object}  domain.Tricycle
// @Router       /api/tricycles/{id} [get]
func (h *TricycleHandler) GetByID(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	t, err := h.svc.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(t)
}

// ListTricycles godoc
// @Summary      List all tricycles
// @Tags         tricycles
// @Security     BearerAuth
// @Success      200  {array}  domain.Tricycle
// @Router       /api/tricycles [get]
func (h *TricycleHandler) List(c *fiber.Ctx) error {
	list, err := h.svc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}

// TokenizeTricycle godoc
// @Summary      Tokenize a tricycle on the Canton ledger
// @Tags         tricycles
// @Security     BearerAuth
// @Param        id   path      int  true  "Tricycle ID"
// @Success      200  {object}  domain.Tricycle
// @Router       /api/tricycles/{id}/tokenize [post]
func (h *TricycleHandler) Tokenize(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	t, err := h.svc.Tokenize(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(t)
}

// FractionalizeTricycle godoc
// @Summary      Fractionalize a tokenized tricycle into investor shares
// @Tags         tricycles
// @Security     BearerAuth
// @Param        id    path      int                              true  "Tricycle ID"
// @Param        body  body      object{total_fractions=integer}  true  "Fractions"
// @Success      200   {object}  domain.Tricycle
// @Router       /api/tricycles/{id}/fractionalize [post]
func (h *TricycleHandler) Fractionalize(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var body struct {
		TotalFractions int `json:"total_fractions"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	t, err := h.svc.Fractionalize(c.Context(), uint(id), body.TotalFractions)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(t)
}
