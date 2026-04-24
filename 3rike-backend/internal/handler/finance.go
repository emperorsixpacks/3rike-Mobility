package handler

import (
	"strconv"

	"github.com/3rike12/3rike-backend/internal/domain"
	"github.com/gofiber/fiber/v2"
)

// --- Payment ---

type PaymentHandler struct{ svc domain.PaymentService }

func (h *PaymentHandler) Record(c *fiber.Ctx) error {
	var p domain.Payment
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	result, err := h.svc.RecordPayment(c.Context(), &p)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(result)
}

func (h *PaymentHandler) GetByDriver(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("driverID"))
	list, err := h.svc.GetByDriver(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}

// --- Loan ---

type LoanHandler struct{ svc domain.LoanService }

func (h *LoanHandler) Apply(c *fiber.Ctx) error {
	var l domain.Loan
	if err := c.BodyParser(&l); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	result, err := h.svc.Apply(c.Context(), &l)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(result)
}

func (h *LoanHandler) GetByID(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	l, err := h.svc.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(l)
}

func (h *LoanHandler) Repay(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var body struct {
		AmountUSDC float64 `json:"amount_usdc"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	l, err := h.svc.Repay(c.Context(), uint(id), body.AmountUSDC)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(l)
}

// --- Savings ---

type SavingsHandler struct{ svc domain.SavingsService }

func (h *SavingsHandler) Deposit(c *fiber.Ctx) error {
	var body struct {
		DriverID   uint    `json:"driver_id"`
		AmountUSDC float64 `json:"amount_usdc"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	result, err := h.svc.Deposit(c.Context(), body.DriverID, body.AmountUSDC)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func (h *SavingsHandler) GetBalance(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("driverID"))
	s, err := h.svc.GetBalance(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(s)
}

// --- Yield ---

type YieldHandler struct{ svc domain.YieldService }

func (h *YieldHandler) GetByInvestor(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("investorID"))
	list, err := h.svc.GetByInvestor(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}
