package routes

import (
	models "api-users/models"

	ztm "github.com/devcoons/go-ztm"
	"github.com/gin-gonic/gin"
)

func RouteGatewayRefreshUserNonce(gc *gin.Context) {
	c, s, ok := ztm.InitServiceSJWT(gc)

	if !ok {
		gc.AbortWithStatus(401)
		return
	}
	if c.Service != "api-gateway" {
		gc.AbortWithStatus(401)
		return
	}

	usr := models.UsersGetById(s.Database, c.UserId)
	if usr == nil {
		gc.AbortWithStatus(404)
		return
	}
	s.ReloadUserNonceFromDB(c.UserId, usr.Nonce)
	gc.Done()
}
