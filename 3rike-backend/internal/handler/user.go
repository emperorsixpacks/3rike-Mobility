package handler

import (
	"github.com/3rike12/3rike-backend/internal/domain"
	"github.com/3rike12/3rike-backend/pkg/canton"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	svc          domain.UserService
	canton       *canton.Client
	validatorURL string
}

// Me godoc
// @Summary      Get current user profile
// @Tags         user
// @Security     BearerAuth
// @Success      200  {object}  domain.User
// @Router       /auth/me [get]
func (h *UserHandler) Me(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	u, err := h.svc.Me(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}
	return c.JSON(u)
}

// UpdateProfile godoc
// @Summary      Update email / profile
// @Tags         user
// @Security     BearerAuth
// @Param        body  body  object  true  "email"
// @Success      200   {object}  domain.User
// @Router       /auth/profile [put]
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var body struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	u, err := h.svc.UpdateProfile(c.Context(), userID, body.Email)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(u)
}

// ChangePassword godoc
// @Summary      Change password
// @Tags         user
// @Security     BearerAuth
// @Router       /auth/password [put]
func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var body struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.svc.ChangePassword(c.Context(), userID, body.OldPassword, body.NewPassword); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "password updated"})
}

// DeleteAccount godoc
// @Summary      Delete own account
// @Tags         user
// @Security     BearerAuth
// @Router       /auth/account [delete]
func (h *UserHandler) DeleteAccount(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	if err := h.svc.DeleteAccount(c.Context(), userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// WalletBalance godoc
// @Summary      Get CC balance from the user's Canton wallet
// @Tags         user
// @Security     BearerAuth
// @Success      200  {object}  canton.WalletBalance
// @Router       /auth/wallet/balance [get]
func (h *UserHandler) WalletBalance(c *fiber.Ctx) error {
	if h.validatorURL == "" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "validator URL not configured"})
	}
	bal, err := h.canton.GetWalletBalance(c.Context(), h.validatorURL, c.Get("Authorization"))
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(bal)
}

// LinkWallet godoc
// @Summary      Link a Canton wallet party ID to the current user
// @Tags         user
// @Security     BearerAuth
// @Param        body  body  object{canton_party_id=string}  true  "Party ID"
// @Success      200   {object}  domain.User
// @Router       /auth/wallet [put]
func (h *UserHandler) LinkWallet(c *fiber.Ctx) error {	userID := c.Locals("userID").(uint)
	var body struct {
		CantonPartyID string `json:"canton_party_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if body.CantonPartyID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "canton_party_id required"})
	}
	u, err := h.svc.LinkWallet(c.Context(), userID, body.CantonPartyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(u)
}
