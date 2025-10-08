package entity

import (
	"time"
)

type SeederEntity struct {
	ID        uint64    `orm:"table=seeder;redisCache;redisSearch=search_pool"`
	Name      string    `orm:"required;unique=Seeder_Name;searchable;"`
	CreatedAt time.Time `orm:"time=true"`
}
