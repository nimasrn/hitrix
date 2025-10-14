package acl

//type UserRoleSetter interface {
//	SetRole(roleEntity *entity.RoleEntity)
//}
//
//func CreateRole(c *gin.Context, request *acl.CreateOrUpdateRoleRequestDTO) error {
//	ormService := service.DI().OrmForContext(c.Request.Context())
//
//	resourcesMapping, permissionsMapping, err := validateResourcesAndPermissions(ormService, request.Resources)
//	if err != nil {
//		return err
//	}
//
//	now := service.DI().Clock().Now()
//
//	roleEntity := &entity.RoleEntity{
//		Name:      request.Name,
//		CreatedAt: now,
//	}
//
//	ormService.NewEntity(roleEntity)
//
//	if err := createPrivileges(ormService, roleEntity, request.Resources, resourcesMapping, permissionsMapping, now); err != nil {
//		return err
//	}
//
//	err = ormService.FlushWithCheck()
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func UpdateRole(c *gin.Context, roleID *acl.RoleRequestDTO, request *acl.CreateOrUpdateRoleRequestDTO) error {
//	ormService := service.DI().OrmForContext(c.Request.Context())
//
//	roleEntity, found := fluxaorm.GetByID[entity.RoleEntity](ormService, roleID.ID)
//	if !found {
//		return fmt.Errorf("role with ID: %d not found", roleID.ID)
//	}
//
//	resourcesMapping, permissionsMapping, err := validateResourcesAndPermissions(ormService, request.Resources)
//	if err != nil {
//		return err
//	}
//
//	//TODO Krasi ORM: pager
//	privilegeEntitiesToDelete := fluxaorm.GetByIndex[entity.PrivilegeEntity](
//		ormService,
//		"RoleID_FakeDelete",
//		roleID,
//		false)
//
//	now := service.DI().Clock().Now()
//
//	for _, privilegeEntity := range privilegeEntitiesToDelete.All() {
//		ormService.DeleteEntity(privilegeEntity)
//	}
//
//	err = ormService.FlushWithCheck()
//	if err != nil {
//		return err
//	}
//
//	ormService.EditEntity(roleEntity)
//
//	roleEntity.Name = request.Name
//
//	if err := createPrivileges(ormService, roleEntity, request.Resources, resourcesMapping, permissionsMapping, now); err != nil {
//		return err
//	}
//
//	return ormService.FlushWithCheck()
//}
//
//func DeleteRole(c *gin.Context, roleID *acl.RoleRequestDTO) error {
//	ormService := service.DI().OrmForContext(c.Request.Context())
//
//	roleEntity, found := fluxaorm.GetByID[entity.RoleEntity](ormService, roleID.ID)
//	if !found {
//		return fmt.Errorf("role with ID: %d not found", roleID.ID)
//	}
//
//	privilegeEntitiesIterator := fluxaorm.GetByIndex[entity.PrivilegeEntity](
//		ormService,
//		"RoleID_FakeDelete",
//		roleID,
//		false,
//	)
//
//	ormService.DeleteEntity(roleEntity)
//
//	for _, privilegeEntity := range privilegeEntitiesIterator.All() {
//		ormService.DeleteEntity(privilegeEntity)
//	}
//
//	return ormService.FlushWithCheck()
//}
//
//func PostAssignRoleToUserAction(c *gin.Context, getUserFunc func() beeorm.Entity, request *acl.AssignRoleToUserRequestDTO) error {
//	ormService := service.DI().OrmForContext(c.Request.Context())
//
//	roleEntity := &entity.RoleEntity{}
//	if !ormService.LoadByID(request.RoleID, roleEntity) {
//		return fmt.Errorf("role with ID: %d not found", request.RoleID)
//	}
//
//	userEntity := getUserFunc()
//	if !ormService.LoadByID(request.UserID, userEntity) {
//		return fmt.Errorf("user with ID: %d not found", request.UserID)
//	}
//
//	userWithSettableRole, ok := userEntity.(UserRoleSetter)
//	if !ok {
//		panic("user entity does not implement UserRoleSetter interface")
//	}
//
//	userWithSettableRole.SetRole(roleEntity)
//
//	userEntityWithNewRole, _ := userWithSettableRole.(beeorm.Entity)
//
//	ormService.Flush(userEntityWithNewRole)
//
//	return nil
//}
//
//type resourceMapping map[uint64]*entity.ResourceEntity
//
//type permissionMapping map[uint64]*entity.PermissionEntity
//
//func validateResourcesAndPermissions(ormService fluxaorm.Context, resources []*acl.RoleResourceRequestDTO) (resourceMapping, permissionMapping, error) {
//	resourceIDs := make([]uint64, len(resources))
//	permissionIDs := make([]uint64, 0)
//
//	for i, resource := range resources {
//		resourceIDs[i] = resource.ResourceID
//
//		permissionIDs = append(permissionIDs, resource.PermissionIDs...)
//	}
//
//	resourcesQuery := beeorm.NewRedisSearchQuery()
//	resourcesQuery.FilterUint("ID", resourceIDs...)
//
//	resourceEntities := make([]*entity.ResourceEntity, 0)
//	ormService.RedisSearch(&resourceEntities, resourcesQuery, beeorm.NewPager(1, 1000))
//
//	if len(resourceEntities) != len(resourceIDs) {
//		return nil, nil, fmt.Errorf("some of the provided resources is not found")
//	}
//
//	allPermissionEntities := make([]*entity.PermissionEntity, 0)
//	ormService.LoadByIDs(permissionIDs, &allPermissionEntities)
//
//	permissionEntities := make([]*entity.PermissionEntity, 0)
//
//	for _, permissionEntity := range allPermissionEntities {
//		if permissionEntity.FakeDelete {
//			continue
//		}
//
//		permissionEntities = append(permissionEntities, permissionEntity)
//	}
//
//	if len(permissionEntities) != len(permissionIDs) {
//		return nil, nil, fmt.Errorf("some of the provided permissions is not found")
//	}
//
//	resourcesMapping := resourceMapping{}
//
//	for _, resourceEntity := range resourceEntities {
//		resourcesMapping[resourceEntity.ID] = resourceEntity
//	}
//
//	permissionsMapping := permissionMapping{}
//
//	for _, permissionEntity := range permissionEntities {
//		permissionsMapping[permissionEntity.ID] = permissionEntity
//	}
//
//	return resourcesMapping, permissionsMapping, nil
//}
//
//func createPrivileges(
//	ormService fluxaorm.Context,
//	roleEntity *entity.RoleEntity,
//	resources []*acl.RoleResourceRequestDTO,
//	resourcesMapping resourceMapping,
//	permissionsMapping permissionMapping,
//	now time.Time,
//) error {
//	for _, resource := range resources {
//		if len(resource.PermissionIDs) == 0 {
//			continue
//		}
//
//		resourceEntity, ok := resourcesMapping[resource.ResourceID]
//		if !ok {
//			return fmt.Errorf("resource with ID: %d not found in mapping", resource.ResourceID)
//		}
//
//		privilegeEntity := &entity.PrivilegeEntity{
//			RoleID:     roleEntity,
//			ResourceID: resourceEntity,
//			CreatedAt:  now,
//		}
//
//		for _, permissionID := range resource.PermissionIDs {
//			permissionEntity, ok := permissionsMapping[permissionID]
//			if !ok {
//				return fmt.Errorf("permission with ID: %d not found in mapping", permissionID)
//			}
//
//			if permissionEntity.ResourceID.ID != resourceEntity.ID {
//				return fmt.Errorf("permission with ID: %d does not belong to resource with ID: %d", permissionID, resource.ResourceID)
//			}
//
//			privilegeEntity.PermissionIDs = append(privilegeEntity.PermissionIDs, permissionEntity)
//		}
//
//		flusher.Track(privilegeEntity)
//	}
//
//	return nil
//}
