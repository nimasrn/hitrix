package account

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/latolukasz/fluxaorm"

	"github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/pkg/binding"
	"github.com/coretrix/hitrix/pkg/view/account"
	"github.com/coretrix/hitrix/service"
)

type LoginDevForm struct {
	Username string `binding:"required,min=6,max=60" json:"Username"  form:"username"`
	Password string `binding:"required,min=8,max=60" json:"Password"  form:"password"`
}

func (l *LoginDevForm) Login(c *gin.Context) (string, string, error) {
	err := binding.ShouldBindJSON(c, l)
	if err != nil {
		return "", "", err
	}

	ormService := service.DI().OrmForContext(c.Request.Context())
	devPanelUserEntity, found := fluxaorm.GetByUniqueIndex[entity.DevPanelUserEntity](ormService, "Username", l.Username)

	if !found {
		return "", "", errors.New("invalid username or password")
	}

	//TODO Krasi ORM: check possible null
	if !service.DI().Password().VerifyPassword(l.Password, devPanelUserEntity.Password) {
		return "", "", errors.New("invalid username or password")
	}

	token, refreshToken, err := account.GenerateDevTokenAndRefreshToken(ormService, devPanelUserEntity.GetID())
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, err
}
