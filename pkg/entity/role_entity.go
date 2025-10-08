package entity

import (
	"time"
)

type RoleEntity struct {
	ID           uint64    `orm:"table=roles;redisCache;redisSearch=search_pool;sortable"`
	Name         string    `orm:"required;searchable;unique=Name_FakeDelete:1"`
	IsPredefined bool      `orm:"searchable"`
	CreatedAt    time.Time `orm:"time=true"`
	FakeDelete   bool      `orm:"unique=Name_FakeDelete:2"`
}
