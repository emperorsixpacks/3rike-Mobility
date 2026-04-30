package handler

import (
	"strconv"

	"github.com/3rike12/3rike-backend/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct{ svc domain.PaymentService }

// RecordPayment godoc
// @Summary      Record a driver weekly repayment
// @Tags         payments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      domain.Payment  true  "Payment"
// @Success      201   {object}  domain.Payment
// @Router       /api/payments [post]
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

// GetPaymentsByDriver godoc
// @Summary      Get all payments for a driver
// @Tags         payments
// @Security     BearerAuth
// @Param        driverID  path      int  true  "Driver ID"
// @Success      200       {array}   domain.Payment
// @Router       /api/payments/driver/{driverID} [get]
func (h *PaymentHandler) GetByDriver(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("driverID"))
	list, err := h.svc.GetByDriver(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}

type LoanHandler struct{ svc domain.LoanService }

// ApplyLoan godoc
// @Summary      Apply for a loan (requires credit score >= 500)
// @Tags         loans
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      domain.Loan  true  "Loan"
// @Success      201   {object}  domain.Loan
// @Router       /api/loans [post]
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

// GetLoan godoc
// @Summary      Get loan by ID
// @Tags         loans
// @Security     BearerAuth
// @Param        id   path      int  true  "Loan ID"
// @Success      200  {object}  domain.Loan
// @Router       /api/loans/{id} [get]
func (h *LoanHandler) GetByID(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	l, err := h.svc.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(l)
}

// RepayLoan godoc
// @Summary      Make a loan repayment
// @Tags         loans
// @Security     BearerAuth
// @Param        id    path      int                            true  "Loan ID"
// @Param        body  body      object{amount_usdc=number}     true  "Repayment"
// @Success      200   {object}  domain.Loan
// @Router       /api/loans/{id}/repay [put]
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

type SavingsHandler struct{ svc domain.SavingsService }

// Deposit godoc
// @Summary      Deposit USDC into driver savings account
// @Tags         savings
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      object{driver_id=integer,amount_usdc=number}  true  "Deposit"
// @Success      200   {object}  domain.Savings
// @Router       /api/savings/deposit [post]
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

// GetSavingsBalance godoc
// @Summary      Get savings balance for a driver
// @Tags         savings
// @Security     BearerAuth
// @Param        driverID  path      int  true  "Driver ID"
// @Success      200       {object}  domain.Savings
// @Router       /api/savings/{driverID}/balance [get]
func (h *SavingsHandler) GetBalance(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("driverID"))
	s, err := h.svc.GetBalance(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(s)
}

type YieldHandler struct{ svc domain.YieldService }

// GetYieldByInvestor godoc
// @Summary      Get all yield payouts for an investor
// @Tags         yield
// @Security     BearerAuth
// @Param        investorID  path      int  true  "Investor ID"
// @Success      200         {array}   domain.YieldPayout
// @Router       /api/yield/investor/{investorID} [get]
func (h *YieldHandler) GetByInvestor(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("investorID"))
	list, err := h.svc.GetByInvestor(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}
