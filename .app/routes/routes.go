package routes

import (
	middleware "api-users/middleware"
	models "api-users/models"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

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

func UnmashalBody(body io.ReadCloser) map[string]interface{} {
	var values map[string]interface{}

	bbody, err := ioutil.ReadAll(body)

	if err != nil {
		return nil
	}

	json.Unmarshal([]byte(bbody), &values)
	return values
}

func InitServiceSJWT(c *gin.Context) (*middleware.SJWTClaims, *middleware.Service, bool) {

	srv, ok := c.MustGet("service").(*middleware.Service)

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
