package entity

type DevPanelUserEntity struct {
	ID       uint64 `orm:"table=dev_panel_users;redisCache;redisSearch=search_pool"`
	Email    string `orm:"unique=Email;searchable"`
	Password string
}

func (u *DevPanelUserEntity) GetUniqueFieldName() string {
	return "Email"
}

func (u *DevPanelUserEntity) GetUsername() string {
	return u.Email
}

func (u *DevPanelUserEntity) GetPassword() string {
	return u.Password
}

func (u *DevPanelUserEntity) CanAuthenticate() bool {
	return true
}
