package entity

import (
	"time"
)

type PermissionEntity struct {
	ID         uint64          `orm:"table=permissions;redisCache"`
	ResourceID *ResourceEntity `orm:"required;unique=ResourceID_Name_FakeDelete:1"`
	Name       string          `orm:"required;unique=ResourceID_Name_FakeDelete:3"`
	CreatedAt  time.Time       `orm:"time=true"`
	FakeDelete bool            `orm:"unique=ResourceID_Name_FakeDelete:2"`

	//CachedQueryAll        *beeorm.CachedQuery `query:"1 ORDER BY ID"`
	//CachedQueryResourceID *beeorm.CachedQuery `query:":ResourceID = ?"`
}
