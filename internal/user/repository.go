package user

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go-template/pkg/app"
	"gorm.io/gorm"
)

type Repository interface {
}

type repository struct {
	db   *gorm.DB
	rDB  *redis.Client
	rCtx context.Context
}

func NewRepository(app *app.AppContext) Repository {
	return &repository{
		db:   app.DB,
		rDB:  app.RDB,
		rCtx: app.RContext,
	}
}
