package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/latolukasz/fluxaorm"
	"github.com/stretchr/testify/assert"
	"github.com/xorcare/pointer"

	entityExample "github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/pkg/dto/acl"
	"github.com/coretrix/hitrix/pkg/entity"
	aclView "github.com/coretrix/hitrix/pkg/view/acl"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/clock/mocks"
	"github.com/coretrix/hitrix/service/component/crud"
	registryMocks "github.com/coretrix/hitrix/service/registry/mocks"
)

func TestListResourcesAction(t *testing.T) {
	now := time.Unix(1, 0)

	fakeClock := &mocks.FakeSysClock{}
	fakeClock.On("Now").Return(now)

	mockServices := []*service.DefinitionGlobal{
		registryMocks.ServiceProviderMockClock(fakeClock),
	}

	ctx := createContextMyApp(t, "my-app", mockServices, nil)

	ormService := service.DI().Orm().Clone()

	resource1 := CreateResource(ormService, map[string]interface{}{})

	ormService.Flush()

	resource2 := CreateResource(ormService, map[string]interface{}{"Name": "car"})

	CreatePermission(ormService, map[string]interface{}{"ResourceID": resource1, "Name": "create"})
	CreatePermission(ormService, map[string]interface{}{"ResourceID": resource1, "Name": "view"})

	CreatePermission(ormService, map[string]interface{}{"ResourceID": resource2, "Name": "unlock"})
	CreatePermission(ormService, map[string]interface{}{"ResourceID": resource2, "Name": "lock"})
	CreatePermission(ormService, map[string]interface{}{"ResourceID": resource2, "Name": "drive"})

	ormService.Flush()

	got := &acl.ResourcesResponseDTO{}

	err := SendHTTPRequest(ctx, http.MethodGet, "/acl/resources/", false, got)
	assert.Nil(t, err)

	want := &acl.ResourcesResponseDTO{
		Resources: []*acl.ResourceResponseDTO{
			{
				ID:   1,
				Name: "user",
				Permissions: []*acl.PermissionResponseDTO{
					{
						ID:   1,
						Name: "create",
					},
					{
						ID:   2,
						Name: "view",
					},
				},
			},
			{
				ID:   2,
				Name: "car",
				Permissions: []*acl.PermissionResponseDTO{
					{
						ID:   3,
						Name: "unlock",
					},
					{
						ID:   4,
						Name: "lock",
					},
					{
						ID:   5,
						Name: "drive",
					},
				},
			},
		},
	}

	assert.Equal(t, want, got)

	fakeClock.AssertExpectations(t)
}

func TestListRolesAction(t *testing.T) {
	now := time.Unix(1, 0)

	fakeClock := &mocks.FakeSysClock{}
	fakeClock.On("Now").Return(now)

	mockServices := []*service.DefinitionGlobal{
		registryMocks.ServiceProviderMockClock(fakeClock),
	}

	ctx := createContextMyApp(t, "my-app", mockServices, nil)

	ormService := service.DI().Orm().Clone()

	resource1 := CreateResource(ormService, map[string]interface{}{})

	CreatePermission(ormService, map[string]interface{}{"ResourceID": resource1, "Name": "create"})
	CreatePermission(ormService, map[string]interface{}{"ResourceID": resource1, "Name": "view"})

	CreateRole(ormService, map[string]interface{}{})
	CreateRole(ormService, map[string]interface{}{"Name": "super-admin"})
	CreateRole(ormService, map[string]interface{}{"Name": "super-mega-admin"})

	ormService.Flush()

	got := &acl.RolesResponseDTO{}

	request := &crud.ListRequest{
		Page:     pointer.Int(1),
		PageSize: pointer.Int(2),
	}

	err := SendHTTPRequestWithBody(ctx, http.MethodPost, "/acl/roles/", request, false, got)
	assert.Nil(t, err)

	want := &acl.RolesResponseDTO{
		Total:   3,
		Columns: aclView.RolesColumns(),
		Rows: []*acl.RoleResponseDTO{
			{
				ID:   1,
				Name: "admin",
			},
			{
				ID:   2,
				Name: "super-admin",
			},
		},
		PageContext: &acl.ResourcesResponseDTO{
			Resources: []*acl.ResourceResponseDTO{
				{
					ID:   1,
					Name: "user",
					Permissions: []*acl.PermissionResponseDTO{
						{
							ID:   1,
							Name: "create",
						},
						{
							ID:   2,
							Name: "view",
						},
					},
				},
			},
		},
	}

	assert.Equal(t, want, got)

	fakeClock.AssertExpectations(t)
}

func TestGetRoleAction(t *testing.T) {
	now := time.Unix(1, 0)

	fakeClock := &mocks.FakeSysClock{}
	fakeClock.On("Now").Return(now)

	mockServices := []*service.DefinitionGlobal{
		registryMocks.ServiceProviderMockClock(fakeClock),
	}

	ctx := createContextMyApp(t, "my-app", mockServices, nil)

	ormService := service.DI().Orm().Clone()

	role := CreateRole(ormService, map[string]interface{}{})
	resource := CreateResource(ormService, map[string]interface{}{})

	permission1 := CreatePermission(ormService, map[string]interface{}{"ResourceID": resource, "Name": "create"})
	permission2 := CreatePermission(ormService, map[string]interface{}{"ResourceID": resource, "Name": "view"})

	CreatePrivilege(ormService, map[string]interface{}{"RoleID": role, "ResourceID": resource, "PermissionIDs": []*entity.PermissionEntity{
		permission1,
		permission2,
	}})

	ormService.Flush()

	got := &acl.RoleResponseDTO{}

	err := SendHTTPRequest(ctx, http.MethodGet, "/acl/role/1/", false, got)
	assert.Nil(t, err)

	want := &acl.RoleResponseDTO{
		ID:   1,
		Name: "admin",
		Resources: []*acl.ResourceResponseDTO{
			{
				ID:   1,
				Name: "user",
				Permissions: []*acl.PermissionResponseDTO{
					{
						ID:   1,
						Name: "create",
					},
					{
						ID:   2,
						Name: "view",
					},
				},
			},
		},
	}

	assert.Equal(t, want, got)

	fakeClock.AssertExpectations(t)
}

func TestCreateRoleAction(t *testing.T) {
	now := time.Unix(1, 0)

	fakeClock := &mocks.FakeSysClock{}
	fakeClock.On("Now").Return(now)

	mockServices := []*service.DefinitionGlobal{
		registryMocks.ServiceProviderMockClock(fakeClock),
	}

	ctx := createContextMyApp(t, "my-app", mockServices, nil)

	ormService := service.DI().Orm().Clone()

	resource := CreateResource(ormService, map[string]interface{}{})

	permission1 := CreatePermission(ormService, map[string]interface{}{"ResourceID": resource, "Name": "create"})
	permission2 := CreatePermission(ormService, map[string]interface{}{"ResourceID": resource, "Name": "view"})

	ormService.Flush()

	request := &acl.CreateOrUpdateRoleRequestDTO{
		Name: "admin",
		Resources: []*acl.RoleResourceRequestDTO{
			{
				ResourceID:    resource.ID,
				PermissionIDs: []uint64{permission1.ID, permission2.ID},
			},
		},
	}

	err = SendHTTPRequestWithBody(ctx, http.MethodPost, "/acl/role/", request, false, nil)
	assert.Nil(t, err)

	privilegeEntity, found := fluxaorm.GetByID[entity.PrivilegeEntity](ormService, 1)
	assert.True(t, found)

	assert.Equal(t, privilegeEntity.RoleID.GetEntity(ormService).ID, uint64(1))
	assert.Equal(t, privilegeEntity.ResourceID.GetEntity(ormService).ID, uint64(1))
	assert.Equal(t, privilegeEntity.PermissionIDs[0].ID, uint64(1))
	assert.Equal(t, privilegeEntity.PermissionIDs[1].ID, uint64(2))

	fakeClock.AssertExpectations(t)
}

func TestUpdateRoleAction(t *testing.T) {
	now := time.Unix(1, 0)

	fakeClock := &mocks.FakeSysClock{}
	fakeClock.On("Now").Return(now)

	mockServices := []*service.DefinitionGlobal{
		registryMocks.ServiceProviderMockClock(fakeClock),
	}

	ctx := createContextMyApp(t, "my-app", mockServices, nil)

	ormService := service.DI().Orm().Clone()

	role := CreateRole(ormService, map[string]interface{}{})
	resource := CreateResource(ormService, map[string]interface{}{})

	permission1 := CreatePermission(ormService, map[string]interface{}{"ResourceID": resource, "Name": "create"})
	permission2 := CreatePermission(ormService, map[string]interface{}{"ResourceID": resource, "Name": "view"})

	CreatePrivilege(ormService, map[string]interface{}{"RoleID": role, "ResourceID": resource, "PermissionIDs": []*entity.PermissionEntity{
		permission1,
		permission2,
	}})

	ormService.Flush()

	request := &acl.CreateOrUpdateRoleRequestDTO{
		Name: "super-admin",
		Resources: []*acl.RoleResourceRequestDTO{
			{
				ResourceID:    resource.ID,
				PermissionIDs: []uint64{permission2.ID},
			},
		},
	}

	err = SendHTTPRequestWithBody(ctx, http.MethodPut, "/acl/role/1/", request, false, nil)
	assert.Nil(t, err)

	_, found := fluxaorm.GetByID[entity.PrivilegeEntity](ormService, 1)
	assert.False(t, found)

	privilegeEntity, found := fluxaorm.GetByID[entity.PrivilegeEntity](ormService, 2)
	assert.True(t, found)

	assert.Equal(t, privilegeEntity.RoleID.GetEntity(ormService).ID, uint64(1))
	assert.Equal(t, privilegeEntity.RoleID.GetEntity(ormService).Name, "super-admin")
	assert.Equal(t, privilegeEntity.ResourceID.GetEntity(ormService).ID, uint64(1))
	assert.Equal(t, privilegeEntity.PermissionIDs[0].ID, permission2.ID)
	assert.Equal(t, len(privilegeEntity.PermissionIDs), 1)

	fakeClock.AssertExpectations(t)
}

func TestDeleteRoleAction(t *testing.T) {
	now := time.Unix(1, 0)

	fakeClock := &mocks.FakeSysClock{}
	fakeClock.On("Now").Return(now)

	mockServices := []*service.DefinitionGlobal{
		registryMocks.ServiceProviderMockClock(fakeClock),
	}

	ctx := createContextMyApp(t, "my-app", mockServices, nil)

	ormService := service.DI().Orm().Clone()

	role := CreateRole(ormService, map[string]interface{}{})
	resource := CreateResource(ormService, map[string]interface{}{})

	permission1 := CreatePermission(ormService, map[string]interface{}{"ResourceID": resource, "Name": "create"})
	permission2 := CreatePermission(ormService, map[string]interface{}{"ResourceID": resource, "Name": "view"})

	CreatePrivilege(ormService, map[string]interface{}{"RoleID": role, "ResourceID": resource, "PermissionIDs": []*entity.PermissionEntity{
		permission1,
		permission2,
	}})

	ormService.Flush()

	err := SendHTTPRequest(ctx, http.MethodDelete, "/acl/role/1/", false, nil)
	assert.Nil(t, err)

	roleEntity, found := fluxaorm.GetByID[entity.RoleEntity](ormService, 1)
	assert.True(t, found)
	assert.Equal(t, roleEntity.FakeDelete, true)

	privilegeEntity, found := fluxaorm.GetByID[entity.PrivilegeEntity](ormService, 1)
	assert.True(t, found)
	assert.Equal(t, privilegeEntity.FakeDelete, true)

	fakeClock.AssertExpectations(t)
}

func TestPostAssignRoleToUserAction(t *testing.T) {
	now := time.Unix(1, 0)

	fakeClock := &mocks.FakeSysClock{}
	fakeClock.On("Now").Return(now)

	mockServices := []*service.DefinitionGlobal{
		registryMocks.ServiceProviderMockClock(fakeClock),
	}

	ctx := createContextMyApp(t, "my-app", mockServices, nil)

	ormService := service.DI().Orm().Clone()

	role1 := CreateRole(ormService, map[string]interface{}{})
	role2 := CreateRole(ormService, map[string]interface{}{"Name": "super-admin"})

	user := CreateAdminUser(ormService, map[string]interface{}{"RoleID": role1})

	ormService.Flush()

	request := &acl.AssignRoleToUserRequestDTO{
		UserID: user.ID,
		RoleID: role2.ID,
	}

	err := SendHTTPRequestWithBody(ctx, http.MethodPost, "/acl/assign-role/", request, false, nil)
	assert.Nil(t, err)

	userEntity, found := fluxaorm.GetByID[entityExample.AdminUserEntity](ormService, 1)
	assert.True(t, found)

	assert.Equal(t, userEntity.RoleID.GetEntity(ormService).ID, role2.ID)
	assert.Equal(t, userEntity.RoleID.GetEntity(ormService).Name, role2.Name)

	fakeClock.AssertExpectations(t)
}
