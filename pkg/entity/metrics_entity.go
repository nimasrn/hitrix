package entity

import (
	"time"
)

type MetricsEntity struct {
	ID        uint64 `orm:"table=metrics"`
	AppName   string
	Metrics   string    `orm:"mediumblob"`
	CreatedAt time.Time `orm:"time=true;"`
}
