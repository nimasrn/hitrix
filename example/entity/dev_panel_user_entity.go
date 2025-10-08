package entity

type DevPanelUserEntity struct {
	ID       uint64 `orm:"table=dev_panel_users;localCache;redisCache;"`
	Username string `orm:"unique=Username;searchable;cached"`
	Password string
}

func (u *DevPanelUserEntity) GetID() uint64 {
	return u.ID
}

func (u *DevPanelUserEntity) GetUniqueFieldName() string {
	return "Email"
}

func (u *DevPanelUserEntity) GetUsername() string {
	return u.Username
}

func (u *DevPanelUserEntity) GetPassword() string {
	return u.Password
}

func (u *DevPanelUserEntity) CanAuthenticate() bool {
	return true
}
