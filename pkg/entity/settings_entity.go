package entity

const (
	SettingsValueTypeText     = "text"
	SettingsValueTypeNumber   = "number"
	SettingsValueTypeEmail    = "email"
	SettingsValueTypeTel      = "tel"
	SettingsValueTypeURI      = "uri"
	SettingsValueTypePassword = "password"
	SettingsValueTypeDateTime = "datetime"
	SettingsValueTypeJSON     = "json"
)

type settingsValueType struct {
	SettingsValueTypeText     string
	SettingsValueTypeNumber   string
	SettingsValueTypeEmail    string
	SettingsValueTypeTel      string
	SettingsValueTypeURI      string
	SettingsValueTypePassword string
	SettingsValueTypeDateTime string
	SettingsValueTypeJSON     string
}

var SettingsValueTypeAll = settingsValueType{
	SettingsValueTypeText:     SettingsValueTypeText,
	SettingsValueTypeNumber:   SettingsValueTypeNumber,
	SettingsValueTypeEmail:    SettingsValueTypeEmail,
	SettingsValueTypeTel:      SettingsValueTypeTel,
	SettingsValueTypeURI:      SettingsValueTypeURI,
	SettingsValueTypePassword: SettingsValueTypePassword,
	SettingsValueTypeDateTime: SettingsValueTypeDateTime,
	SettingsValueTypeJSON:     SettingsValueTypeJSON,
}

type SettingsEntity struct {
	ID        uint64 `orm:"table=settings;redisCache"`
	Key       string `orm:"required;unique=SettingsKey"`
	Value     string `orm:"required;length=max"`
	ValueType string `orm:"enum=entity.SettingsValueTypeAll"`
	Editable  bool
	Deletable bool
	Hidden    bool

	//CachedQuerySettingsKey *beeorm.CachedQuery `queryOne:":Key = ?"`
}
