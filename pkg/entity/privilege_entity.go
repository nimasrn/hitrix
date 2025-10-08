package entity

import (
	"time"

	"github.com/latolukasz/fluxaorm"
)

type PrivilegeEntity struct {
	ID            uint64                             `orm:"table=privileges;redisCache"`
	RoleID        fluxaorm.Reference[RoleEntity]     `orm:"required;unique=RoleID_ResourceID_FakeDelete:1;cached"`
	ResourceID    fluxaorm.Reference[ResourceEntity] `orm:"required;unique=RoleID_ResourceID_FakeDelete:2;cached"`
	PermissionIDs []*PermissionEntity                `orm:"required"`
	CreatedAt     time.Time                          `orm:"time=true"`
	FakeDelete    bool                               `orm:"unique=RoleID_ResourceID_FakeDelete:3"`

	//CachedQueryPrivilegeRoleIDResourceID *beeorm.CachedQuery `queryOne:":RoleID = ? AND :ResourceID = ?"`
	//CachedQueryPrivilegeRoleID           *beeorm.CachedQuery `query:":RoleID = ?"`
	//CachedQueryPrivilegeResourceID       *beeorm.CachedQuery `query:":ResourceID = ?"`
}
