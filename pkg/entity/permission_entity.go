package entity

import (
	"time"

	"github.com/latolukasz/fluxaorm"
)

type PermissionEntity struct {
	ID         uint64                             `orm:"table=permissions;redisCache"`
	ResourceID fluxaorm.Reference[ResourceEntity] `orm:"required;unique=ResourceID_Name_FakeDelete:1;cached"`
	Name       string                             `orm:"required;unique=ResourceID_Name_FakeDelete:3"`
	CreatedAt  time.Time                          `orm:"time=true"`
	FakeDelete bool                               `orm:"unique=ResourceID_Name_FakeDelete:2"`
}
