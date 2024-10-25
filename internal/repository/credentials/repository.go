package credentials

import (
	"context"
	"github.com/upassed/upassed-authentication-service/internal/config"
	domain "github.com/upassed/upassed-authentication-service/internal/repository/model"
	"gorm.io/gorm"
	"log/slog"
	"reflect"
	"runtime"
)

type Repository interface {
	CheckDuplicatesExists(ctx context.Context, username string) (bool, error)
	Save(context.Context, domain.Credentials) error
}

type credentialsRepositoryImpl struct {
	db  *gorm.DB
	cfg *config.Config
	log *slog.Logger
}

func New(db *gorm.DB, cfg *config.Config, log *slog.Logger) Repository {
	op := runtime.FuncForPC(reflect.ValueOf(New).Pointer()).Name()

	log = log.With(
		slog.String("op", op),
	)

	return &credentialsRepositoryImpl{
		db:  db,
		cfg: cfg,
		log: log,
	}
}
