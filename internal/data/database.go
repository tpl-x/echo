package data

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/tpl-x/echo/internal/config"
	"github.com/tpl-x/echo/internal/ent"
	"go.uber.org/fx"
	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib-x/entsqlite"
)

type Database struct {
	client *ent.Client
	logger *zap.Logger
}

func NewDatabase(cfg *config.DatabaseConfig, log *zap.Logger) (*Database, error) {
	var drv *sql.Driver
	var err error

	switch cfg.Driver {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
		drv, err = sql.Open(dialect.MySQL, dsn)
	case "sqlite3":
		drv, err = sql.Open(dialect.SQLite, cfg.Database)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// Set connection pool settings for non-SQLite databases
	if cfg.Driver != "sqlite3" {
		db := drv.DB()
		db.SetMaxIdleConns(cfg.MaxIdle)
		db.SetMaxOpenConns(cfg.MaxOpen)
	}

	// Test connection
	if err := drv.DB().Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create ent client with the driver
	client := ent.NewClient(ent.Driver(drv))

	// Run auto migration
	if err := client.Schema.Create(context.Background()); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	log.Info("Database connected and migrated successfully",
		zap.String("driver", cfg.Driver),
		zap.String("database", cfg.Database))

	return &Database{
		client: client,
		logger: log,
	}, nil
}

func (d *Database) Client() *ent.Client {
	return d.client
}

func (d *Database) Close() error {
	return d.client.Close()
}

func ProvideDatabase(lc fx.Lifecycle, cfg *config.DatabaseConfig, log *zap.Logger) (*Database, error) {
	db, err := NewDatabase(cfg, log)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Info("Closing database connection")
			return db.Close()
		},
	})

	return db, nil
}

var Module = fx.Module("data",
	fx.Provide(
		fx.Annotate(
			func(cfg *config.AppConfig) *config.DatabaseConfig {
				return &cfg.Database
			},
			fx.As(new(*config.DatabaseConfig)),
		),
		ProvideDatabase,
	),
)
