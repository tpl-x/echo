package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tpl-x/echo/internal/biz"
	"go.uber.org/zap"
)

type UserHandler struct {
	userUseCase *biz.UserUseCase
	logger      *zap.Logger
}

func NewUserHandler(userUseCase *biz.UserUseCase, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		logger:      logger.Named("user_handler"),
	}
}

func (h *UserHandler) GetUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Warn("Invalid user ID", zap.String("id", idStr), zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	user, err := h.userUseCase.GetUser(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("Failed to get user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	user, err := h.userUseCase.CreateUser(c.Request().Context(), req.Name, req.Email)
	if err != nil {
		h.logger.Error("Failed to create user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) ListUsers(c echo.Context) error {
	users, err := h.userUseCase.ListUsers(c.Request().Context())
	if err != nil {
		h.logger.Error("Failed to list users", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"users": users,
		"count": len(users),
	})
}
