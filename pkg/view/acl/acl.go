package acl

import (
	"github.com/latolukasz/fluxaorm"

	"github.com/coretrix/hitrix/pkg/entity"
)

func ACL(ormService fluxaorm.Context, roleEntity *entity.RoleEntity, resource string, permissions ...string) bool {
	resourceEntity, found := fluxaorm.GetByUniqueIndex[entity.ResourceEntity](ormService, "Name", resource)
	if !found {
		return false
	}

	allPermissionEntitiesIterator := fluxaorm.GetByIndex[entity.PermissionEntity](ormService, "ResourceID", resourceEntity.ID)

	permissionEntities := make([]*entity.PermissionEntity, 0)

	for _, permissionEntity := range allPermissionEntitiesIterator.All() {
		for _, permission := range permissions {
			if permissionEntity.Name == permission {
				permissionEntities = append(permissionEntities, permissionEntity)
			}
		}
	}

	if len(permissions) != len(permissionEntities) {
		return false
	}

	permissionIDs := make([]uint64, len(permissionEntities))

	for i, permissionEntity := range permissionEntities {
		permissionIDs[i] = permissionEntity.ID
	}

	privilegeEntity, found := fluxaorm.GetByUniqueIndex[entity.PrivilegeEntity](
		ormService,
		"RoleID_ResourceID_FakeDelete",
		roleEntity.ID,
		resourceEntity.ID,
		0,
	)
	if !found {
		return false
	}

	hasPrivilege := false

	for _, permissionEntity := range privilegeEntity.PermissionIDs {
		for _, permissionID := range permissionIDs {
			if permissionEntity.ID == permissionID {
				hasPrivilege = true

				break
			}
		}
	}

	return hasPrivilege
}
