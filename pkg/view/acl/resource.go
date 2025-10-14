package acl

//
//func ListResources(c *gin.Context) *acl.ResourcesResponseDTO {
//	ormService := service.DI().OrmForContext(c.Request.Context())
//
//	permissionsEntityIterator := fluxaorm.GetAll[entity.PermissionEntity](ormService)
//
//	resourceDTOsMapping := resourceDTOsMapping{}
//
//	for _, permissionEntity := range permissionsEntityIterator.All() {
//		//TODO Krasi ORM: fix after references are back
//		resourceEntity := permissionEntity.ResourceID.GetEntity(ormService)
//		dto, ok := resourceDTOsMapping[resourceEntity.ID]
//		if !ok {
//			dto = &acl.ResourceResponseDTO{
//				ID:          resourceEntity.ID,
//				Name:        resourceEntity.Name,
//				Permissions: make([]*acl.PermissionResponseDTO, 0),
//			}
//		}
//
//		dto.Permissions = append(dto.Permissions, &acl.PermissionResponseDTO{
//			ID:   permissionEntity.ID,
//			Name: permissionEntity.Name,
//		})
//
//		resourceDTOsMapping[resourceEntity.ID] = dto
//	}
//
//	resultDTOs := make([]*acl.ResourceResponseDTO, 0)
//
//	for _, resourceDTO := range resourceDTOsMapping {
//		resultDTOs = append(resultDTOs, resourceDTO)
//	}
//
//	sort.Slice(resultDTOs, func(i, j int) bool {
//		return resultDTOs[i].ID < resultDTOs[j].ID
//	})
//
//	return &acl.ResourcesResponseDTO{Resources: resultDTOs}
//}
//
//type UserRoleGetter interface {
//	GetRole() *entity.RoleEntity
//}
//
//func ListUserResources(c *gin.Context, getUserFunc func(c *gin.Context) beeorm.Entity) *acl.ResourcesResponseDTO {
//	ormService := service.DI().OrmForContext(c.Request.Context())
//
//	userEntity := getUserFunc(c)
//
//	userWithGettableRole, ok := userEntity.(UserRoleGetter)
//	if !ok {
//		panic("user entity does not implement UserRoleSetter interface")
//	}
//
//	privilegeEntities := make([]*entity.PrivilegeEntity, 0)
//
//	ormService.CachedSearchWithReferences(
//		&privilegeEntities,
//		"CachedQueryPrivilegeRoleID",
//		beeorm.NewPager(1, 4000),
//		[]interface{}{userWithGettableRole.GetRole().ID},
//		[]string{"ResourceID", "PermissionIDs"},
//	)
//
//	resourceDTOsMapping := resourceDTOsMapping{}
//
//	for _, privilegeEntity := range privilegeEntities {
//		dto, ok := resourceDTOsMapping[privilegeEntity.ResourceID.ID]
//		if !ok {
//			dto = &acl.ResourceResponseDTO{
//				ID:          privilegeEntity.ResourceID.ID,
//				Name:        privilegeEntity.ResourceID.Name,
//				Permissions: make([]*acl.PermissionResponseDTO, 0),
//			}
//		}
//
//		for _, permissionEntity := range privilegeEntity.PermissionIDs {
//			dto.Permissions = append(dto.Permissions, &acl.PermissionResponseDTO{
//				ID:   permissionEntity.ID,
//				Name: permissionEntity.Name,
//			})
//		}
//
//		resourceDTOsMapping[privilegeEntity.ResourceID.ID] = dto
//	}
//
//	resultDTOs := make([]*acl.ResourceResponseDTO, 0)
//
//	for _, resourceDTO := range resourceDTOsMapping {
//		resultDTOs = append(resultDTOs, resourceDTO)
//	}
//
//	sort.Slice(resultDTOs, func(i, j int) bool {
//		return resultDTOs[i].ID < resultDTOs[j].ID
//	})
//
//	return &acl.ResourcesResponseDTO{Resources: resultDTOs}
//}
//
//type resourceDTOsMapping map[uint64]*acl.ResourceResponseDTO
