package entity

import (
	"time"
)

type ResourceEntity struct {
	ID         uint64    `orm:"table=resources;redisCache;redisSearch=search_pool;searchable"`
	Name       string    `orm:"required;searchable;unique=Name_FakeDelete:1;cached"`
	CreatedAt  time.Time `orm:"time=true"`
	FakeDelete bool      `orm:"unique=Name_FakeDelete:2"`

	//CachedQueryName *beeorm.CachedQuery `queryOne:":Name = ?"`
}
