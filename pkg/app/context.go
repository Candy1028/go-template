package app

import (
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type AppContext struct {
	DB       *gorm.DB
	RDB      *redis.Client
	RContext context.Context
}

func (a *AppContext) GetRContext() context.Context {
	return a.RContext
}

func (a *AppContext) SetAppContext(db *gorm.DB, rdb *redis.Client, rContext context.Context) {
	a.DB = db
	a.RDB = rdb
	a.RContext = rContext
}
func NewAppContext(db *gorm.DB, rdb *redis.Client, rContext context.Context) *AppContext {
	return &AppContext{
		DB:       db,
		RDB:      rdb,
		RContext: rContext,
	}
}
