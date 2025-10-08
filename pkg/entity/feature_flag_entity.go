package entity

import (
	"time"
)

type FeatureFlagEntity struct {
	ID         uint64     `orm:"table=feature_flags;localCache"`
	Name       string     `orm:"length=100;required;unique=Name"`
	Registered bool       `orm:"index=Registered_Enabled:1"`
	Enabled    bool       `orm:"index=Registered_Enabled:2"`
	UpdatedAt  *time.Time `orm:"time=true"`
	CreatedAt  time.Time  `orm:"time=true"`

	//CachedQueryAll               *beeorm.CachedQuery `query:"1 ORDER BY ID"`
	//CachedQueryName              *beeorm.CachedQuery `queryOne:":Name = ?"`
	//CachedQueryRegisteredEnabled *beeorm.CachedQuery `query:":Registered = ? AND :Enabled = ?"`
}
