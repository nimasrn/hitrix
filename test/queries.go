package main

import (
	"time"

	"github.com/latolukasz/fluxaorm"

	"github.com/coretrix/hitrix/example/entity"
	entityHitrix "github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service"
)

func CreateAdminUser(ormService fluxaorm.Context, row map[string]interface{}) *entity.AdminUserEntity {
	adminUserEntity := &entity.AdminUserEntity{}

	if len(row) != 0 {
		for field, value := range row {
			switch field {
			case "RoleID":
				adminUserEntity.RoleID = fluxaorm.Reference[entityHitrix.RoleEntity](value.(*entityHitrix.RoleEntity).ID)
			}
		}
	}

	fluxaorm.NewEntityFromSource(ormService, adminUserEntity)

	return adminUserEntity
}

func CreateRole(ormService fluxaorm.Context, row map[string]interface{}) *entityHitrix.RoleEntity {
	roleEntity := &entityHitrix.RoleEntity{
		Name:      "admin",
		CreatedAt: service.DI().Clock().Now(),
	}

	if len(row) != 0 {
		for field, value := range row {
			switch field {
			case "Name":
				roleEntity.Name = value.(string)
			case "CreatedAt":
				roleEntity.CreatedAt = value.(time.Time)
			}
		}
	}

	fluxaorm.NewEntityFromSource(ormService, roleEntity)

	return roleEntity
}

func CreateResource(ormService fluxaorm.Context, row map[string]interface{}) *entityHitrix.ResourceEntity {
	resourceEntity := &entityHitrix.ResourceEntity{
		Name:      "user",
		CreatedAt: service.DI().Clock().Now(),
	}

	if len(row) != 0 {
		for field, value := range row {
			switch field {
			case "Name":
				resourceEntity.Name = value.(string)
			case "CreatedAt":
				resourceEntity.CreatedAt = value.(time.Time)
			}
		}
	}

	fluxaorm.NewEntityFromSource(ormService, resourceEntity)

	return resourceEntity
}

func CreatePermission(ormService fluxaorm.Context, row map[string]interface{}) *entityHitrix.PermissionEntity {
	permissionEntity := &entityHitrix.PermissionEntity{
		Name:      "view",
		CreatedAt: service.DI().Clock().Now(),
	}

	if len(row) != 0 {
		for field, value := range row {
			switch field {
			case "ResourceID":
				permissionEntity.ResourceID = fluxaorm.Reference[entityHitrix.ResourceEntity](value.(*entityHitrix.ResourceEntity).ID)
			case "Name":
				permissionEntity.Name = value.(string)
			case "CreatedAt":
				permissionEntity.CreatedAt = value.(time.Time)
			}
		}
	}

	fluxaorm.NewEntityFromSource(ormService, permissionEntity)

	return permissionEntity
}

func CreatePrivilege(ormService fluxaorm.Context, row map[string]interface{}) *entityHitrix.PrivilegeEntity {
	privilegeEntity := &entityHitrix.PrivilegeEntity{
		PermissionIDs: nil,
		CreatedAt:     service.DI().Clock().Now(),
	}

	if len(row) != 0 {
		for field, value := range row {
			switch field {
			case "RoleID":
				privilegeEntity.RoleID = fluxaorm.Reference[entityHitrix.RoleEntity](value.(*entityHitrix.RoleEntity).ID)
			case "ResourceID":
				privilegeEntity.ResourceID = fluxaorm.Reference[entityHitrix.ResourceEntity](value.(*entityHitrix.ResourceEntity).ID)
			case "PermissionIDs":
				privilegeEntity.PermissionIDs = value.([]*entityHitrix.PermissionEntity)
			case "CreatedAt":
				privilegeEntity.CreatedAt = value.(time.Time)
			}
		}
	}

	fluxaorm.NewEntityFromSource(ormService, privilegeEntity)

	return privilegeEntity
}
