package data

import (
	"context"
	"fmt"

	"github.com/tpl-x/echo/internal/ent"
	"github.com/tpl-x/echo/internal/ent/user"
	"go.uber.org/zap"
)

type UserRepo struct {
	db     *Database
	logger *zap.Logger
}

func NewUserRepo(db *Database) *UserRepo {
	return &UserRepo{
		db:     db,
		logger: db.logger.Named("user_repo"),
	}
}

func (r *UserRepo) GetUser(ctx context.Context, id int64) (*ent.User, error) {
	r.logger.Info("Getting user", zap.Int64("user_id", id))

	u, err := r.db.client.User.
		Query().
		Where(user.ID(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("user not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return u, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, name, email string) (*ent.User, error) {
	r.logger.Info("Creating user", zap.String("name", name), zap.String("email", email))

	u, err := r.db.client.User.
		Create().
		SetName(name).
		SetEmail(email).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	r.logger.Info("User created successfully", zap.Int64("id", u.ID))
	return u, nil
}

func (r *UserRepo) ListUsers(ctx context.Context) ([]*ent.User, error) {
	r.logger.Info("Listing all users")

	users, err := r.db.client.User.
		Query().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	r.logger.Info("Listed users successfully", zap.Int("count", len(users)))
	return users, nil
}
