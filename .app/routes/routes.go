package routes

import (
	models "api-users/models"
	"net/http"

	ztm "github.com/devcoons/go-ztm"

	"github.com/gin-gonic/gin"
)

func RouteGatewayRefreshUserNonce(gc *gin.Context) {
	c, s, ok := InitServiceSJWT(gc)

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

func InitServiceSJWT(c *gin.Context) (*ztm.SJWTClaims, *ztm.Service, bool) {

	srv, ok := c.MustGet("service").(*ztm.Service)

	if !ok || srv.Database == nil {
		c.IndentedJSON(http.StatusExpectationFailed, nil)
		return nil, nil, false
	}

	claims := srv.ValidateServiceJWT(c.Request)

	if claims == nil {
		return nil, nil, false
	}

	return claims, srv, true
}

func RouteGET(c *gin.Context) {

	claims, _, ok := InitServiceSJWT(c)

	if !ok || claims == nil {
		c.Data(503, "application/json", nil)
		return
	}

	claims.Hop = claims.Hop - 1

	c.IndentedJSON(200, nil)
}
