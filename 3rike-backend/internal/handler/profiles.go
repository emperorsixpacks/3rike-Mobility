package handler

import (
	"strconv"

	"github.com/3rike12/3rike-backend/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type DriverHandler struct{ svc domain.DriverService }

// CreateDriver godoc
// @Summary      Create driver profile
// @Tags         drivers
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      domain.Driver  true  "Driver"
// @Success      201   {object}  domain.Driver
// @Router       /api/drivers [post]
func (h *DriverHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var d domain.Driver
	if err := c.BodyParser(&d); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	result, err := h.svc.Create(c.Context(), userID, &d)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(result)
}

// GetDriver godoc
// @Summary      Get driver by ID
// @Tags         drivers
// @Security     BearerAuth
// @Param        id   path      int  true  "Driver ID"
// @Success      200  {object}  domain.Driver
// @Router       /api/drivers/{id} [get]
func (h *DriverHandler) GetByID(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	d, err := h.svc.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(d)
}

// ListDrivers godoc
// @Summary      List all drivers
// @Tags         drivers
// @Security     BearerAuth
// @Success      200  {array}  domain.Driver
// @Router       /api/drivers [get]
func (h *DriverHandler) List(c *fiber.Ctx) error {
	list, err := h.svc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}

type InvestorHandler struct{ svc domain.InvestorService }

// CreateInvestor godoc
// @Summary      Create investor profile
// @Tags         investors
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      domain.Investor  true  "Investor"
// @Success      201   {object}  domain.Investor
// @Router       /api/investors [post]
func (h *InvestorHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var i domain.Investor
	if err := c.BodyParser(&i); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	result, err := h.svc.Create(c.Context(), userID, &i)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(result)
}

// GetInvestor godoc
// @Summary      Get investor by ID
// @Tags         investors
// @Security     BearerAuth
// @Param        id   path      int  true  "Investor ID"
// @Success      200  {object}  domain.Investor
// @Router       /api/investors/{id} [get]
func (h *InvestorHandler) GetByID(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	inv, err := h.svc.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(inv)
}

// ListInvestors godoc
// @Summary      List all investors
// @Tags         investors
// @Security     BearerAuth
// @Success      200  {array}  domain.Investor
// @Router       /api/investors [get]
func (h *InvestorHandler) List(c *fiber.Ctx) error {
	list, err := h.svc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}

// ListInvestments godoc
// @Summary      List all fractions (investments) owned by an investor
// @Tags         investors
// @Security     BearerAuth
// @Param        id   path      int  true  "Investor ID"
// @Success      200  {array}   domain.Fraction
// @Router       /api/investors/{id}/investments [get]
func (h *InvestorHandler) ListInvestments(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	fractions, err := h.svc.ListInvestments(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fractions)
}
