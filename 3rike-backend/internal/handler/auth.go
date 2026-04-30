package handler

import (
	"github.com/3rike12/3rike-backend/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct{ svc domain.AuthService }

// Register godoc
// @Summary      Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      object{email=string,password=string,role=string}  true  "Register"
// @Success      201   {object}  domain.User
// @Failure      422   {object}  map[string]string
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var body struct {
		Email    string      `json:"email"`
		Password string      `json:"password"`
		Role     domain.Role `json:"role"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	user, err := h.svc.Register(c.Context(), body.Email, body.Password, body.Role)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(user)
}

// Login godoc
// @Summary      Login and get JWT + session
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      object{email=string,password=string}  true  "Login"
// @Success      200   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	token, sessionID, err := h.svc.Login(c.Context(), body.Email, body.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"token": token, "session_id": sessionID})
}

// Logout godoc
// @Summary      Logout current session
// @Tags         auth
// @Security     BearerAuth
// @Success      200  {object}  map[string]string
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	sessionID, _ := c.Locals("sessionID").(string)
	if err := h.svc.Logout(c.Context(), sessionID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "logged out"})
}

// ListSessions godoc
// @Summary      List all active sessions for current user
// @Tags         auth
// @Security     BearerAuth
// @Success      200  {array}   domain.Session
// @Router       /auth/sessions [get]
func (h *AuthHandler) ListSessions(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	sessions, err := h.svc.ListSessions(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(sessions)
}

// RevokeSession godoc
// @Summary      Revoke a specific session (log out other device)
// @Tags         auth
// @Security     BearerAuth
// @Param        sessionID  path  string  true  "Session ID"
// @Success      200        {object}  map[string]string
// @Router       /auth/sessions/{sessionID} [delete]
func (h *AuthHandler) RevokeSession(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	sessionID := c.Params("sessionID")
	if err := h.svc.RevokeSession(c.Context(), sessionID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "session revoked"})
}
