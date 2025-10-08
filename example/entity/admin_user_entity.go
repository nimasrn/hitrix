package entity

import (
	hitrixEntity "github.com/coretrix/hitrix/pkg/entity"
	"github.com/latolukasz/fluxaorm"
)

type AdminUserEntity struct {
	ID     uint64                                      `orm:"table=admin_users;log=log_db_pool;redisCache;redisSearch=search_pool"`
	RoleID fluxaorm.Reference[hitrixEntity.RoleEntity] `orm:"required;cached"`
}

func (u *AdminUserEntity) SetRole(roleEntity *hitrixEntity.RoleEntity) {
	u.RoleID = fluxaorm.Reference[hitrixEntity.RoleEntity](roleEntity.ID)
}
