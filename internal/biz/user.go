package biz

import (
	"context"
	"fmt"
	"strings"

	"github.com/tpl-x/echo/internal/data"
	"github.com/tpl-x/echo/internal/ent"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type UserUseCase struct {
	userRepo *data.UserRepo
	logger   *zap.Logger
}

func NewUserUseCase(userRepo *data.UserRepo, logger *zap.Logger) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
		logger:   logger.Named("user_usecase"),
	}
}

func (uc *UserUseCase) GetUser(ctx context.Context, id int64) (*ent.User, error) {
	uc.logger.Info("GetUser called", zap.Int64("id", id))

	if id <= 0 {
		return nil, fmt.Errorf("invalid user id: %d", id)
	}

	user, err := uc.userRepo.GetUser(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get user", zap.Error(err), zap.Int64("id", id))
		return nil, err
	}

	return user, nil
}

func (uc *UserUseCase) CreateUser(ctx context.Context, name, email string) (*ent.User, error) {
	uc.logger.Info("CreateUser called", zap.String("name", name), zap.String("email", email))

	// Business logic validation
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	if strings.TrimSpace(email) == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	if !strings.Contains(email, "@") {
		return nil, fmt.Errorf("invalid email format")
	}

	// Call repository
	user, err := uc.userRepo.CreateUser(ctx, strings.TrimSpace(name), strings.TrimSpace(email))
	if err != nil {
		uc.logger.Error("Failed to create user", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("User created successfully", zap.Int64("id", user.ID))
	return user, nil
}

func (uc *UserUseCase) ListUsers(ctx context.Context) ([]*ent.User, error) {
	uc.logger.Info("ListUsers called")

	users, err := uc.userRepo.ListUsers(ctx)
	if err != nil {
		uc.logger.Error("Failed to list users", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("Listed users successfully", zap.Int("count", len(users)))
	return users, nil
}

var Module = fx.Module("biz",
	fx.Provide(
		data.NewUserRepo,
		NewUserUseCase,
	),
)
